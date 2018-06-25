#!/bin/sh 


services="gateway  currency_service public_service token_service user_service wallet_service  kline_service "
for service in $services;do 
    service_id=`ps -ef |grep -v grep | grep $service | awk -F " "  '{print $2}'`
    echo "kill $service,$service_id"
    kill -9 $service_id
done
