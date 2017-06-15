#!/bin/sh
kubectl create configmap ca-certificates --from-file=/etc/ssl/certs/ca-certificates.crt --namespace=evedata