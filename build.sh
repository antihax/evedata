#!/bin/bash

go get -u ./...
git --git-dir=../../bwmarrin/discordgo checkout develop

for dir in cmd/*/; do
    dir=${dir%/}
    if grep -q '^package main$' $dir/*.go 2>/dev/null; then
        echo "building $dir"
        CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o bin/$(basename $dir) ./$dir
        docker build -t antihax/evedata-$(basename $dir) -f docker/Dockerfile.$(basename $dir) .
    else
        echo "(skipped $dir)"
    fi
done
