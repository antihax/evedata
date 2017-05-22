#!/bin/bash
set -e

for dir in apps/*/; do
    dir=${dir%/}
    if grep -q '^package main$' $dir/*.go 2>/dev/null; then
        docker push antihax
    else
        echo "(skipped $dir)"
    fi
done
