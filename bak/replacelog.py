#!/usr/bin/env python
# encoding:utf-8


import os 
import sys 

services = ["gateway", "user_service", "wallet_service", "ws_service", "config_service", "kline_service", "price_service"]


def alter(file,old_str,new_str):
    file_data = ""
    with open(file, "r") as f:
        for line in f:
            if old_str in line:
                line = line.replace(old_str,new_str)
            file_data += line
    with open(file,"w") as f:
        f.write(file_data)

def repfilename(filename):
    repstr = '''. "digicon/currency_service/log"'''
    new_str = '''log "github.com/sirupsen/logrus"'''
    alter(filename, repstr, new_str)
    repstr = '''"digicon/currency_service/log"'''
    alter(filename, repstr, new_str)

    repstr = '''. "digicon/ws_service/log"'''
    alter(filename,repstr, new_str)

    for service in services:
        repstr = '''. "digicon/'''+service+'''/log"'''
        alter(filename, repstr, new_str)

    repstr = '''log.Log.'''
    new_lrstr = '''log.'''
    alter(filename, repstr, new_lrstr)
    repstr = '''Log.'''
    new_lrstr = '''log.'''
    alter(filename, repstr, new_lrstr)

    
def mywalk(root):
    for (root, dirs, files) in os.walk(root):
        for filename in files:
            rfilename = os.path.join(root, filename)
            print (rfilename)
            repfilename(rfilename)


if __name__ == "__main__":
    if len(sys.argv) >= 2:
        path = sys.argv[1]
    else:
        path = "./"
    mywalk(path)

