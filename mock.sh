#!/bin/bash
# remove all containers
docker ps -aq | xargs docker stop | xargs docker rm

# Mock ESI Server
docker run -d -h mock-esi -p 127.0.0.1:8080:8080 antihax/mock-esi

# NSQ Test Service
docker run --net=host -d -p 4160:4160 -p 4161:4161 -h nsqlookupd1.nsq nsqio/nsq /nsqlookupd
docker run --net=host -d -p 4171:4171 -h nsqadmin.nsq nsqio/nsq /nsqadmin --lookupd-http-address=127.0.0.1:4161 -max-msg-size=4194304
docker run --net=host -d -p 4151:4151 -p 4150:4150 -h localhost nsqio/nsq /nsqd --lookupd-tcp-address=127.0.0.1:4160
docker run --net=host --name=teamspeak -d -p 9987:9987/udp -p 30033:30033 -p 10011:10011 -p 41144:41144 mbentley/teamspeak clear_database=1 serveradmin_password=nothinguseful
sleep 10
docker exec teamspeak /bin/bash -c 'cat /data/logs/*.log | grep -oE "token=([A-Za-z0-9+_=]+)" | cut -d= -f2' > teamspeakToken.txt
