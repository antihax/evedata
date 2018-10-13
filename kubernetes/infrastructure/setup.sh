curl -fsSL https://raw.githubusercontent.com/appscode/voyager/7.0.0/hack/deploy/voyager.sh \
    | bash -s -- --provider=baremetal --rbac --enable-analytics=false

kubectl apply -f evedata.yaml
kubectl apply -f monitoring.yaml
kubectl apply -f nsqd.yaml
kubectl apply -f ledis-redis.yaml
kubectl apply -f graylog.yaml

kubectl label nodes --all beta.kubernetes.io/fluentd-ds-ready=true
kubectl label nodes loadbalancer01 loadbalancer=voyager
kubectl label nodes loadbalancer02 loadbalancer=voyager

kubectl label nodes redis01 database=redis
kubectl label nodes database01 database=mysql
kubectl label nodes database02 database=mysql

kubectl label nodes database01 elasticsearch=elasticsearch
kubectl label nodes database02 elasticsearch=elasticsearch
kubectl label nodes database03 elasticsearch=elasticsearch

kubectl label nodes loadbalancer01 worker=worker
kubectl label nodes loadbalancer02 worker=worker
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
