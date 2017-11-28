#!/bin/bash
# Mock ESI Server
docker run -d -h mock-esi -p 127.0.0.1:8080:8080 antihax/mock-esi

# NSQ Test Service
docker run --net=host -d -p 4160:4160 -p 4161:4161 -h nsqlookupd1.nsq nsqio/nsq /nsqlookupd
docker run --net=host -d -p 4171:4171 -h nsqadmin.nsq nsqio/nsq /nsqadmin --lookupd-http-address=127.0.0.1:4161
docker run --net=host -d -p 4151:4151 -p 4150:4150 -h localhost nsqio/nsq /nsqd --lookupd-tcp-address=127.0.0.1:4160
