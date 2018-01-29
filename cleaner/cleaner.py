#!/usr/bin/python
# -*- coding: utf8 -*-
import sys
import traceback
import logging
import re
import json
import threading
import requests
import time
import copy
from vultr import Vultr
from multiprocessing.dummy import Pool as ThreadPool

sys.path.append("..")
from worker import Worker
from httpstatus import HTTPStatus


class Cleaner(Worker):
    def __init__(self):
        super().__init__(__file__)
        self.type = "cleaner"
        self.tasklog = ""
        self.tmpSnapshotID = ""
        self.targetVPS = None
        self.queryInterval = int(self.config.get('cleaner', 'queryInterval'))
        self.createMinTime = int(self.config.get('cleaner', 'createMinTime'))
        self.destroyWaitTime = int(self.config.get('cleaner', 'destroyWaitTime'))
        self.timeZone = int(self.config.get('cleaner', 'timeZone'))
        self.oldVPSList = []
        self.vultrDestroyCoolTime = int(self.config.get('cleaner', 'vultrDestroyCoolTime'))
        self.destroyPool = ThreadPool(processes=self.maxThreads * 4)
        self.destroyResults = []
        self.VULTR512MPLANID = '200'
        self.VULTR1024MPLANID = '201'
        self.lastException = None
        # Init Vultr api instance
        self.vultrApikey = self.config.get('cleaner', 'vultrApikey')
        self.vultr = Vultr(self.vultrApikey)
        # Function dic for different VPS providers
        self.create_tmp_snapshot = {
            'Vultr': self.create_tmp_snapshot_vultr
        }
        self.destroy_tmp_snapshot = {
            'Vultr': self.destroy_tmp_snapshot_vultr
        }
        self.get_server_info = {
            'Vultr': self.get_server_info_vultr
        }
        self.destroy_and_create = {
            'Vultr': self.destroy_and_create_vultr
        }
        self.get_server_ip = {
            'Vultr': self.get_server_ip_vultr
        }

        # These kawaii variables↓ only works for Threads in Threads Pool
        self.cleaning = False

    def get_tasks(self):
        rst = None
        try:
            rst = self._GET(path="task", data_dict={'class': self.type, 'state': 'Queuing'})
            if not (rst['code'] == HTTPStatus.OK and rst['result']):
                logging.error('Get tasks result: %s' % rst)
                raise Exception('Invalid get tasks result')
        except:
            logging.error('Failed while getting tasks')
            traceback.print_exc(file=sys.stderr)
        return rst['data']

    def clean(self, task):
        provider = task['Node']['Provider']
        need_destroy_snapshot = False
        # Get VPS info
        if not self.get_server_info[provider](task):
            return False
        # Set target snapshot ID, create temporary snapshot if needed
        if task['Node']['Snapshot'] is None or task['Node']['Snapshot'] == '':
            if not self.create_tmp_snapshot[provider](task):
                return False
            need_destroy_snapshot = True
        else:
            self.tmpSnapshotID = task['Node']['Snapshot']
        # Main clean loop
        shiny = False
        while not shiny:
            # Destroy and Create VPS
            if not self.destroy_and_create[provider](task):
                return False
            # Call watcher
            watcher_task_dict = {
                'callback_id': 0,
                'class': 'watcher',
                'node_id': 0,
                'ip_ver': 4,
                'server_name': self.get_server_ip[provider]()['ipv4']
            }
            rst = self.assign_task_and_wait(task=task, new_task_dict=watcher_task_dict)
            # Check watch task result
            if isinstance(rst, bool) and not rst:
                return False
            if rst['State'] == 'Passing':
                shiny = True
        # Destroy temporary snapshot if needed
        if need_destroy_snapshot:
            if not self.destroy_tmp_snapshot[provider](task):
                return False
        # Update task state and returun True
        return True

    def assign_task_and_wait(self, task, new_task_dict=None):
        # Assign new task
        logging.debug('Assigning new task for %s for node %s...' % (new_task_dict['class'], task['Node']['Name']))
        rst = self._POST(path='task', data_dict=new_task_dict)
        if rst is not None and rst.code == HTTPStatus.OK and rst['result']:
            logging.debug('Successfully assigned new task for %s for node %s, task info: %s'
                          % (new_task_dict['class'], task['Node']['Name'], rst['data']))
            new_task = rst['data']
        else:
            logging.error(
                'Failed while assigning new task for %s for node %s' % (new_task_dict['class'], task['Node']['Name']))
            return False
        # Wait for the result
        logging.debug(
            'Waiting for result of new task for %s for node %s' % (new_task_dict['class'], task['Node']['Name']))
        while True:
            rst = self._GET(path='task/' + str(new_task['ID']), data_dict={})
            if rst is not None and rst.code == HTTPStatus.OK and rst['result']:
                logging.debug('The state for new task for %s for node %s is %s' % (
                    new_task_dict['class'], task['Node']['Name'], rst['data']['State']))
                if rst['data']['State'] in ['Passing', 'Failing']:
                    logging.info('New task for %s for node %s is done and the result is %s' % (
                        new_task_dict['class'], task['Node']['Name'], rst['data']['State']))
                    return rst['data']
            else:
                logging.error(
                    'Failed while waiting for the result of new task for %s for node %s' % (
                        new_task_dict['class'], task['Node']['Name']))
                return False
            time.sleep(self.queryInterval)

    def do_safely(self, method, description, args, depth=0):
        try:
            return method(**args)
        except Exception as e:
            traceback.print_exc(file=sys.stderr)
            if depth < self.maxTry:
                logging.warning('Failed while %s after trying %d times, retrying...'
                                % (description, (depth + 1)))
                return self.do_safely(method, description, args, depth + 1)
            else:
                logging.error('Failed while %s after trying %d times, exiting...'
                              % (description, self.maxTry))
                self.lastException = e
                return False

    def get_server_info_vultr(self, task, subid=None):
        # Get server info from task's IPv4 or subid
        logging.debug("Starting getting server info for node %s..." % task['Node']['Name'])
        servers = self.do_safely(method=self.vultr.server.list,
                                 description="getting server info for node %s" % task['Node']['Name'],
                                 args={
                                     'subid': subid,
                                     'params': {'main_ip': task['Node']['IPv4']} if subid is None else None
                                 })
        if isinstance(servers, bool) and not servers:
            return False
        if subid is None:
            for k in servers:
                if (subid is None and servers[k]['main_ip'] == task['Node']['IPv4']) or (subid == k):
                    logging.debug("Server info for ip %s: %s" % (task['Node']['IPv4'], servers[k]))
                    self.targetVPS = servers[k]
                    return True
            logging.error("Failed while getting server info for node %s" % (task['Node']['Name']))
            return False
        else:
            logging.debug("Server info for subid %s: %s" % (subid, servers))
            self.targetVPS = servers
            return True

    def get_server_ip_vultr(self):
        return {
            'ipv4': self.targetVPS['main_ip'],
            'ipv6': self.targetVPS['v6_networks'][0]['v6_main_ip']
        }

    def create_tmp_snapshot_vultr(self, task):
        # Start create snapshot
        logging.debug("Starting creating snapshot for node %s..." % task['Node']['Name'])
        rst = self.do_safely(method=self.vultr.snapshot.create,
                             description="creating snapshot for node %s" % task['Node']['Name'],
                             args={
                                 'subid': self.targetVPS['SUBID'],
                                 'params': {'description': 'Temporary snapshot created by %s for node %s'
                                                           % (self.name, task['Node']['Name'])}
                             })
        if isinstance(rst, bool) and not rst:
            return False
        logging.debug("Snapshot id for node %s is %s" % (task['Node']['Name'], rst['SNAPSHOTID']))
        self.tmpSnapshotID = rst['SNAPSHOTID']
        # Wait until finish creating snapshot
        logging.debug("Starting waiting snapshot for node %s..." % task['Node']['Name'])
        finish = False
        while not finish:
            rst = self.do_safely(method=self.vultr.snapshot.list,
                                 description="waiting snapshot for node %s" % task['Node']['Name'],
                                 args={
                                     'params': {'SNAPSHOTID': self.tmpSnapshotID}
                                 })
            if isinstance(rst, bool) and not rst:
                return False
            logging.debug("Snapshot status is %s" % rst[self.tmpSnapshotID]['status'])
            finish = True if rst[self.tmpSnapshotID]['status'] == 'complete' else False
            time.sleep(self.queryInterval)
        logging.debug("Successfully create snapshot for node %s" % task['Node']['Name'])
        return True

    def destroy_tmp_snapshot_vultr(self, task):
        # Start create snapshot
        logging.debug("Starting destroying snapshot for node %s..." % task['Node']['Name'])
        rst = self.do_safely(method=self.vultr.snapshot.destroy,
                             description="destroying snapshot for node %s" % task['Node']['Name'],
                             args={
                                 'snapshotid': self.tmpSnapshotID
                             })
        if isinstance(rst, bool) and not rst:
            return False
        logging.debug("Successfully destroy snapshot for node %s" % task['Node']['Name'])
        return True

    def destroy_worker_vultr(self, task):
        # Copy targetVPS to prevent chaos
        targetVPS = copy.deepcopy(self.targetVPS)
        # Calc sleep time
        time_array = time.strptime(targetVPS['date_created'], '%Y-%m-%d %H:%M:%S')
        create_time = int(time.mktime(time_array)) + self.timeZone * 3600
        sleep_time = self.vultrDestroyCoolTime - (int(time.time()) - create_time)
        logging.debug('sleep_time: %s, now_time: %s, create_time: %s' % (
            sleep_time, int(time.time()), create_time))
        if sleep_time > 0:
            logging.debug("Cooling down %d seconds for destroying %s Vultr VPS %s created at %s for node %s..."
                          % (sleep_time, targetVPS['ram'],
                             targetVPS['SUBID'],
                             targetVPS['date_created'],
                             task['Node']['Name']))
            time.sleep(sleep_time)
        # Perform destroy
        logging.debug("Destroying %s Vultr VPS %s created at %s for node %s..." % (targetVPS['ram'],
                                                                                   targetVPS['SUBID'],
                                                                                   targetVPS['date_created'],
                                                                                   task['Node']['Name']))
        rst = self.do_safely(method=self.vultr.server.destroy,
                             description="destroying vps for node %s" % task['Node']['Name'],
                             args={
                                 'subid': targetVPS['SUBID'],
                                 'params': {}
                             })
        if isinstance(rst, bool) and not rst:
            return False
        # Wait for destroying
        logging.debug("Waiting for destroying complete for node %s..." % (task['Node']['Name']))
        finish = False
        while not finish:
            servers = self.do_safely(method=self.vultr.server.list,
                                     description="getting server info for node %s" % task['Node']['Name'],
                                     args={})
            if isinstance(servers, bool) and not servers:
                return False
            logging.debug("Server list: %s" % servers)
            finish = targetVPS['SUBID'] not in servers
            time.sleep(self.queryInterval)
        logging.debug("Destroy complete for node %s" % (task['Node']['Name']))
        logging.debug("Waiting for %d seconds to prevent false positive..." % self.destroyWaitTime)
        time.sleep(self.destroyWaitTime)  # Wait for a while to prevent false positive
        return True

    def destroy_vultr(self, task, wait=False):
        logging.debug("Starting destroying %s Vultr VPS %s created at %s for node %s..." % (self.targetVPS['ram'],
                                                                                            self.targetVPS['SUBID'],
                                                                                            self.targetVPS[
                                                                                                'date_created'],
                                                                                            task['Node']['Name']))
        rst = self.destroyPool.apply_async(func=self.destroy_worker_vultr, args=(task,))
        self.destroyResults.append(rst)
        if wait:
            logging.debug("Waiting for destroying worker return for node %s..." % (task['Node']['Name']))
            rst.wait()
            return rst.successful() and rst.get()
        logging.debug("Wait = False, returning...")
        return True

    def create_vultr(self, task):
        if task['Node']['Plan'] == self.VULTR512MPLANID:
            logging.debug("Creating 512M Vultr VPS using snapshot %s for node %s..." % (self.tmpSnapshotID,
                                                                                        task['Node']['Name']))
        elif task['Node']['Plan'] == self.VULTR1024MPLANID:
            logging.debug("Creating 1024M Vultr VPS using snapshot %s for node %s..." % (self.tmpSnapshotID,
                                                                                         task['Node']['Name']))
        rst = self.do_safely(method=self.vultr.server.create,
                             description="creating vps for node %s" % task['Node']['Name'],
                             args={
                                 'dcid': task['Node']['DataCenter'],
                                 'vpsplanid': task['Node']['Plan'],
                                 'osid': 164,
                                 'params': {
                                     'SNAPSHOTID': self.tmpSnapshotID,
                                     'enable_ipv6': 'yes',
                                     'label': task['Node']['Name'] + '-' + str(int(time.time()))
                                 }
                             })
        fall_back = False
        new_vps_id = ''
        if isinstance(rst, bool) and not rst:
            err = self.lastException
            if 'plan' in str(err) and 'unavailable' in str(err):
                logging.warning("Vultr VPS create failed! Falling back to create 1024M" +
                                " for node %s..." % task['Node']['Name'])
                fall_back = True
            else:
                logging.warning("Vultr VPS create failed due to unknown error, returning...")
                return False
        else:
            new_vps_id = rst['SUBID']
        if fall_back:
            rst = self.do_safely(method=self.vultr.server.create,
                                 description="creating vps for node %s" % task['Node']['Name'],
                                 args={
                                     'dcid': task['Node']['DataCenter'],
                                     'vpsplanid': self.VULTR1024MPLANID,
                                     'osid': 164,
                                     'params': {
                                         'SNAPSHOTID': self.tmpSnapshotID,
                                         'enable_ipv6': 'yes',
                                         'label': task['Node']['Name'] + str(int(time.time()))
                                     }
                                 })
            if isinstance(rst, bool) and not rst:
                return False
            new_vps_id = rst['SUBID']
        # Clear target VPS
        self.targetVPS = None
        logging.debug("New VPS ID is %s" % new_vps_id)
        # Wait for creating and update target VPS info
        logging.debug("Waiting for creating complete for node %s..." % (task['Node']['Name']))
        finish = False
        while not finish:
            if not self.get_server_info_vultr(task=task, subid=new_vps_id):
                return False
            logging.debug("Server info: %s" % self.targetVPS)
            time_array = time.strptime(self.targetVPS['date_created'], '%Y-%m-%d %H:%M:%S')
            create_time = int(time.mktime(time_array)) + self.timeZone * 3600
            lifespan = int(time.time()) - create_time
            logging.debug('Lifespan of the new server: %s' % lifespan)
            finish = lifespan > self.createMinTime and \
                     self.targetVPS['status'] == 'active' and self.targetVPS['power_status'] == 'running' and \
                     self.targetVPS['server_state'] == 'ok'
            time.sleep(self.queryInterval)
        logging.debug("Successfully create new VPS for node %s" % task['Node']['Name'])
        return True

    def destroy_and_create_vultr(self, task):
        # Check availability for 1024M Vultr VPS plan in given DataCenter
        rst = self.do_safely(method=self.vultr.regions.list,
                             description="getting regions availability list for node %s" % task['Node']['Name'],
                             args={
                                 'params': {
                                     'availability': 'yes'
                                 }
                             })
        if isinstance(rst, bool) and not rst:
            return False
        logging.debug("Regions availability: %s" % rst)
        if task['Node']['DataCenter'] not in rst or int(self.VULTR1024MPLANID) not in \
                rst[task['Node']['DataCenter']]['availability']:
            logging.error(
                "Even 1024M Vultr VPS is not available in DataCenter %s, returning..." % task['Node']['DataCenter'])
            return False
        # Update old VPS list
        self.oldVPSList.append(self.targetVPS)
        # Destroy useless VPS
        if self.targetVPS['VPSPLANID'] == self.VULTR512MPLANID:
            if not self.destroy_vultr(task, wait=True):
                return False
        elif self.targetVPS['VPSPLANID'] == self.VULTR1024MPLANID:
            if not self.destroy_vultr(task, wait=False):
                return False
        # Create new VPS using temporary snapshot
        if not self.create_vultr(task):
            return False
        return True
