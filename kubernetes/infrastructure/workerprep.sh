sed -i 's/^#*Port 22/Port 44/' /etc/ssh/sshd_config
service sshd restart

sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF > /etc/apt/sources.list.d/kubernetes.list  
deb http://apt.kubernetes.io/ kubernetes-xenial main  
EOF

mkdir -p /etc/cni/net.d

export DEBIAN_FRONTEND=noninteractive; 
apt-get update; 
apt-get install -q -y docker.io apt-transport-https kubelet kubeadm kubectl kubernetes-cni; 
apt-get autoremove -q -y; 
apt-get upgrade -q -y; 
apt-get --with-new-pkgs upgrade -q -y; 
apt-get autoremove -q -y; 
docker system prune -a -f;

systemctl enable docker.service;
