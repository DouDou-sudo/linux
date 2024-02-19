##### 一、kubeadm init时报错(swap is enable)
```
[WARNING Swap]: swap is enabled; production deployments should disable swap unless testing the NodeSwap feature gate of the kubelet
```
解决方案，执行如下命令，临时关闭swap
swapoff -a
永久关闭swap
sed -i '/swap/s/^/#/g' /etc/fstab

##### 二、kubeadm init时报错(container runtime is not running)
```
error execution phase preflight: [preflight] Some fatal errors occurred:
[ERROR CRI]: container runtime is not running: output: time="2023-07-21T09:20:07Z" level=fatal msg="validate service connection: CRI v1 runtime API is not implemented for endpoint \"unix:///var/run/containerd/containerd.sock\": rpc error: code = Unimplemented desc = unknown service runtime.v1.RuntimeService"
, error: exit status 1
```
查看当前containerd的配置containerd config dump,确认是以下哪个问题

1. 导出默认配置，config.toml这个文件默认是不存在的
    containerd config default > /etc/containerd/config.toml
    查看配置文件cat /etc/containerd/config.toml,如果禁用了cri插件，删除cri
    disabled_plugins = ["cri"]
    将cri删除重启containerd服务
    systemctl daemon-reload
    systemctl restart containerd
2. 查看配置文件cat /etc/containerd/config.toml
```
    [grpc]
    address = "/run/containerd/containerd.sock"
    修改如下
    [grpc]
    address = "/var/run/containerd/containerd.sock"
    # address = "/run/containerd/containerd.sock"
    重启containerd服务
    systemctl daemon-reload
    systemctl restart containerd
```

##### 三、kubeadm init时报错(node \"master\" not found)
```
[kubelet.go:2419] "Error getting node" err="node \"master\" not found"
```
1. --apiserver-advertise-address=检查ip地址是否写正确

2. 查看主机名和/etc/hosts文件，主机名解析是否正确
kubeadm init \
  --apiserver-advertise-address=192.168.189.200  \
  --image-repository registry.aliyuncs.com/google_containers \
  --control-plane-endpoint=cluster-endpoint \
  --kubernetes-version v1.24.1 \
  --service-cidr=10.1.0.0/16 \
  --pod-network-cidr=10.244.0.0/16 \
  --v=5
此次使用的--control-plane-endpoint为设置的dns，设置为dns是为了后面部署为高可用集群，我没有在hosts文件中设置dns解析导致一直报这个错，添加如下解析，重新init就正常了
[root@master ~]# cat /etc/hosts
192.168.189.200 cluster-endpoint


3. 需要配置cri-docker.service文件，ExecStart=/usr/bin/cri-dockerd项后面指定你的指定你的pause版本，没遇到过
例如：–pod-infra-container-image=registry.aliyuncs.com/google_containers/pause:3.7
cat /usr/lib/systemd/system/cri-docker.service 
ExecStart=/usr/bin/cri-dockerd --pod-infra-container-image=registry.aliyuncs.com/google_containers/pause:3.7 --container-runtime-endpoint
重启服务
systemctl daemon-reload && systemctl restart cri-docker.service

##### 四、kubeadm join时报错(/usr/bin/kubeadm: No such file or directory)
```
[root@k8snode2 ~]# kubeadm join cluster-endpoint:6443 --token b8ufi2.1qs3rw48dehh5yqx \
> --discovery-token-ca-cert-hash sha256:05fff7ddf0b5a6b9d4a7fd40de6d926019190773e8a067462625314d3548c987 \
> --control-plane 
-bash: /usr/bin/kubeadm: No such file or directory
```
yum安装kubeadm-1.24.1，会自动安装kubelet和kubectl，自动安装的kubectl和kubelet为新版本，需要指定具体版本安装

