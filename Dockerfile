FROM golang:latest

ADD . /go/src/github.com/antihax/evedata

RUN go install github.com/antihax/evedata

ENTRYPOINT /go/bin/evedata-server

EXPOSE 3000