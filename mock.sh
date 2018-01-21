#!/bin/bash

# Remove any currently running containers
docker stop mysql teamspeak mock-esi redis nsqlookup nsqadmin nsqd | xargs docker rm

# MySQL Server
docker run --net=host --name=mysql --health-cmd='mysqladmin ping --silent' -d -p 127.0.0.1:3306:3306 -h evedata.sql -e MYSQL_ALLOW_EMPTY_PASSWORD=true percona

# Teamspeak Server
docker run --net=host --name=teamspeak -d -p 127.0.0.1:9987:9987/udp -p 127.0.0.1:30033:30033 -p 127.0.0.1:10011:10011 -p 127.0.0.1:41144:41144 mbentley/teamspeak clear_database=1 serveradmin_password=nothinguseful

# Mock ESI Server
docker run --net=host --name=mock-esi -d  -h mock-esi -p 127.0.0.1:8080:8080 antihax/mock-esi

# Redis 
docker run --net=host --name=redis -d -p 127.0.0.1:6379:6379 -h evedata.sql redis

# NSQ Service
docker run --net=host --name=nsqlookup -d -p 127.0.0.1:4160:4160 -p 4161:4161 -h nsqlookupd1.nsq nsqio/nsq /nsqlookupd
docker run --net=host --name=nsqadmin -d -p 127.0.0.1:4171:4171 -h nsqadmin.nsq nsqio/nsq /nsqadmin --lookupd-http-address=127.0.0.1:4161 -max-msg-size=4194304
docker run --net=host --name=nsqd -d -p 127.0.0.1:4151:4151 -p 4150:4150 -h localhost nsqio/nsq /nsqd --lookupd-tcp-address=127.0.0.1:4160

# Get the admin token for the TS server
until [ `docker inspect -f "{{.State.Status}}" teamspeak | grep -c running` -eq 1 ]
do
    sleep 1
done
docker exec teamspeak /bin/bash -c 'cat /data/logs/*.log | grep -oE "token=([A-Za-z0-9+_=]+)" | cut -d= -f2' > teamspeakToken.txt

# Populate SQL
until [ `docker inspect -f "{{json .State.Health.Status }}" mysql | grep -c healthy` -eq 1  ]
do
    sleep 1
done
echo "create database eve; create database evedata; set sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO';" | docker exec -i mysql /bin/bash -c mysql
cat ./services/vanguard/sql/evedata.sql | docker exec -i mysql /bin/bash -c 'mysql -Devedata'
unzip -p ./services/vanguard/sql/eve.zip | docker exec -i mysql /bin/bash -c 'mysql -Deve'
