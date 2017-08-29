#!/bin/bash
set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

# Build everything in apps to test for compile errors
for dir in cmd/*/; do
    dir=${dir%/}
    if grep -q '^package main$' $dir/*.go 2>/dev/null; then
        echo "building $dir"
        CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o bin/$(basename $dir) ./$dir
    else
        echo "(skipped $dir)"
    fi
done
