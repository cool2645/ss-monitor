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
from vultr import Vultr

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
        self.createWaitTime = int(self.config.get('cleaner', 'createWaitTime'))
        self.destroyWaitTime = int(self.config.get('cleaner', 'destroyWaitTime'))
        self.oldVPSList = []
        # Init Vultr api instance
        self.vultrApikey = self.config.get('cleaner', 'vultrApikey')
        self.vultr = Vultr(self.vultrApikey)
        # Function dic for different VPS providers
        self.create_tmp_snapshot = {
            'Vultr': self.create_tmp_snapshot_vultr
        }
        self.get_server_info = {
            'Vultr': self.get_server_info_vultr
        }
        self.destroy_and_create = {
            'Vultr': self.destroy_and_create_vultr
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
        pass

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
        except:
            traceback.print_exc(file=sys.stderr)
            if depth < self.maxTry:
                logging.warning('Failed while %s after trying %d times, retrying...'
                                % (description, (depth + 1)))
                return self.do_safely(method, description, args, depth + 1)
            else:
                logging.error('Failed while %s after trying %d times, exiting...'
                              % (description, self.maxTry))
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

    def destroy_and_create_vultr(self, task):
        # Destroy useless VPS
        logging.debug("Destroying %s Vultr VPS %s created at %s for node %s..." % (self.targetVPS['ram'],
                                                                                   self.targetVPS['SUBID'],
                                                                                   self.targetVPS['date_created'],
                                                                                   task['Node']['Name']))
        rst = self.do_safely(method=self.vultr.server.destroy,
                             description="destroying vps for node %s" % task['Node']['Name'],
                             args={
                                 'subid': self.targetVPS['SUBID'],
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
            finish = self.targetVPS['SUBID'] not in servers
            time.sleep(self.queryInterval)
        logging.debug("Destroy complete for node %s" % (task['Node']['Name']))
        logging.debug("Waiting for %d seconds to prevent false positive..." % self.destroyWaitTime)
        time.sleep(self.destroyWaitTime)  # Wait for a while to prevent false positive
        # Create new VPS using temporary snapshot
        logging.debug("Creating %s Vultr VPS using snapshot %s for node %s..." % (self.targetVPS['ram'],
                                                                                  self.tmpSnapshotID,
                                                                                  task['Node']['Name']))
        rst = self.do_safely(method=self.vultr.server.create,
                             description="creating vps for node %s" % task['Node']['Name'],
                             args={
                                 'dcid': self.targetVPS['DCID'],
                                 'vpsplanid': self.targetVPS['VPSPLANID'],
                                 'osid': 164,
                                 'params': {
                                     'SNAPSHOTID': self.tmpSnapshotID,
                                     'enable_ipv6': 'yes',
                                     'label': task['Node']['Name']
                                 }
                             })
        if isinstance(rst, bool) and not rst:
            return False
        # Update old VPS list and target VPS
        self.oldVPSList.append(self.targetVPS)
        self.targetVPS = None
        new_id = rst['SUBID']
        logging.debug("New VPS ID is %s" % new_id)
        # Wait for creating
        logging.debug("Waiting for %d seconds to prevent false positive..." % self.createWaitTime)
        time.sleep(self.createWaitTime)  # Wait for a while to prevent false positive
        logging.debug("Waiting for creating complete for node %s..." % (task['Node']['Name']))
        finish = False
        while not finish:
            if not self.get_server_info_vultr(task=task, subid=new_id):
                return False
            logging.debug("Server info: %s" % self.targetVPS)
            finish = self.targetVPS['status'] == 'active' and self.targetVPS['power_status'] == 'running' and \
                     self.targetVPS['server_state'] == 'ok'
            time.sleep(self.queryInterval)
        logging.debug("Successfully create new VPS for node %s" % task['Node']['Name'])
        return True
