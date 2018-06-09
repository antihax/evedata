#!/bin/bash

mkdir ovpn && cd ovpn

docker run -e EASYRSA_KEY_SIZE=4096 --net=none --rm -t -i -v $PWD:/etc/openvpn kylemanna/openvpn ovpn_genconfig -u udp://vpn.evedata.org
docker run -e EASYRSA_KEY_SIZE=4096 --net=none --rm -t -i -v $PWD:/etc/openvpn kylemanna/openvpn ovpn_initpki
docker run -e EASYRSA_KEY_SIZE=4096 --net=none --rm -t -i -v $PWD:/etc/openvpn kylemanna/openvpn ovpn_copy_server_files
kubectl create namespace ovpn
kubectl create secret generic ovpn-key -n ovpn --from-file=server/pki/private/vpn.evedata.org.key
kubectl create secret generic ovpn-cert -n ovpn --from-file=server/pki/issued/vpn.evedata.org.crt
kubectl create secret generic ovpn-pki -n ovpn \
    --from-file=server/pki/ca.crt --from-file=server/pki/dh.pem --from-file=server/pki/ta.key
kubectl create configmap ovpn-conf -n ovpn --from-file=server/
kubectl create configmap ccd -n ovpn --from-file=server/ccd
kubectl apply -f ../ovpn.yaml