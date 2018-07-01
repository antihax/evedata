# move SSH port to stop search bots hitting it and poluting logs
sed -i 's/^#*Port 22/Port 44/' /etc/ssh/sshd_config
service sshd restart

# turn swap off
sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

# add kubernetes repo to apt
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF > /etc/apt/sources.list.d/kubernetes.list  
deb http://apt.kubernetes.io/ kubernetes-xenial main  
EOF

# bugfix: won't create CNI on occasion if this dir is not present
mkdir -p /etc/cni/net.d

# update apt repos and install docker, pwgen, and kubernetes
apt update -y
apt install -y docker.io apt-transport-https pwgen
apt update -y
apt install -y kubelet kubeadm kubectl kubernetes-cni
apt upgrade -y
apt autoremove -y

# auto start docker
systemctl enable docker.service

# start master node
kubeadm init

# copy kubectl config
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
export kubeconfig=/etc/kubernetes/admin.conf

# create weavenet and set random encryption password
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
kubectl create secret -n kube-system generic weave-passwd --from-literal=weave-passwd=`pwgen 32 1 -sy`
kubectl patch daemonset -n kube-system --type=json weave-net -p '[{"op": "add", "path": "/spec/template/spec/containers/0/env/1", "value": {"name":"WEAVE_PASSWORD","valueFrom":{"secretKeyRef":{"name":"weave-passwd","key":"weave-passwd"}}}}]'
kubectl patch daemonset -n kube-system --type=json weave-net -p '[{"op": "add", "path": "/spec/template/spec/containers/0/env/1", "value": {"name":"EXTRA_ARGS","value":"--log-level=error"}}]'

# fix local host DNS issues
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
data:
  upstreamNameservers: |-
    ["8.8.8.8", "1.1.1.1"]
EOF

# restart kubernetes pods to make sure everything is now clean
kubectl delete pod --all -n kube-system