##### 五、kubeadm join时报错(hostname "k8snode2" could not be reached)
```
[root@k8snode2 ~]# kubeadm join cluster-endpoint:6443 --token b8ufi2.1qs3rw48dehh5yqx \
> --discovery-token-ca-cert-hash sha256:05fff7ddf0b5a6b9d4a7fd40de6d926019190773e8a067462625314d3548c987 \
> --control-plane 
[preflight] Running pre-flight checks
	[WARNING Hostname]: hostname "k8snode2" could not be reached
	[WARNING Hostname]: hostname "k8snode2": lookup k8snode2 on 192.168.189.2:53: no such host
error execution phase preflight: [preflight] Some fatal errors occurred:
	[ERROR FileAvailable--etc-kubernetes-kubelet.conf]: /etc/kubernetes/kubelet.conf already exists
	[ERROR FileAvailable--etc-kubernetes-bootstrap-kubelet.conf]: /etc/kubernetes/bootstrap-kubelet.conf already exists
	[ERROR Port-10250]: Port 10250 is in use
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
To see the stack trace of this error execute with --v=5 or higher
```
主机名不对，hosts写的是node2
[root@k8snode2 ~]# hostnamectl set-hostname node2
[root@node2 ~]# su
##### 六、kubeadm join时报错(kubelet.conf already exists)
```
[root@node2 ~]# kubeadm join cluster-endpoint:6443 --token b8ufi2.1qs3rw48dehh5yqx \
> --discovery-token-ca-cert-hash sha256:05fff7ddf0b5a6b9d4a7fd40de6d926019190773e8a067462625314d3548c987 \
> --control-plane 
[preflight] Running pre-flight checks
error execution phase preflight: [preflight] Some fatal errors occurred:
	[ERROR FileAvailable--etc-kubernetes-kubelet.conf]: /etc/kubernetes/kubelet.conf already exists
	[ERROR FileAvailable--etc-kubernetes-bootstrap-kubelet.conf]: /etc/kubernetes/bootstrap-kubelet.conf already exists
	[ERROR Port-10250]: Port 10250 is in use
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
To see the stack trace of this error execute with --v=5 or higher
```
配置文件已存在，删除相关配置
[root@node2 ~]# rm -rf /etc/kubernetes/
##### 七、kubeadm join时报错([ERROR Port-10250]: Port 10250 is in use)
```
[root@node2 ~]# kubeadm join cluster-endpoint:6443 --token b8ufi2.1qs3rw48dehh5yqx --discovery-token-ca-cert-hash sha256:05fff7ddf0b5a6b9d4a7fd40de6d926019190773e8a067462625314d3548c987 --control-plane 
[preflight] Running pre-flight checks
error execution phase preflight: [preflight] Some fatal errors occurred:
	[ERROR Port-10250]: Port 10250 is in use
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
To see the stack trace of this error execute with --v=5 or higher
```
10250端口已被占用，查看kubelet是running状态
后来发现node2节点被其他集群使用了，里面有配置。。。。
[root@node2 ~]# kubeadm reset
[root@master ~]# kubectl delete node node2
node "node2" deleted

##### 八、kubeadm join时报错(kubelet service is not enabled，[ERROR CRI]: container runtime is not running)
```
[root@node2 ~]#   kubeadm join cluster-endpoint:6443 --token b8ufi2.1qs3rw48dehh5yqx \
> --discovery-token-ca-cert-hash sha256:05fff7ddf0b5a6b9d4a7fd40de6d926019190773e8a067462625314d3548c987 \
> --control-plane 
[preflight] Running pre-flight checks
	[WARNING Service-Kubelet]: kubelet service is not enabled, please run 'systemctl enable kubelet.service'
error execution phase preflight: [preflight] Some fatal errors occurred:
	[ERROR CRI]: container runtime is not running: output: time="2023-12-28T10:25:29+08:00" level=fatal msg="validate service connection: CRI v1 runtime API is not implemented for endpoint \"unix:///var/run/containerd/containerd.sock\": rpc error: code = Unimplemented desc = unknown service runtime.v1.RuntimeService", error: exit status 1
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
To see the stack trace of this error execute with --v=5 or higher
```
[root@node2 ~]# systemctl enable kubelet.service
修改containerd配置文件，重启containerd
修改socket文件路径
```
cat /etc/containerd/config.toml 
[grpc]
  address = "/var/run/containerd/containerd.sock"
#  address = "/run/containerd/containerd.sock"
```

