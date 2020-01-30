# lunara-k8s
Lunara - Kubernetes

#### official repository
```sh
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
```

#### aliyun
```sh
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1 
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
```

#### Set SELinux in permissive mode (effectively disabling it)
```sh
setenforce 0
sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes

systemctl enable --now kubelet
```
#### start minikube with aliyun registry
```sh
minikube start --image-mirror-country=cn --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers
```

## Create secret
```sh
kubectl create secret docker-registry my-secret --docker-server=123.456.789.0:9595 --docker-username=admin --docker-password=XXXX --docker-email=test@xyz.com
```

#### useful links
https://github.com/kubernetes/kubernetes/issues/56850

## TODO
1. create private hub by harbor
2. k8s/client-go #
3. pull images form harbor in kunbernetes
4. generate business's configuration files
5. specific shared volume
6. configMap

#### points 
1. create redis instances include master-slave
2. create mysql instances

#### service
1. RBAC
2. create & list & update & delete all deployments & pods & services