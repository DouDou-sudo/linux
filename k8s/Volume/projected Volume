在 Kubernetes 中，有几种特殊的 Volume，它们存在的意义不是为了存放容器里的数据，也不是用来进行容器和宿主机之间的数据交换。这些特殊 Volume 的作用，是为容器提供预先定义好的数据。所以，从容器的角度来看，这些 Volume 里的信息就是仿佛是被 Kubernetes“投射”（Project）进入容器当中的。这正是 Projected Volume 的含义。
到目前为止，Kubernetes 支持的 Projected Volume 一共有四种：
Secret；
ConfigMap；
Downward API；
ServiceAccountToken。
1、secret
它的作用，是帮你把 Pod 想要访问的加密数据，存放到 Etcd 中。然后，你就可以通过在 Pod 的容器里挂载 Volume 的方式，访问到这些 Secret 里保存的信息了。
比如；莫过于存放数据库的 Credential 信息，比如下面这个例子
[root@k8smaster build-yaml]# cat test-projected-volume.yml 

apiVersion: v1
kind: Pod
metadata:
  name: test-projected-volume 
spec:
  containers:
  - name: test-secret-volume
    image: busybox
    args:
    - sleep
    - "86400"
    volumeMounts:
    - name: mysql-cred
      mountPath: "/projected-volume"
      readOnly: true
  volumes:
  - name: mysql-cred
    projected:
      sources:
      - secret:
          name: user
      - secret:
          name: password
在这个pod中，定义了一个简单的容器，它声明挂载的 Volume，并不是常见的 emptyDir 或者 hostPath 类型，而是 projected 类型。而这个 Volume 的数据来源（sources），则是名为 user 和 pass 的 Secret 对象，分别对应的是数据库的用户名和密码。
[root@k8smaster build-yaml]# cat username.txt 
admin
[root@k8smaster build-yaml]# cat password.txt 
123wqeq
建立user和pass的键值对关系
[root@k8smaster build-yaml]# kubectl create secret generic user --from-file=./username.txt
[root@k8smaster build-yaml]#  kubectl create secret generic password --from-file=./password.txt
查看secret对象
[root@k8smaster build-yaml]# kubectl get secret
NAME                  TYPE                                  DATA   AGE
default-token-8d4kq   kubernetes.io/service-account-token   3      17h
password              Opaque                                1      73m
user                  Opaque                                1      73m
创建这个pod
[root@k8smaster build-yaml]# kubectl apply -f test-projected-volume.yml 
pod/test-projected-volume created
[root@k8smaster build-yaml]#  kubectl exec -it test-projected-volume -- /bin/sh
/ # cd projected-volume/
/projected-volume # ls
password.txt  username.txt
/projected-volume # cat password.txt 
123wqeq
/projected-volume # cat username.txt 
admin
2、configmap
ConfigMap 保存的是不需要加密的、应用所需的配置信息。而 ConfigMap 的用法几乎与 Secret 完全相同：你可以使用 kubectl create configmap 从文件或者目录创建 ConfigMap，也可以直接编写 ConfigMap 对象的 YAML 文件。

[root@k8smaster build-yaml]# cat example/ui.properties
color.good=purple
color.bad=yellow
allow.textmode=true
创建configmap
[root@k8smaster build-yaml]# kubectl create configmap ui-config --from-file=example/ui.properties 
configmap/ui-config created
[root@k8smaster build-yaml]# kubectl get configmaps 
NAME        DATA   AGE
ui-config   1      11s
查看这个configmap里保存的data
[root@k8smaster build-yaml]# kubectl get configmaps ui-config -o yaml
apiVersion: v1
data:
  ui.properties: |
    color.good=purple
    color.bad=yellow
    allow.textmode=true
kind: ConfigMap
metadata:
  creationTimestamp: "2022-08-25T04:21:01Z"
  name: ui-config
  namespace: default
  resourceVersion: "34527"
  selfLink: /api/v1/namespaces/default/configmaps/ui-config
  uid: 18cb8da2-3b6f-4c3a-bc14-0f7fde3c84ae
修改test-projected-volume.yml文件的volume
[root@k8smaster build-yaml]# cat test-projected-volume.yml 
apiVersion: v1
kind: Pod
metadata:
  name: test-projected-volume 
spec:
  containers:
  - name: test-secret-volume
    image: busybox
    args:
    - sleep
    - "86400"
    volumeMounts:
    - name: java
      mountPath: "/projected-volume"
      readOnly: true
  volumes:
  - name: java
    projected:
      sources:
      - configMap:      //注意此处的大小写
          name: ui-config

[root@k8smaster build-yaml]# kubectl exec -it test-projected-volume -- /bin/sh
/projected-volume # ls
ui.properties
/projected-volume # cat ui.properties 
color.good=purple
color.bad=yellow
allow.textmode=true

