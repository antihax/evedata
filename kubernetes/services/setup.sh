#!/bin/bash
kubectl create namespace evedata

kubectl apply -f tokenserver.yaml
kubectl apply -f conservator.yaml
kubectl apply -f hammer.yaml
kubectl apply -f nail.yaml
kubectl apply -f mailserver.yaml
kubectl apply -f zkilboard.yaml
kubectl apply -f vanguard.yaml
kubectl apply -f artifice.yaml