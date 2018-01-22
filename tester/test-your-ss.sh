#!/bin/bash
#########################################################################
# File Name: test-your-ss.sh
# Author: LI JIAHAO
# ###############
# mail: lijiahao@cool2645.com
# Created Time: Wed 19 Apr 2017 01:46:21 PM CST
#########################################################################

sslocal -c /shadowsocks.json &
sleep 10
proxychains /usr/bin/curl http://ipv4.vm0.test-ipv6.com/ip/
proxychains /usr/bin/curl http://ipv6.vm0.test-ipv6.com/ip/
