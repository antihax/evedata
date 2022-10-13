#!/bin/bash
kubectl --namespace storage port-forward service/sql-mariadb 3306:3306 &
kubectl --namespace storage port-forward service/redis 6379:6379  &
kubectl --namespace nsq port-forward service/nsqadmin 4171:4171  &
kubectl --namespace nsq port-forward service/nsqlookupd1 4160:4160  &
kubectl --namespace nsq port-forward service/nsqlookupd1 4161:4161  &
kubectl --namespace nsq port-forward service/nsqd 4151:4151  &
kubectl --namespace nsq port-forward service/nsqd 4150:4150  &