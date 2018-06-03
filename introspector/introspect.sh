#!/bin/bash
#########################################################################
# File Name: introspect.sh
# Author: LI JIAHAO
# ###############
# Mail: lijiahao@cool2645.com
# Created Time: Mon 04 Jun 2018 12:49:38 AM CST
#########################################################################
filename="beacon.txt"
timeout=10
count=0
crisis=30
while read -r line
do
    name="$line"
    curl --connect-timeout $timeout --output /dev/null $name
    code=$?
    echo "{\"host\":\"$name\",\"code\":$code}"
    if (( $code != 0 )); then
        count=$(($count + 1))
    fi
    if (( $count > $crisis )); then
        echo "{\"result\":$count,\"msg\":\"break due to too many failures\"}"
        exit
    fi
done < "$filename"
echo "{\"result\":$count,\"msg\":\"\"}"