##### 九、kubeadm join master时报错(failure loading certificate for CA)手动证书发放
```
[root@node2 ~]#   kubeadm join cluster-endpoint:6443 --token b8ufi2.1qs3rw48dehh5yqx \
> --discovery-token-ca-cert-hash sha256:05fff7ddf0b5a6b9d4a7fd40de6d926019190773e8a067462625314d3548c987 \
> --control-plane 
[preflight] Running pre-flight checks
	[WARNING Service-Kubelet]: kubelet service is not enabled, please run 'systemctl enable kubelet.service'
[preflight] Reading configuration from the cluster...
[preflight] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -o yaml'
error execution phase preflight: 
One or more conditions for hosting a new control plane instance is not satisfied.

[failure loading certificate for CA: couldn't load the certificate file /etc/kubernetes/pki/ca.crt: open /etc/kubernetes/pki/ca.crt: no such file or directory, failure loading key for service account: couldn't load the private key file /etc/kubernetes/pki/sa.key: open /etc/kubernetes/pki/sa.key: no such file or directory, failure loading certificate for front-proxy CA: couldn't load the certificate file /etc/kubernetes/pki/front-proxy-ca.crt: open /etc/kubernetes/pki/front-proxy-ca.crt: no such file or directory, failure loading certificate for etcd CA: couldn't load the certificate file /etc/kubernetes/pki/etcd/ca.crt: open /etc/kubernetes/pki/etcd/ca.crt: no such file or directory]
Please ensure that:
* The cluster has a stable controlPlaneEndpoint address.
* The certificates that must be shared among control plane instances are provided.

To see the stack trace of this error execute with --v=5 or higher
```
相关证书文件没有拷贝到要加入的master节点上，创建目录从master节点拷贝证书到新的master节点上，再重新join
[root@node2 ~]# mkdir -p  /etc/kubernetes/pki/etcd
[root@master pki]# scp -rp /etc/kubernetes/pki/ca.* node2:/etc/kubernetes/pki
[root@master pki]#  scp -rp /etc/kubernetes/pki/sa.* node2:/etc/kubernetes/pki
[root@master pki]#  scp -rp /etc/kubernetes/pki/front-proxy-ca.* node2:/etc/kubernetes/pki
[root@master pki]#  scp -rp /etc/kubernetes/pki/etcd/ca.* node2:/etc/kubernetes/pki/etcd

##### 十、kubectl get node报错(server 192.168.X.X:6443 was refused -did you specify)
查看端口没有被监听
netstat -tunlp | grep 6443
查看kubelet服务是running状态
systemctl status kubelet
journalctl -xeu kubelet输出如下
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.008302    6876 kubelet.go:2248] node "k8smaster" not found
提示找不到k8smaster节点
网上查找下面几种情况也可能导致这个报错
1. 查看hosts文件和ip地址，都正常
2. 查看.kube目录配置文件，正常
3. image丢失，没有丢失
4. 根分区占满，没有占满
5. 查看cat /etc/kubernetes/kubelet.conf 配置文件的node名和server的ip地址
查看各组件状态，apiserver启动正常，过一会apiserver失败，etcd启动失败
查找apiserver的容器id
```
[root@k8smaster ~]# docker ps -a
CONTAINER ID        IMAGE                                               COMMAND                  CREATED             STATUS                       PORTS               NA
MES551661d93303        2c4adeb21b4f                                        "etcd --advertise-cl…"   2 minutes ago       Exited (2) 2 minutes ago                         k
8s_etcd_etcd-k8smaster_kube-system_aef1d0dcc0b4258f3863a07b0f691b97_110fe4fa6faeab0        201c7a840312                                        "kube-apiserver --ad…"   4 minutes ago       Exited (255) 4 minutes ago                       k
8s_kube-apiserver_kube-apiserver-k8smaster_kube-system_589b0a320bff7edb0b69c0c65ce3f075_111
```
查看apiserver的logs
```
docker logs fe4fa6faeab0
W0830 06:25:56.434071       1 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {127.0.0.1:2379 0  <nil>}. Err :connection error: desc = "transpor
t: Error while dialing dial tcp 127.0.0.1:2379: connect: connection refused". Reconnecting...W0830 06:25:56.443032       1 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {127.0.0.1:2379 0  <nil>}. Err :connection error: desc = "transpor
```
查看etcd的logs
```
[root@k8smaster ~]# docker logs 551661d93303
panic: freepages: failed to get all reachable pages (page 3482471299534041133: out of bounds: 594)

goroutine 64 [running]:
github.com/coreos/etcd/cmd/vendor/github.com/coreos/bbolt.(*DB).freepages.func2(0xc4200805a0)
	/tmp/etcd-release-3.3.10/etcd/release/etcd/gopath/src/github.com/coreos/etcd/cmd/vendor/github.com/coreos/bbolt/db.go:976 +0xfb
created by github.com/coreos/etcd/cmd/vendor/github.com/coreos/bbolt.(*DB).freepages
	/tmp/etcd-release-3.3.10/etcd/release/etcd/gopath/src/github.com/coreos/etcd/cmd/vendor/github.com/coreos/bbolt/db.go:974 +0x1b7
```
此panic是因为虚拟机异常掉电，导致etcd数据损坏导致的
在故障节点上停止etcd服务并删除损坏的etcd数据，这时etcd本身没有启动，将数据备份，然后reload，重启kubelet
>注意:这样会导致集群配置文件都丢失
[root@k8smaster ~]# mv /var/lib/etcd/member /var/lib/etcd/member.bak
[root@k8smaster ~]# systemctl daemon-reload
[root@k8smaster ~]# systemctl restart kubelet
