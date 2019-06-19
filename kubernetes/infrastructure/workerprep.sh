sed -i 's/^#*Port 22/Port 44/' /etc/ssh/sshd_config
service sshd restart

sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF > /etc/apt/sources.list.d/kubernetes.list  
deb http://apt.kubernetes.io/ kubernetes-xenial main  
EOF

mkdir -p /etc/cni/net.d

apt update -y
apt install -y docker.io apt-transport-https kubelet kubeadm kubectl kubernetes-cni
apt upgrade -y
apt autoremove -y

systemctl enable docker.service
