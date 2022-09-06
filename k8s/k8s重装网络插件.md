如果你想要在我们前面部署的集群中实践 Flannel 的话，可以在 Master 节点上执行如下命令来替换网络插件。
第一步，执行$ rm -rf /etc/cni/net.d/*；
第二步，执行$ kubectl delete -f "https://cloud.weave.works/k8s/net?k8s-version=1.11"；
第三步，在/etc/kubernetes/manifests/kube-controller-manager.yaml里，为容器启动命令添加如下两个参数：--allocate-node-cidrs=true--cluster-cidr=10.244.0.0/16；此处的ip地址如果自定义，需要修改下载好的kube-flannel.yml文件，将ip地址改为对应ip，然后再create
第四步， 重启所有 kubelet；第五步， 执行$ kubectl create -f https://raw.githubusercontent.com/coreos/flannel/bc79dd1505b0c8681ece4de4c0d86c5cd2643275/Documentation/kube-flannel.yml。




