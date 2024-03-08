k8s升级
在 master 节点执行
kubeadm config view
创建kubeadm-config-upgrade.yaml配置文件，文件内容如下，根据前面 kubeadm config view 的执行结果，修改了如下字段：
imageRepository 的值修改为：registry.aliyuncs.com/google_containers
#kubernetesVersion 的值修改为： v1.20.7

apiServer:
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns:
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.aliyuncs.com/google_containers
kind: ClusterConfiguration
kubernetesVersion: v1.20.7
networking:
  dnsDomain: cluster.local
  podSubnet: 10.44.0.0/16
  serviceSubnet: 10.22.0.0/16
scheduler: {}

# 查看配置文件差异
kubeadm upgrade diff --config kubeadm-config-upgrade.yaml
     
# 执行升级前试运行
kubeadm upgrade apply --config kubeadm-config-upgrade.yaml --dry-run
     
# 执行升级动作
kubeadm upgrade apply --config kubeadm-config-upgrade.yaml

