FROM golang:latest

RUN go get -u github.com/antihax/evedata
RUN go install github.com/antihax/evedata
COPY /go/src/github.com/antihax/evedata/static /go/static
COPY /go/src/github.com/antihax/evedata/templates /go/templates

ENTRYPOINT /go/bin/evedata

EXPOSE 3000
