# !/usr/bin/env python
# -*- coding:utf-8 -*-

import threading
import jobcontroller
import time

th = threading.Thread(target=jobcontroller.main_loop)
th.start()
try:
    while True:
        time.sleep(1)
except KeyboardInterrupt:
    print("Main: Interrupt signal received, exiting")
    jobcontroller.loop = False
    th.join()
