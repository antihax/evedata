FROM golang:latest

RUN go get -u github.com/antihax/evedata
RUN go install github.com/antihax/evedata

ENTRYPOINT /go/bin/evedata

EXPOSE 3000
