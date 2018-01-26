#!/usr/bin/python
# -*- coding: utf8 -*-
import sys
import traceback
import logging
import re
import json
import threading
from http import HTTPStatus

sys.path.append("..")
from worker import Worker

loop = True


class Watcher(Worker):
    def __init__(self):
        super().__init__(__file__)
        self.type = "watcher"
        self.pingApiUrl = "https://www.ipip.net"
        self.parseRE = re.compile(r"parent\.call_ping\((.*?)\);")
        self.lostAllPercentThreshold = float(self.config.get('watcher', 'lostAllPercentThreshold'))

    def ping(self, host=None):
        if host is None:
            raise Exception("Try to ping an empty host")
        bak_api_url = self.apiUrl
        self.apiUrl = self.pingApiUrl
        ping_data_dict = {'a': 'send', 'host': host, 'area[]': 'china'}
        try:
            raw_rst = self._GET(path="ping.php", data_dict=ping_data_dict, isDeserialize=False)
            # logging.info("Raw result for pinging host %s:" % host)
            # logging.info(str(raw_rst, 'utf8'))
        except:
            logging.error("Failed while pinging host %s" % host)
            traceback.print_exc(file=sys.stderr)
            self.apiUrl = bak_api_url
        try:
            rst = self.parse_ping_rst(str(raw_rst, 'utf8'))
        except:
            logging.error("Failed while parsing raw result")
            traceback.print_exc(file=sys.stderr)
            self.apiUrl = bak_api_url
        try:
            data = self.cal_node_status(rst)
        except:
            logging.error("Failed while calculating result")
            traceback.print_exc(file=sys.stderr)
            self.apiUrl = bak_api_url
            return "", False
        data['data'] = rst
        is_passing = True if data['lost_all_percent'] < self.lostAllPercentThreshold else False
        try:
            self.apiUrl = bak_api_url
            return json.dumps(data), is_passing
        except:
            logging.error("Failed while dumping json")
            traceback.print_exc(file=sys.stderr)
            self.apiUrl = bak_api_url
            return "", False

    def parse_ping_rst(self, raw_rst):
        rst_list = self.parseRE.findall(raw_rst)
        rst_list_de = []
        for r in rst_list:
            rst_list_de.append(json.loads(r))
        return rst_list_de

    def cal_node_status(self, rst):
        n = len(rst)
        lost_all_cnt = 0
        lost_cnt = 0
        res_time_cnt = 0.0
        for r in rst:
            if int(r['loss']) == 100:
                lost_all_cnt += 1
            res_time_cnt += float(r['rtt_avg'])
            lost_cnt += int(r['loss'])
        lost_all_percent = float(lost_all_cnt / n) * 100
        avg_lost_percent = float(lost_cnt / n)
        avg_res_time = float(res_time_cnt / n)
        rst_dict = {'lost_all_percent': lost_all_percent,
                    'avg_lost_percent': avg_lost_percent,
                    'avg_res_time': avg_res_time,
                    'pinger_cnt': n
                    }
        return rst_dict

    def heartbeat(self):
        logging.debug("Sending heartbeat package...")
        try:
            rst = self._POST(path="status/worker/" + self.name, data_dict={'class': self.type})
            if not (rst['code'] == HTTPStatus.OK and rst['result']):
                logging.error('Heartbeat result: %s' % rst)
                raise Exception('Invalid heartbeat result')
        except:
            logging.error('Failed while heartbeating')
            traceback.print_exc(file=sys.stderr)

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

    def task_loop(self):
        while loop:
            self.heartbeat()
            tasks = self.get_tasks()
            if len(tasks) > 0:
                c = 0
                while len(self.threads) < self.maxThreads and c < len(tasks):
                    th = threading.Thread()

    def watch(self, task):
        # Try to assign one task
        try:
            rst = self._PUT(path='task/' + str(task['ID']) + '/assign', data_dict={'worker': self.name})
            if rst.code != HTTPStatus.OK:
                logging.error('Task assign result: %s' % rst)
                raise Exception('Invalid task assign result')
            elif rst['result']:
                try:
                    rst = self._PUT(path='task/' + str(task['ID']), data_dict={'worker': self.name,
                                                                               'state': 'Pending'})
                    if not (rst.code == HTTPStatus.OK and rst['result']):
                        logging.error('Update task status to Pending result: %s' % rst)
                        raise Exception('Update task status to Pending failed')
                except:
                    logging.error('Failed while updating task %s status to Pending' % task['ID'])
                    traceback.print_exc(file=sys.stderr)
            elif not rst['result']:
                return False
        except:
            logging.error('Failed while assigning task %s' % task['ID'])
            traceback.print_exc(file=sys.stderr)
        host = task['Node']['DomainPrefix4'] + '.' + task['Node']['DomainRoot']
        logging.info('Start pinging %s' % host)
        # Perform ping
        ping_rst, is_passing = self.ping(host)
        state = 'Passing' if is_passing else 'Failing'
        # Update task status and log
        try:
            rst = self._PUT(path='task/' + str(task['ID']), data_dict={'worker': self.name,
                                                                       'state': state,
                                                                       'log': ping_rst})
            if not (rst.code == HTTPStatus.OK and rst['result']):
                logging.error('Update task result: %s' % rst)
                raise Exception('Failed to update task result')
        except:
            logging.error('Failed to update task %s result' % task['ID'])
            traceback.print_exc(file=sys.stderr)
            return False
        return True

# w = Watcher()
# print(w.get_tasks())
# print(w.ping(host='us0.ss.2645net.work'))
