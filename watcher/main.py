#!/usr/bin/python
# -*- coding: utf8 -*-
from watcher import Watcher
from multiprocessing.dummy import Pool as ThreadPool
import threading
import time
import logging

w = Watcher()
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
        # print('ready: ', result.ready())  # 线程函数是否已经启动了
        # print('map_async: 不堵塞')
        # result.wait()  # 等待所有线程函数执行完毕
        # print('after wait')
        # if result.ready():  # 线程函数是否已经启动了
        #     if result.successful():  # 线程函数是否执行成功
        #         print('result: ', result.get())  # 线程函数返回值
        time.sleep(w.heartRate)


def do(task):
    watcher = Watcher()
    return watcher.watch(task)


main()

# w = Watcher()
# w.heartbeat()
# tasks = w.get_tasks()
# t = tasks[0]
# print(t)
# print(w.watch(t))
