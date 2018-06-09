#!/bin/bash
cd ovpn
echo "enter unique client name:"
read CLIENTNAME
docker run --net=none --rm -it -v $PWD:/etc/openvpn kylemanna/openvpn easyrsa build-client-full $CLIENTNAME
docker run --net=none --rm -v $PWD:/etc/openvpn kylemanna/openvpn ovpn_getclient $CLIENTNAME > ../$CLIENTNAME.ovpn
sed -i 's/ 1194/ 31194/g' ../$CLIENTNAME.ovpn

IP=(kubectl get service kube-dns -n kube-system -o custom-columns=:.spec.clusterIP | sed -n 2p)

cat <<EOT >> ../$CLIENTNAME.ovpn
script-security 2
dhcp-option DNS $IP
dhcp-option DOMAIN evedata

up /etc/openvpn/update-resolv-conf
down /etc/openvpn/update-resolv-conf
EOT

echo "created $CLIENTNAME"