#!/usr/bin/python
# -*- coding: utf8 -*-
import sys
import traceback
import logging
import re
import json
import threading
import requests
from urllib.parse import urlencode
from http import HTTPStatus

sys.path.append("..")
from worker import Worker

loop = True


class Watcher(Worker):
    def __init__(self):
        super().__init__(__file__)
        self.type = "watcher"
        # These kawaii variables↓ only works for Threads in Threads Pool
        self.pingApiUrl = "https://www.ipip.net/ping.php"
        self.parseRE = re.compile(r"parent\.call_ping\((.*?)\);")
        self.lostAllPercentThreshold = float(self.config.get('watcher', 'lostAllPercentThreshold'))
        self.pinging = True
        self.toSync = False
        self.rawRst = ""
        self.logSyncingChunckSize = int(self.config.get('watcher', 'logSyncingChunckSize'))

    def ping(self, host=None, depth=0):
        if host is None:
            raise Exception("Try to ping an empty host")
        ping_data_dict = {'a': 'send', 'host': host, 'area[]': 'china'}
        try:
            r = requests.get(self.pingApiUrl + "?" + urlencode(ping_data_dict), stream=True, timeout=self.connTimeout)
            if r.encoding is None:
                r.encoding = 'utf-8'
            for chunck in r.iter_content(self.logSyncingChunckSize, decode_unicode=True):
                # filter out keep-alive new chuncks
                if chunck:
                    self.rawRst += chunck
                    self.toSync = True
        except:
            traceback.print_exc(file=sys.stderr)
            if depth < self.maxTry:
                logging.warning("Failed while pinging host %s after trying %d times, reconnecting..." % (host, depth))
                self.rawRst = ""
                return self.ping(host, depth + 1)
            else:
                logging.critical("Failed while pinging host %s after trying %d times, exiting..." % (host, self.maxTry))
                return False
        finally:
            self.pinging = False

    def sync_log(self, task):
        while self.pinging:
            if self.toSync:
                ping_rst = self.generate_report(self.rawRst)
                try:
                    ping_rst_str = json.dumps(ping_rst)
                except:
                    logging.error("Failed while dumping json")
                    traceback.print_exc(file=sys.stderr)
                    continue
                try:
                    logging.info('Updating task %s log' % task['ID'])
                    rst = self._PUT(path='task/' + str(task['ID']), data_dict={'worker': self.name,
                                                                               'log': ping_rst_str})
                    if not (rst.code == HTTPStatus.OK and rst['result']):
                        logging.error('Update task log: %s' % rst)
                        raise Exception('Failed to update task log')
                except:
                    logging.error('Failed to update task %s log' % task['ID'])
                    traceback.print_exc(file=sys.stderr)
                    return False
        ping_rst = self.generate_report(self.rawRst)
        state = 'Passing' if ping_rst['lost_all_percent'] < self.lostAllPercentThreshold else 'Failing'
        try:
            rst = self._PUT(path='task/' + str(task['ID']), data_dict={'worker': self.name,
                                                                       'state': state,
                                                                       'log': ping_rst})
            if not (rst.code == HTTPStatus.OK and rst['result']):
                logging.error('Update task result: %s' % rst)
                raise Exception('Failed to update task result')
        except:
            logging.error('Failed to update task %s result ' % task['ID'])
            traceback.print_exc(file=sys.stderr)
            return False

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
        if n != 0:
            lost_all_percent = float(lost_all_cnt / n) * 100
            avg_lost_percent = float(lost_cnt / n)
            avg_res_time = float(res_time_cnt / n)
        else:
            lost_all_percent = 0
            avg_lost_percent = 0
            avg_res_time = 0
        rst_dict = {'lost_all_percent': lost_all_percent,
                    'avg_lost_percent': avg_lost_percent,
                    'avg_res_time': avg_res_time,
                    'pinger_cnt': n
                    }
        return rst_dict

    def generate_report(self, raw_rst):
        try:
            rst = self.parse_ping_rst(raw_rst)
        except:
            logging.error("Failed while parsing raw result")
            traceback.print_exc(file=sys.stderr)
        try:
            data = self.cal_node_status(rst)
        except:
            logging.error("Failed while calculating result")
            traceback.print_exc(file=sys.stderr)
            return "", False
        data['data'] = rst
        return data

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

        th = threading.Thread(target=self.sync_log, kwargs={'task': task})
        logging.info('Start log syncing for ping %s' % host)
        self.pinging = True
        th.start()
        logging.info('Start pinging %s' % host)
        # Perform ping
        res = self.ping(host)
        th.join()
        return res