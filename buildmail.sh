#!/bin/bash


    dir=cmd/mailserver/
        echo "building $dir"
        CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o bin/$(basename $dir) ./$dir
        docker build -t antihax/evedata-$(basename $dir) -f docker/Dockerfile.$(basename $dir) .
        docker push antihax/evedata-$(basename $dir)