3、dowmward api
它的作用是：让 Pod 里的容器能够直接获取到这个 Pod API 对象本身的信息。
[root@k8smaster build-yaml]# cat test-downwardapi-volume.yml 
apiVersion: v1
kind: Pod
metadata:
  name: test-downwardapi-volume
  labels:
    zone: us-est-coast
    cluster: test-cluster1
    rack: rack-22
spec:
  containers:
    - name: client-container
      image: busybox
      command: ["sh", "-c"]
      args:
      - while true; do
          if [[ -e /etc/podinfo/labels ]]; then
            echo -en '\n\n'; cat /etc/podinfo/labels; fi;
          sleep 5;
        done;
      volumeMounts:
        - name: podinfo
          mountPath: /etc/podinfo
          readOnly: false
  volumes:
    - name: podinfo
      projected:
        sources:
        - downwardAPI:
            items:
              - path: "labels"
                fieldRef:
                  fieldPath: metadata.labels

[root@k8smaster build-yaml]# kubectl logs test-downwardapi-volume 
通过这样的声明方式，当前 Pod 的 Labels 字段的值，就会被 Kubernetes 自动挂载成为容器里的 /etc/podinfo/labels 文件。
[root@k8smaster build-yaml]# kubectl exec -it test-downwardapi-volume -- /bin/sh
/etc/podinfo # cd ..
/etc # cat /etc/podinfo/labels 
cluster="test-cluster1"
rack="rack-22"
zone="us-est-coast"/etc # 
而这个容器的启动命令，则是不断打印出 /etc/podinfo/labels 里的内容。所以，当我创建了这个 Pod 之后，就可以通过 kubectl logs 指令，查看到这些 Labels 字段被打印出来，如下所示：
cluster="test-cluster1"
rack="rack-22"
zone="us-est-coast"

cluster="test-cluster1"
rack="rack-22"
zone="us-est-coast"
目前，Downward API 支持的字段已经非常丰富了，比如：

1. 使用fieldRef可以声明使用:
spec.nodeName - 宿主机名字
status.hostIP - 宿主机IP
metadata.name - Pod的名字
metadata.namespace - Pod的Namespace
status.podIP - Pod的IP
spec.serviceAccountName - Pod的Service Account的名字
metadata.uid - Pod的UID
metadata.labels['<KEY>'] - 指定<KEY>的Label值
metadata.annotations['<KEY>'] - 指定<KEY>的Annotation值
metadata.labels - Pod的所有Label
metadata.annotations - Pod的所有Annotation

2. 使用resourceFieldRef可以声明使用:
容器的CPU limit
容器的CPU request
容器的memory limit
容器的memory request
不过，需要注意的是，Downward API 能够获取到的信息，一定是 Pod 里的容器进程启动之前就能够确定下来的信息。而如果你想要获取 Pod 容器运行后才会出现的信息，比如，容器进程的 PID，那就肯定不能使用 Downward API 了，而应该考虑在 Pod 里定义一个 sidecar 容器。

其实，Secret、ConfigMap，以及 Downward API 这三种 Projected Volume 定义的信息，大多还可以通过环境变量的方式出现在容器里。但是，通过环境变量获取这些信息的方式，不具备自动更新的能力。所以，一般情况下，我都建议你使用 Volume 文件的方式获取这些信息。
Service Account 对象的作用，就是 Kubernetes 系统内置的一种“服务账户”，它是 Kubernetes 进行权限分配的对象。
比如，Service Account A，可以只被允许对 Kubernetes API 进行 GET 操作，而 Service Account B，则可以有 Kubernetes API 的所有操作权限。像这样的 Service Account 的授权信息和文件，实际上保存在它所绑定的一个特殊的 Secret 对象里的。
这个特殊的 Secret 对象，就叫作 ServiceAccountToken。任何运行在 Kubernetes 集群上的应用，都必须使用这个 ServiceAccountToken 里保存的授权信息，也就是 Token，才可以合法地访问 API Server。所以说，Kubernetes 项目的 Projected Volume 其实只有三种，因为第四种 ServiceAccountToken，只是一种特殊的 Secret 而已。
另外，为了方便使用，Kubernetes 已经为你提供了一个默认“服务账户”（default Service Account）。并且，任何一个运行在 Kubernetes 里的 Pod，都可以直接使用这个默认的 Service Account，而无需显示地声明挂载它。

Service Account默认开启，在 /var/run/secrets/kubernetes.io/serviceaccount/目录下，可以禁用
[root@k8smaster build-yaml]# kubectl exec -it test-projected-volume  -- /bin/sh
/ # ls /var/run/secrets/kubernetes.io/serviceaccount/
ca.crt     namespace  token

