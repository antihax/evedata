#!/bin/bash
set -e
docker run --privileged=true --rm -v /sys:/hostsys busybox sh -c "echo never > /hostsys/kernel/mm/transparent_hugepage/enabled"
docker run --privileged=true --rm -v /sys:/hostsys busybox sh -c "echo never > /hostsys/kernel/mm/transparent_hugepage/defrag"
set +e

# Remove any currently running containers
docker stop mysql mock-esi redis nsqlookup nsqadmin nsqd | xargs docker rm

# MySQL Server
docker run --net=host --name=mysql --health-cmd='mysqladmin ping --silent' -d -p 127.0.0.1:3306:3306 -e INIT_TOKUDB=1 -e MYSQL_ALLOW_EMPTY_PASSWORD=true percona/percona-server

# Mock ESI Server
docker run --net=host --name=mock-esi -d  -h mock-esi -p 127.0.0.1:8080:8080 antihax/mock-esi

# Redis 
docker run --net=host --name=redis -d -p 127.0.0.1:6379:6379  redis

# NSQ Service
docker run --net=host --name=nsqlookup -d -p 127.0.0.1:4160:4160 -p 4161:4161  nsqio/nsq /nsqlookupd
docker run --net=host --name=nsqadmin -d -p 127.0.0.1:4171:4171 nsqio/nsq /nsqadmin --lookupd-http-address=127.0.0.1:4161
docker run --net=host --name=nsqd -d -p 127.0.0.1:4151:4151 -p 4150:4150  nsqio/nsq /nsqd --lookupd-tcp-address=127.0.0.1:4160 -max-msg-size=8388608

# Populate SQL
until [ `docker inspect -f "{{json .State.Health.Status }}" mysql | grep -c healthy` -eq 1  ]
do
    sleep 1
    echo Percona not ready yet.
done
echo Percona Ready

echo "create database eve; create database evedata; set sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO';" | docker exec -i mysql /bin/bash -c 'mysql -uroot'
cat ./services/vanguard/sql/evedata.sql | docker exec -i mysql /bin/bash -c 'mysql -uroot -Devedata'
gzip -dc ./services/vanguard/sql/eve.gz | docker exec -i mysql /bin/bash -c 'mysql -uroot -Deve'
echo "SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));" | docker exec -i mysql /bin/bash -c 'mysql -uroot'
echo "SET GLOBAL sql_mode = 'NO_ENGINE_SUBSTITUTION';" | docker exec -i mysql /bin/bash -c 'mysql -uroot'

