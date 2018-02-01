#!/usr/bin/python
# -*- coding: utf8 -*-
from cleaner import Cleaner
from multiprocessing.dummy import Pool as ThreadPool
import threading
import time
import logging

w = Cleaner()
pool = ThreadPool(w.maxThreads)
loop = True
result = None


def main():
    global loop
    global result
    th = threading.Thread(target=main_loop)
    th.start()
    try:
        while True:
            time.sleep(1)
            if not th.is_alive():
                logging.critical("Main loop thread is dead, exiting")
                return
    except KeyboardInterrupt:
        logging.warning("Interrupt signal received, exiting")
        loop = False
        result.wait()
        return


def main_loop():
    global loop
    global result
    while loop:
        w.heartbeat()
        tasks = w.get_tasks()
        result = pool.map_async(do, tasks)
        time.sleep(w.heartRate)


def do(task):
    cleaner = Cleaner()
    return cleaner.clean(task)


if __name__ == '__main__':
    main()

# def test():
# test_task = {
#     'ID': 1,
#     'Node': {
#         'ID': 1,
#         'Name': 'Fake-US3',
#         'IPv4': '207.148.19.134',
#         'IPv6': '',
#         'DataCenter': '1',
#         'Plan': '200',
#         'Snapshot': '',
#         'Provider': 'Vultr',
#         'DNSProvider': 'DNSimple',
#         'DomainPrefix4': 'fake-us3.ss',
#         'DomainPrefix6': 'fake-us3.ss6',
#         'DomainRoot': '2645net.work'
#     }
# }
# test_vps_vultr = {
#     'main_ip': '123.123.123.123',
#     'v6_networks': [
#         {'v6_main_ip': '1234:1234:1234:1234:1234:1234:1234:1234'}
#     ]
# }
# c = Cleaner()

# c.targetVPS = test_vps_vultr

# print(c.update_dns['DNSimple'](test_task))

# print(c.update_node_info(test_task))

# print(c.broadcast('🔶🔵🔴 做女孩子吧！'))

# print(c.clean(test_task))

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
