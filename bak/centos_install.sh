#!/bin/sh


remote_ip=47.106.136.96

services="currency_service"

remote_path="/root/go/src/dig/"


function build_service(){
    for service in $services;
    do
        cd ../$service
        GOOS=linux GOARCH=amd64 go build
    done
}



function push_to_remote(){
    for service in $services;
    do
        scp ../$service/$service root@$remote_ip:$remote_path$service.nw
        result=`ssh root@$remote_ip "cd $remote_path && ls $service.nw"`
        echo "push result:$result"
    done
}

function restart_service(){
    for service in $services;
    do
        service_id=`ps -ef |grep -v grep | grep $service | awk -F " "  '{print $2}'`
        result=`ssh root@$remote_ip "ps -ef | grep -v grep | grep $service"`
        service_id=`echo $result | awk -F " " '{print $2}'`
        echo  "$service $service_id"
        #ssh root@$remote_ip  "kill -9 $service_id && cd $remote_path && mv -f $service.nw $service && nohup ./$service &"
        restart_result=`ssh root@$remote_ip "ps -ef | grep -v grep | grep $service"`
        service_id=`echo $restart_result | awk -F " " '{print $2}'`
        echo  "$service $service_id"        
    done
}

#build_service
push_to_remote
restart_service
