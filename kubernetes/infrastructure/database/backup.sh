#!/bin/bash
kubectl create namespace backup
kubectl create secret generic backup-secrets -n backup \
   --from-literal=accountID="" --from-literal=applicationKey="" --from-literal=dbuser="" --from-literal=dbpass=""
kubectl apply -f backup.yaml