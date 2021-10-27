kubectl apply -f evedata.yaml
kubectl apply -f monitoring.yaml
kubectl apply -f nsqd.yaml
kubectl apply -f redis.yaml
#kubectl apply -f graylog.yaml

kubectl label nodes --all beta.kubernetes.io/fluentd-ds-ready=true
kubectl label nodes loadbalancer01 loadbalancer=voyager
kubectl label nodes loadbalancer02 loadbalancer=voyager

kubectl label nodes database01 database=mysql
kubectl label nodes database02 database=mysql
kubectl label nodes database03 database=mysql

kubectl label nodes database01 elasticsearch=elasticsearch
kubectl label nodes database02 elasticsearch=elasticsearch
kubectl label nodes database03 elasticsearch=elasticsearch

kubectl label node loadbalancer01 worker=worker
kubectl label node loadbalancer02 worker=worker
kubectl label nodes worker01 worker=worker
kubectl label nodes worker02 worker=worker
kubectl label nodes worker03 worker=worker
kubectl label nodes worker04 worker=worker
kubectl label nodes worker05 worker=worker
kubectl label nodes worker06 worker=worker
kubectl label nodes worker07 worker=worker
kubectl label nodes worker08 worker=worker
kubectl label nodes worker09 worker=worker
kubectl label nodes worker10 worker=worker
kubectl label nodes worker11 worker=worker
kubectl label nodes worker12 worker=worker
kubectl label nodes worker13 worker=worker
kubectl label nodes worker14 worker=worker
kubectl label nodes worker15 worker=worker
kubectl label nodes worker16 worker=worker
kubectl label nodes worker17 worker=worker
kubectl label nodes worker18 worker=worker
kubectl label nodes worker19 worker=worker
kubectl label nodes worker20 worker=worker
kubectl label nodes worker21 worker=worker
kubectl label nodes worker22 worker=worker
kubectl label nodes worker23 worker=worker
kubectl label nodes worker24 worker=worker
kubectl label nodes worker25 worker=worker
kubectl label nodes worker26 worker=worker
kubectl label nodes worker27 worker=worker
kubectl label nodes worker28 worker=worker
kubectl label nodes worker29 worker=worker
