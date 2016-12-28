FROM golang:latest

ADD . /go/src/github.com/shijuvar/golang-docker

RUN go install github.com/antihax/evedata

ENTRYPOINT /go/bin/golang-docker

EXPOSE 3000