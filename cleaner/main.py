#!/usr/bin/python
# -*- coding: utf8 -*-
from cleaner import Cleaner
from multiprocessing.dummy import Pool as ThreadPool
import threading
import time
import logging


def test():
    test_task = {
        'Node': {
            'Name': 'Fake-US3',
            'IPv4': '45.77.148.98',
            'IPv6': ''
        }
    }
    c = Cleaner()
    c.tmpSnapshotID = '9e65a6e32c09f'
    print(c.get_server_info['Vultr'](task=test_task))
    # print(c.create_tmp_snapshot['Vultr'](task=test_task))
    # print(c.destroy_and_create['Vultr'](task=test_task))
    # print("oldVPSList: \n", c.oldVPSList)
    task_dict = {
        'callback_id': 0,
        'class': 'watcher',
        'node_id': 0,
        'ip_ver': 4,
        'server_name': '45.77.148.98'
    }
    print(c.assign_task_and_wait(task=test_task, new_task_dict=task_dict))
