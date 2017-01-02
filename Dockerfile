FROM golang:latest

RUN go get -u github.com/antihax/evedata
RUN go install github.com/antihax/evedata
RUN mv /go/src/github.com/antihax/evedata/static /go/static
RUN mv /go/src/github.com/antihax/evedata/templates /go/templates

ENTRYPOINT /go/bin/evedata

EXPOSE 3000
