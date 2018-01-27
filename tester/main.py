#!/usr/bin/python
# -*- coding: utf8 -*-
from tester import Tester
from multiprocessing.dummy import Pool as ThreadPool
import threading
import time
import logging

w = Tester()
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
    tester = Tester()
    return tester.test(task)


if __name__ == '__main__':
    main()
