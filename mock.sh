#!/bin/bash
set -e

if test -f /sys/kernel/mm/transparent_hugepage/enabled; then
   echo "Need root permissions to disable /sys/kernel/mm/transparent_hugepage"
   echo "If this fails, please run this script with sudo"
   read -p "Press enter to continue or CTRL+C to cancel"
   echo never > /sys/kernel/mm/transparent_hugepage/enabled
fi
if test -f /sys/kernel/mm/transparent_hugepage/defrag; then
   echo "Need root permissions to disable /sys/kernel/mm/transparent_hugepage/defrag"
   echo "If this fails, please run this script with sudo"
   read -p "Press enter to continue or CTRL+C to cancel"
   echo never > /sys/kernel/mm/transparent_hugepage/defrag
fi

set +e

# Remove any currently running containers
docker stop mysql teamspeak mock-esi redis nsqlookup nsqadmin nsqd | xargs docker rm

# MySQL Server
docker run --net=host --name=mysql --health-cmd='mysqladmin ping --silent' -d -p 127.0.0.1:3306:3306 -e INIT_TOKUDB=1 -e MYSQL_ALLOW_EMPTY_PASSWORD=true percona/percona-server

# Teamspeak Server
docker run --net=host --name=teamspeak -d -p 127.0.0.1:9987:9987/udp -p 127.0.0.1:30033:30033 -p 127.0.0.1:10011:10011 -p 127.0.0.1:41144:41144 mbentley/teamspeak clear_database=1 license_accepted=1 serveradmin_password=nothinguseful

# Mock ESI Server
docker run --net=host --name=mock-esi -d  -h mock-esi -p 127.0.0.1:8080:8080 antihax/mock-esi

# Redis 
docker run --net=host --name=redis -d -p 127.0.0.1:6379:6379  redis

# NSQ Service
docker run --net=host --name=nsqlookup -d -p 127.0.0.1:4160:4160 -p 4161:4161  nsqio/nsq /nsqlookupd
docker run --net=host --name=nsqadmin -d -p 127.0.0.1:4171:4171 nsqio/nsq /nsqadmin --lookupd-http-address=127.0.0.1:4161
docker run --net=host --name=nsqd -d -p 127.0.0.1:4151:4151 -p 4150:4150  nsqio/nsq /nsqd --lookupd-tcp-address=127.0.0.1:4160 -max-msg-size=8388608

# Get the admin token for the TS server
until [ `docker inspect -f "{{.State.Status}}" teamspeak | grep -c running` -eq 1 ]
do
    sleep 1
    docker logs teamspeak
    echo TeamSpeak not ready yet.
done
echo TeamSpeak Ready
docker exec teamspeak /bin/bash -c 'cat /data/logs/*.log | grep -oE "token=([A-Za-z0-9+_=]+)" | cut -d= -f2' > teamspeakToken.txt

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

