#!/usr/bin/python
# -*- coding: utf8 -*-
from cleaner import Cleaner
from multiprocessing.dummy import Pool as ThreadPool
import threading
import time
import logging


def test():
    test_task = {
        'ID': 1,
        'Node': {
            'Name': 'Fake-US3',
            'IPv4': '207.148.19.134',
            'IPv6': '',
            'DataCenter': '1',
            'Plan': '200',
            'Snapshot': '',
            'Provider': 'Vultr'
        }
    }
    c = Cleaner()

    # print(c.broadcast('🔶🔵🔴 做女孩子吧！'))

    print(c.clean(test_task))

    # c.tmpSnapshotID = test_task['Node']['Snapshot']

    # print(c.get_server_info['Vultr'](task=test_task))

    # print(c.create_tmp_snapshot['Vultr'](task=test_task))

    # print(c.destroy_and_create['Vultr'](task=test_task))
    #
    # print('Waiting for destroying Threads...')
    # for r in c.destroyResults:
    #     r.wait()
    # for r in c.destroyResults:
    #     if r.ready() and r.successful():
    #         print(r.get())

    # print("oldVPSList: \n", c.oldVPSList)
    # task_dict = {
    #     'callback_id': 0,
    #     'class': 'watcher',
    #     'node_id': 0,
    #     'ip_ver': 4,
    #     'server_name': '45.77.148.98'
    # }
    # print(c.assign_task_and_wait(task=test_task, new_task_dict=task_dict))


test()
