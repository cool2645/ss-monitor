#!/usr/bin/python
# -*- coding: utf8 -*-
import sys
import traceback
import logging
import os
import threading
import shutil
import subprocess
import time

sys.path.append("..")
from worker import Worker
from httpstatus import HTTPStatus

loop = True


class Tester(Worker):
    def __init__(self):
        super().__init__(__file__)
        self.type = "tester"
        self.ipVer = int(self.config.get('tester', 'ipVer'))
        # These kawaii variables↓ only works for Threads in Threads Pool
        self.running = False
        self.tmpPath = self.config.get('tester', 'tmpPath')
        self.logSyncingInterval = int(self.config.get('tester', 'logSyncingInterval'))

    def run(self, id, json, log_f):
        id = str(id)
        test_path = self.tmpPath + '/' + id
        if os.path.exists(test_path):
            logging.info("Removing directory %s" % test_path)
            try:
                shutil.rmtree(test_path)
            except:
                logging.critical('Failed to remove directory %s' % test_path)
                traceback.print_exc(file=sys.stderr)
                return False

        logging.info("Creating directory %s" % test_path)
        try:
            os.makedirs(test_path)
        except:
            logging.critical('Failed to create directory %s' % test_path)
            traceback.print_exc(file=sys.stderr)
            return False

        logging.info("Start Collecting test file")
        try:
            self.copy_file("test-your-ss.sh", test_path + '/' + "test-your-ss.sh")
            self.copy_file("Dockerfile", test_path + '/' + "Dockerfile")
            self.copy_file("proxychains.conf", test_path + '/' + "proxychains.conf")
        except:
            return False

        logging.info("Write %s" % test_path + '/' + "shadowsocks.json")
        try:
            with open(test_path + '/' + "shadowsocks.json", "w") as f:
                f.write(json)
        except:
            logging.error('Failed to write %s' % test_path + '/' + "shadowsocks.json")
            traceback.print_exc(file=sys.stderr)
            return False

        try:
            rst = self._PUT(path='task/' + id, data_dict={'worker': self.name,
                                                                       'state': 'Running'})
            if not (rst.code == HTTPStatus.OK and rst['result']):
                logging.error('Update task status to Running result: %s' % rst)
                raise Exception('Update task status to Running failed')
        except:
            logging.error('Failed while updating task %s status to Running' % id)
            traceback.print_exc(file=sys.stderr)

        docker_img = "test-your-ss-img" + id
        docker_con = "test-your-ss-con" + id
        logging.info("Start building docker %s" % docker_img)
        try:
            subprocess.call("docker build -t " + docker_img + " .", shell=True, cwd=test_path, stdout=log_f, stderr=log_f)
        except:
            logging.error('Failed to build docker %s' % docker_img)
            traceback.print_exc(file=sys.stderr)

        logging.info("Start running docker %s" % docker_con)
        try:
            subprocess.call("docker run --name=\"" + docker_con + "\" " + docker_img, shell=True, cwd=test_path,
                        stdout=log_f, stderr=log_f)
        except:
            logging.error('Failed to run docker %s' % docker_con)
            traceback.print_exc(file=sys.stderr)

        logging.info("Stopping docker container %s" % docker_con)
        try:
            subprocess.call("docker stop " + docker_con, shell=True, cwd=test_path, stdout=log_f, stderr=log_f)
        except:
            logging.error('Failed to stop docker %s' % docker_con)
            traceback.print_exc(file=sys.stderr)

        logging.info("Cleaning up")
        logging.info("Removing docker container %s" % docker_con)
        try:
            subprocess.call("docker rm -f " + docker_con, shell=True, cwd=test_path, stdout=log_f, stderr=log_f)
        except:
            logging.error('Failed to remove docker %s' % docker_con)
            traceback.print_exc(file=sys.stderr)
        logging.info("Removing docker image %s" % docker_img)
        try:
            subprocess.call("docker image rm -f " + docker_img, shell=True, cwd=test_path, stdout=log_f, stderr=log_f)
        except:
            logging.error('Failed to remove docker image %s' % docker_img)
            traceback.print_exc(file=sys.stderr)

        logging.info("Removing temp directory %s" % test_path)
        try:
            shutil.rmtree(test_path)
        except:
            logging.error('Failed to remove temp directory %s' % test_path)
            traceback.print_exc(file=sys.stderr)

        finally:
            self.running = False

        return True

    def copy_file(self, ori, des):
        logging.info("Copy %s to %s" % (ori, des))
        try:
            shutil.copyfile(ori, des)
        except:
            logging.critical('Failed to copy %s to %s' % (ori, des))
            traceback.print_exc(file=sys.stderr)
            raise Exception('Failed to copy file')

    def sync_log(self, task, logfile):
        log, state = "", "Failing"
        while self.running:
            try:
                log_f = open(logfile, "r")
                log = log_f.read()
            except UnicodeDecodeError:
                log = traceback.format_exc()
            except:
                logging.error('Exception raised reading log %s' % logfile)
                traceback.print_exc(file=sys.stderr)
                continue

            try:
                logging.info('Updating task %s log' % task['ID'])
                rst = self._PUT(path='task/' + str(task['ID']), data_dict={'worker': self.name,
                                                                               'log': log})
                if not (rst.code == HTTPStatus.OK and rst['result']):
                    logging.error('Update task log: %s' % rst)
                    raise Exception('Failed to update task log')
            except:
                logging.error('Failed to update task %s log' % task['ID'])
                traceback.print_exc(file=sys.stderr)
                return False
            finally:
                log_f.close()
                time.sleep(self.logSyncingInterval)

        try:
            log_f = open(logfile, "r")
            log = log_f.read()
            state = 'Passing' if log.find('callback({') >= 0 else 'Failing'
        except UnicodeDecodeError:
            log = traceback.format_exc()
        except:
            logging.error('Exception raised reading log %s' % logfile)
            traceback.print_exc(file=sys.stderr)
        finally:
            log_f.close()

        try:
            logging.info('Updating task %s result' % task['ID'])
            rst = self._PUT(path='task/' + str(task['ID']), data_dict={'worker': self.name,
                                                                       'state': state,
                                                                       'log': log})
            if not (rst.code == HTTPStatus.OK and rst['result']):
                logging.error('Update task result: %s' % rst)
                raise Exception('Failed to update task result')
        except:
            logging.error('Failed to update task %s result ' % task['ID'])
            traceback.print_exc(file=sys.stderr)
            return False

    def heartbeat(self):
        logging.debug("Sending heartbeat package...")
        try:
            rst = self._POST(path="status/worker/" + self.name, data_dict={'class': self.type, 'ip_ver': self.ipVer})
            if not (rst['code'] == HTTPStatus.OK and rst['result']):
                logging.error('Heartbeat result: %s' % rst)
                raise Exception('Invalid heartbeat result')
        except:
            logging.error('Failed while heartbeating')
            traceback.print_exc(file=sys.stderr)

    def get_tasks(self):
        rst = None
        try:
            if self.ipVer == 4 or self.ipVer == 6:
                rst = self._GET(path="task", data_dict={'class': self.type, 'state': 'Queuing', 'ip_ver': self.ipVer})
            else:
                rst = self._GET(path="task", data_dict={'class': self.type, 'state': 'Queuing'})
            if not (rst['code'] == HTTPStatus.OK and rst['result']):
                logging.error('Get tasks result: %s' % rst)
                raise Exception('Invalid get tasks result')
        except:
            logging.error('Failed while getting tasks')
            traceback.print_exc(file=sys.stderr)
        return rst['data']

    def test(self, task):
        # Try to assign one task
        try:
            rst = self._PUT(path='task/' + str(task['ID']) + '/assign', data_dict={'worker': self.name})
            if rst.code != HTTPStatus.OK:
                logging.error('Task assign result: %s' % rst)
                raise Exception('Invalid task assign result')
            elif rst['result']:
                try:
                    rst = self._PUT(path='task/' + str(task['ID']), data_dict={'worker': self.name,
                                                                               'state': 'Building'})
                    if not (rst.code == HTTPStatus.OK and rst['result']):
                        logging.error('Update task status to Building result: %s' % rst)
                        raise Exception('Update task status to Building failed')
                except:
                    logging.error('Failed while updating task %s status to Building' % task['ID'])
                    traceback.print_exc(file=sys.stderr)
            elif not rst['result']:
                return False
        except:
            logging.error('Failed while assigning task %s' % task['ID'])
            traceback.print_exc(file=sys.stderr)

        host = task['Node']['Name'] + ' ' + ('IPv6' if task['IPVer'] == 6 else 'IPv4')
        id = task['ID']
        json = task['Node']['Ss6Json'] if task['IPVer'] == 6 else task['Node']['Ss4Json']

        try:
            self.make_tmp(self.tmpPath)
        except:
            logging.critical('Failed to create temp directory')
            traceback.print_exc(file=sys.stderr)
            return False

        logfile = self.tmpPath + '/log/' + str(id) + ".log"
        try:
            log_f = open(logfile, "w")
        except:
            logging.critical('Failed to open log file')
            traceback.print_exc(file=sys.stderr)
            return False

        th = threading.Thread(target=self.sync_log, kwargs={'task': task, 'logfile': logfile})
        logging.info('Start log syncing for run %d, %s' % (id, host))
        self.running = True
        th.start()
        logging.info('Start running %d, %s' % (id, host))
        # Perform run
        res = self.run(id, json, log_f)
        log_f.close()
        th.join()
        return res

    def make_tmp(self, tmp_path):
        if not os.path.exists(tmp_path):
            print("Creating directory " + tmp_path)
            os.makedirs(tmp_path)
        if not os.path.exists(tmp_path + '/log'):
           print("Creating directory " + tmp_path + '/log')
           os.makedirs(tmp_path + '/log')
        return