
#### 将cgroup驱动迁移到systemd

##### 1.修改 kubelet 的 ConfigMap

    运行 kubectl edit cm kubelet-config -n kube-system。

    修改现有 cgroupDriver 的值，或者新增如下式样的字段：

    cgroupDriver: systemd

    该字段必须出现在 ConfigMap 的 kubelet: 小节下。

##### 2.更新所有节点的 cgroup 驱动

对于集群中的每一个节点：

    执行命令 kubectl drain <node-name> --ignore-daemonsets，以 腾空节点
    执行命令 systemctl stop kubelet，以停止 kubelet
    停止容器运行时
    修改容器运行时 cgroup 驱动为 systemd
    在文件 /var/lib/kubelet/config.yaml 中添加设置 cgroupDriver: systemd
    启动容器运行时
    执行命令 systemctl start kubelet，以启动 kubelet
    执行命令 kubectl uncordon <node-name>，以 取消节点隔离

在节点上依次执行上述步骤，确保工作负载有充足的时间被调度到其他节点。