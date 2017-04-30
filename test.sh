#!/bin/bash
set -e
cd /go/src/github.com/antihax/evedata/
go get github.com/modocache/gover
go get -u
go test -v ./...
go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"{{end}}' ./... | xargs -i sh -c {}
gover . coverprofile.txt
bash <(curl -s https://codecov.io/bash) -f coverprofile.txt