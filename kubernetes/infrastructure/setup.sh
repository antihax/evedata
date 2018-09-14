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