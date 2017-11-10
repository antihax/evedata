#!/bin/bash
# Mock ESI Server
docker run -d -h mock-esi -p 127.0.0.1:8080:8080 antihax/mock-esi

# Mock ESI Server
docker run -d -h tokenstore -p 127.0.0.1:4001:4001 antihax/evedata-tokenstore

# kafka Test Service
docker run -d -p 2181:2181 -p 9092:9092 --env ADVERTISED_HOST="0.0.0.0" --env ADVERTISED_PORT=9092 spotify/kafka

# NSQ Test Service
docker run --net=host -d -p 127.0.0.1:4160:4160 -p 127.0.0.1:4161:4161 -h nsqlookupd1.nsq nsqio/nsq /nsqlookupd
docker run --net=host -d -p 127.0.0.1:4171:4171 -h nsqadmin.nsq nsqio/nsq /nsqadmin --lookupd-http-address=127.0.0.1:4161
docker run --net=host -d -p 127.0.0.1:4151:4151 -p 127.0.0.1:4150:4150 -h nsqd.nsq nsqio/nsq /nsqd --lookupd-tcp-address=127.0.0.1:4160

# zipkin
docker run --net=host -d -p 127.0.0.1:9411:9411 openzipkin/zipkin
