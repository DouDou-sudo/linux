刚开机执行命令，报错如下
kubectl get pods
The connection to the server 192.168.X.X:6443 was refused -did you specify...?
查看端口没有被监听
netstat -tunlp | grep 6443
查看kubelet服务是running状态
systemctl status kubelet
journalctl -xeu kubelet输出如下
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.008302    6876 kubelet.go:2248] node "k8smaster" not found
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.108655    6876 kubelet.go:2248] node "k8smaster" not found
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.209148    6876 kubelet.go:2248] node "k8smaster" not found
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.310081    6876 kubelet.go:2248] node "k8smaster" not found
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.410878    6876 kubelet.go:2248] node "k8smaster" not found
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.504547    6876 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1beta1.CSIDriver: G
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.511171    6876 kubelet.go:2248] node "k8smaster" not found
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.521052    6876 reflector.go:125] k8s.io/kubernetes/pkg/kubelet/config/apiserver.go:47: Failed to list *v1.Pod: G
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.522720    6876 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1beta1.RuntimeClass
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.524296    6876 reflector.go:125] k8s.io/kubernetes/pkg/kubelet/kubelet.go:444: Failed to list *v1.Service: Get h
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.524800    6876 reflector.go:125] k8s.io/kubernetes/pkg/kubelet/kubelet.go:453: Failed to list *v1.Node: Get http
Aug 30 14:26:36 k8smaster kubelet[6876]: E0830 14:26:36.612348    6876 kubelet.go:2248] node "k8smaster" not found
提示k8smaster找不到
网上查找下面几种情况也可能导致这个报错
1、查看hosts文件和ip地址，都正常
2、查看.kube目录配置文件，正常
3、image丢失，没有丢失
4、根分区占满，没有占满
5、查看cat /etc/kubernetes/kubelet.conf 配置文件的node名和server的ip地址
查看k8s需要启动的容器
docker ps
apiserver启动正常，过一会apiserver失败，etcd启动失败
查找apiserver的容器id
[root@k8smaster ~]# docker ps -a
CONTAINER ID        IMAGE                                               COMMAND                  CREATED             STATUS                       PORTS               NA
MES551661d93303        2c4adeb21b4f                                        "etcd --advertise-cl…"   2 minutes ago       Exited (2) 2 minutes ago                         k
8s_etcd_etcd-k8smaster_kube-system_aef1d0dcc0b4258f3863a07b0f691b97_110fe4fa6faeab0        201c7a840312                                        "kube-apiserver --ad…"   4 minutes ago       Exited (255) 4 minutes ago                       k
8s_kube-apiserver_kube-apiserver-k8smaster_kube-system_589b0a320bff7edb0b69c0c65ce3f075_111
查看apiserver的logs
docker logs fe4fa6faeab0输出如下
W0830 06:25:56.434071       1 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {127.0.0.1:2379 0  <nil>}. Err :connection error: desc = "transpor
t: Error while dialing dial tcp 127.0.0.1:2379: connect: connection refused". Reconnecting...W0830 06:25:56.443032       1 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {127.0.0.1:2379 0  <nil>}. Err :connection error: desc = "transpor
t: Error while dialing dial tcp 127.0.0.1:2379: connect: connection refused". Reconnecting...W0830 06:25:57.434599       1 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {127.0.0.1:2379 0  <nil>}. Err :connection error: desc = "transpor
t: Error while dialing dial tcp 127.0.0.1:2379: connect: connection refused". Reconnecting...W0830 06:25:58.278313       1 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {127.0.0.1:2379 0  <nil>}. Err :connection error: desc = "transpor
t: Error while dialing dial tcp 127.0.0.1:2379: connect: connection refused". Reconnecting...
查看etcd的启动日志
[root@k8smaster ~]# docker ps -a
CONTAINER ID        IMAGE                                               COMMAND                  CREATED              STATUS                            PORTS           
    NAMES
94576d4717f4        201c7a840312                                        "kube-apiserver --ad…"   About a minute ago   Exited (255) About a minute ago                  
     k8s_kube-apiserver_kube-apiserver-k8smaster_kube-system_589b0a320bff7edb0b69c0c65ce3f075_112
551661d93303        2c4adeb21b4f                                        "etcd --advertise-cl…"   4 minutes ago        Exited (2) 4 minutes ago                         
     k8s_etcd_etcd-k8smaster_kube-system_aef1d0dcc0b4258f3863a07b0f691b97_110
[root@k8smaster ~]# docker logs 551661d93303
2022-08-30 06:28:21.096721 I | etcdmain: etcd Version: 3.3.10
2022-08-30 06:28:21.096793 I | etcdmain: Git SHA: 27fc7e2
2022-08-30 06:28:21.096799 I | etcdmain: Go Version: go1.10.4
2022-08-30 06:28:21.096803 I | etcdmain: Go OS/Arch: linux/amd64
2022-08-30 06:28:21.096808 I | etcdmain: setting maximum number of CPUs to 2, total number of available CPUs is 2
2022-08-30 06:28:21.096873 N | etcdmain: the server is already initialized as member before, starting as etcd member...
2022-08-30 06:28:21.096904 I | embed: peerTLS: cert = /etc/kubernetes/pki/etcd/peer.crt, key = /etc/kubernetes/pki/etcd/peer.key, ca = , trusted-ca = /etc/kubernetes/pk
i/etcd/ca.crt, client-cert-auth = true, crl-file = 2022-08-30 06:28:21.097898 I | embed: listening for peers on https://192.168.189.200:2380
2022-08-30 06:28:21.097958 I | embed: listening for client requests on 127.0.0.1:2379
2022-08-30 06:28:21.097986 I | embed: listening for client requests on 192.168.189.200:2379
panic: freepages: failed to get all reachable pages (page 3482471299534041133: out of bounds: 594)

goroutine 64 [running]:
github.com/coreos/etcd/cmd/vendor/github.com/coreos/bbolt.(*DB).freepages.func2(0xc4200805a0)
	/tmp/etcd-release-3.3.10/etcd/release/etcd/gopath/src/github.com/coreos/etcd/cmd/vendor/github.com/coreos/bbolt/db.go:976 +0xfb
created by github.com/coreos/etcd/cmd/vendor/github.com/coreos/bbolt.(*DB).freepages
	/tmp/etcd-release-3.3.10/etcd/release/etcd/gopath/src/github.com/coreos/etcd/cmd/vendor/github.com/coreos/bbolt/db.go:974 +0x1b7
此panic是因为虚拟机异常掉电，导致etcd数据损坏导致的
最终解决方法
在故障节点上停止etcd服务并删除损坏的etcd数据，这时etcd本身没有启动，将数据备份，然后reload，重启kubelet
此方法会把你所有的pod都丢失掉
[root@k8smaster lib]# cd etcd/
[root@k8smaster etcd]# cd member/
[root@k8smaster member]# ls
snap  wal
[root@k8smaster member]# mkdir ../bak
[root@k8smaster member]# mv * ../bak/
[root@k8smaster member]# ls
[root@k8smaster member]# systemctl daemon-reload
[root@k8smaster member]# systemctl restart kubelet





