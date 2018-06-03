#!/bin/bash
#########################################################################
# File Name: 1.sh
# Author: LI JIAHAO
# ###############
# Mail: lijiahao@cool2645.com
# Created Time: Sun 03 Jun 2018 11:47:53 PM CST
#########################################################################
filename="$1"
while read -r line
do
    name="$line"
    for i in $(dig +short $name | grep -oE "\b([0-9]{1,3}\.){3}[0-9]{1,3}\b$"); do
        echo "Digging $name..."
        echo -e "$i\t$name" >> hosts
    done
done < "$filename"
