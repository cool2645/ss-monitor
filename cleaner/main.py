#!/usr/bin/python
# -*- coding: utf8 -*-

from cleaner import Cleaner
from multiprocessing.dummy import Pool as ThreadPool
import threading
import time
import sys
import logging
import traceback

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
        if result is not None:
            result.wait()
        return


def main_loop():
    global loop
    global result
    while loop:
        w.heartbeat()
        tasks = w.get_tasks()
        logging.debug("Task list: %s" % tasks)
        result = pool.map_async(do, tasks)
        time.sleep(w.heartRate)


def do(task):
    logging.debug("Start doing task: %s" % task)
    cleaner = Cleaner()
    try:
        return cleaner.clean(task)
    except:
        traceback.print_exc(file=sys.stderr)
        logging.critical('Failed while cleaning for task %s' % (task['ID']))
        return False


if __name__ == '__main__':
    main()
