# !/usr/bin/env python
# -*- coding:utf-8 -*-

import os
import shutil
import subprocess
import config
import httphelper
import threading
import time
import traceback

tmpPath = config.TMP_PATH
waitTime = config.LOG_PERIOD


class TestController:

    sync = True

    def sync_log(self, id, logfile):
        log_f = open(logfile, "r")
        while self.sync:
            try:
                log = log_f.read()
            except UnicodeDecodeError:
                log = traceback.format_exc()
            httphelper.sync_log(id, log)
            time.sleep(waitTime)
        try:
            log = log_f.read()
        except UnicodeDecodeError:
            log = traceback.format_exc()
        httphelper.sync_log(id, log)
        httphelper.sync_log(id, "Program exited")
        httphelper.call_judge(id)
        log_f.close()

    def run_test(self, id, json, docker, callback):
        self.make_tmp()
        log_f = open(tmpPath + '/log/' + str(id) + ".log", "w")
        th = threading.Thread(target=self.sync_log, kwargs={'id': id, 'logfile': tmpPath + '/log/' + str(id) + ".log"})
        try:
            th.start()
            id = str(id)
            testPath = tmpPath + '/' + id
            if os.path.exists(testPath):
                print("Removing directory " + testPath)
                shutil.rmtree(testPath)
            print("Creating directory " + testPath)
            os.makedirs(testPath)

            print("Start Collecting test file")
            print("test-your-ss.sh")
            shutil.copyfile("test-your-ss.sh", testPath + '/' + "test-your-ss.sh")
            print("Dockerfile")
            with open(testPath + '/' + "Dockerfile", "w") as f:
                f.write("FROM " + docker + "\n")
                f.write("COPY . /\n")
                f.write("CMD [\"/bin/sh\", \"test-your-ss.sh\"]")
            print("shadowsocks.json")
            with open(testPath + '/' + "shadowsocks.json", "w") as f:
                f.write(json)

            docker_img = "test-your-ss-img" + id
            docker_con = "test-your-ss-con" + id
            print("Start building docker " + docker_img)

            subprocess.call("docker build -t " + docker_img + " .", shell=True, cwd=testPath, stdout=log_f, stderr=log_f)
            print("Start running docker " + docker_con)
            subprocess.call("docker run --name=\"" + docker_con + "\" " + docker_img, shell=True, cwd=testPath, stdout=log_f, stderr=log_f)
            print("Stopping docker container")
            subprocess.call("docker stop " + docker_con, shell=True, cwd=testPath, stdout=log_f, stderr=log_f)
            print("Cleaning up")
            print("Removing docker container")
            subprocess.call("docker rm -f " + docker_con, shell=True, cwd=testPath, stdout=log_f, stderr=log_f)
            print("Removing docker image")
            subprocess.call("docker image rm -f " + docker_img, shell=True, cwd=testPath, stdout=log_f, stderr=log_f)
            print("Removing temp directory")
            shutil.rmtree(testPath)
        except:
            log_f.write(traceback.format_exc())
        log_f.close()

        self.sync = False
        th.join()
        callback(threading.current_thread())

    def make_tmp(self):
        if not os.path.exists(tmpPath):
            print("Creating directory " + tmpPath)
            os.makedirs(tmpPath)
        if not os.path.exists(tmpPath + '/log'):
           print("Creating directory " + tmpPath + '/log')
           os.makedirs(tmpPath + '/log')
        return

# # test
# with open("/home/lijiahao/ss-cfg/jp1.ss.json", "r") as f:
#     json = f.read()
# run_test(1, json, "cool2645/shadowsocksr-master")
