# !/usr/bin/env python
# -*- coding:utf-8 -*-

import requests
import config
import sys
import traceback

apiKey = config.API_KEY
apiUrl = config.API_URL
tryMax = config.TRY_MAX
runHost = config.HOST_NAME


def get_queuing_job_list():
    try:
        r = requests.get(apiUrl + "/jobs/queue")
        queue = r.json()
    except:
        traceback.print_exc(file=sys.stderr)
        return []
    else:
        return queue


def assign_job(queue, count=0):
    if len(queue) > 0:
        payload = {'key': apiKey, 'run_host': runHost}
        try:
            r = requests.post(apiUrl + "/jobs/" + str(queue[0]['id']), params=payload)
            res = r.json()
        except:
            traceback.print_exc(file=sys.stderr)
            res = {'result': False, 'msg': "Network Exception"}
        if res['result']:
            return {'id': res['id'], 'json': res['json'], 'docker': res['docker']}
        else:
            print(res['msg'])
            if count < tryMax:
                print("Retrying")
                return assign_job(get_queuing_job_list(), count + 1)
            else:
                print("Failed too many times, giving up")
                return False
    else:
        return False


def sync_log(id, log, count=0):
    payload = {'key': apiKey, 'log':  log}
    try:
        r = requests.post(apiUrl + "/jobs/" + str(id) + "/log", params=payload)
        res = r.json()
    except:
        traceback.print_exc(file=sys.stderr)
        res = {'result': False, 'msg': "Network Exception"}
    if res['result']:
        return True
    else:
        print(res['msg'])
        if count < tryMax:
            print("Retrying")
            return sync_log(id, log, count + 1)
        else:
            print("Failed too many times, giving up")
            return False


def cancel_job(id, count=0):
    payload = {'key': apiKey, '_method': 'DELETE'}
    try:
        r = requests.post(apiUrl + "/jobs/" + str(id), params=payload)
        res = r.json()
    except:
        traceback.print_exc(file=sys.stderr)
        res = {'result': False, 'msg': "Network Exception"}
    if res['result']:
        return True
    else:
        print(res['msg'])
        if count < tryMax:
            print("Retrying")
            return cancel_job(id, count + 1)
        else:
            print("Failed too many times, giving up")
            return False


def call_judge(id, count=0):
    payload = {'key': apiKey}
    try:
        r = requests.get(apiUrl + "/jobs/" + str(id) + "/judge", params=payload)
        res = r.json()
    except:
        traceback.print_exc(file=sys.stderr)
        res = {'result': False, 'msg': "Network Exception"}
    if res['result']:
        return True
    else:
        print(res['msg'])
        if count < tryMax:
            print("Retrying")
            return call_judge(id, count + 1)
        else:
            print("Failed too many times, giving up")
            return False
