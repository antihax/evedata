# VPN Setup (optional)

Run setup_ovpn.sh and follow the prompts. It may be necessary to edit and change host names.

Run setup_client.sh and follow the prompt to create the client configuration files from the local CA.

On ubuntu you will need to set `options ndots:2` in /etc/resolv.conf
```
"grep -q -F 'options ndots:2' /etc/resolv.conf || echo 'options ndots:2' >> /etc/resolv.conf"
```