#!/bin/bash
if [[ -z "$1" ]] 
then
    for dir in cmd/*/; do
        dir=${dir%/}
        if grep -q '^package main$' $dir/*.go 2>/dev/null; then
            docker push antihax/evedata-$(basename $dir)
        else
            echo "(skipped $dir)"
        fi
    done
else
    dir=$1
    docker push antihax/evedata-$dir
fi