# !/usr/bin/env python
# -*- coding:utf-8 -*-

import httphelper
import testcontroller
import config
import threading
import time

maxThreads = config.MAX_THREADS
waitTime = config.POLLING_PERIOD

loop = True
threads = []


def thread_call_back(th):
    threads.remove(th)


def main_loop():
    global threads

    while loop:
        print("Polling queue")
        queue = httphelper.get_queuing_job_list()
        print(time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time())) + ": There are " + str(len(queue)) + " jobs in queue")
        print(time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time())) + ": Currently running " + str(len(threads)) + " threads")
        if len(queue) > 0 and len(threads) < maxThreads:
            job = httphelper.assign_job(queue)
            if job:
                print("Starting thread for test " + str(job['id']))
                tc = testcontroller.TestController()
                th = threading.Thread(target=tc.run_test,
                                      kwargs={'id': job['id'], 'json': job['json'], 'docker': job['docker'], 'callback': thread_call_back})
                threads.append(th)
                th.start()
        time.sleep(waitTime)

    print("JobController: Interrupt signal received, proceeding graceful exit")
    while len(threads) > 0:
        threads[0].join()
