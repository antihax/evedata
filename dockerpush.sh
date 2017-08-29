#!/bin/bash
set -e

for dir in cmd/*/; do
    dir=${dir%/}
    if grep -q '^package main$' $dir/*.go 2>/dev/null; then
        docker push antihax/evedata-$(basename $dir)
    else
        echo "(skipped $dir)"
    fi
done
