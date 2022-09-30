ConfigMap&Secret的其他问题
一、  ConfigMap&Secret热更新

如果使用yaml创建，直接修改yaml文件即可，如果是使用文件创建推荐修改原文件后使用文件进行更新。命令如下:

kubectl create configmap nginx-conf --from-file=nginx.conf --dry-run -oyaml | kubectl replace -f -

如果configmap&secret是以文件形式挂载，如果资源内容进行更新，pod内挂载的文件内容也会进行同步，如果是变量方式使用是不会进行更新变量的，需要重启Pod。文件方式热更新后不代表服务的配置也进行了热更新，如果服务本身支持热加载配置，可以实现自动更新配置，如果没有也需要重启Pod。

1.1 configmap热更新
[root@k8smaster example]# kubectl get cm ceshi -oyaml
apiVersion: v1
data:
  ui.properties: |
    color.good=purple
    color.bad=yellow
    allow.textmode=true
kind: ConfigMap
metadata:
  creationTimestamp: "2022-09-18T13:24:51Z"
  name: ceshi
  namespace: default
  resourceVersion: "88354"
  selfLink: /api/v1/namespaces/default/configmaps/ceshi
  uid: 432d6e5c-24ca-41b9-a57d-4af490e3c40d
修改文件
[root@k8smaster example]# cat ui.properties 
color.good=purple
color.bad=yellow
allow.textmode=true
a=1
热更新
[root@k8smaster example]# kubectl create configmap ceshi --from-file=ui.properties --dry-run -oyaml | kubectl replace -f -
configmap/ceshi replaced
更新后
[root@k8smaster example]# kubectl get cm ceshi -oyaml
apiVersion: v1
data:
  ui.properties: |
    color.good=purple
    color.bad=yellow
    allow.textmode=true
    a=1
kind: ConfigMap
metadata:
  creationTimestamp: "2022-09-18T13:24:51Z"
  name: ceshi
  namespace: default
  resourceVersion: "325641"
  selfLink: /api/v1/namespaces/default/configmaps/ceshi
  uid: 432d6e5c-24ca-41b9-a57d-4af490e3c40d

1.2 secret热更新
原本的
[root@k8smaster secret]# kubectl get secrets test -oyaml
apiVersion: v1
data:
  password: a3lsaW4uMjAyMAo=
  re.ip: MTkyLjE2OC4xODkuMTQzOjg0NDMK
kind: Secret
metadata:
  creationTimestamp: "2022-09-22T21:06:25Z"
  name: test
  namespace: default
  resourceVersion: "321712"
  selfLink: /api/v1/namespaces/default/secrets/test
  uid: 8460cce5-07ae-445d-b632-06ab6ebe53d1
type: Opaque
热更新
[root@k8smaster secret]# kubectl create secret generic test --from-file=file/ --dry-run -oyaml | kubectl replace -f -
secret/test replaced
更新后
[root@k8smaster secret]# kubectl get secrets test -oyaml
apiVersion: v1
data:
  password: S1lsaW4+MjAyMgo=
  re.ip: MTkyLjE2OC4xODkuMTQzOjg0NDMK
kind: Secret
metadata:
  creationTimestamp: "2022-09-22T21:06:25Z"
  name: test
  namespace: default
  resourceVersion: "326095"
  selfLink: /api/v1/namespaces/default/secrets/test
  uid: 8460cce5-07ae-445d-b632-06ab6ebe53d1
type: Opaque
二、 ConfigMap&Secret使用限制

    如果Pod使用ConfigMap&Secret，需要提前创建，应用的key必须存在
    使用的ConfigMap&Secre的Pod必须与Pod在一个命名空间
    envFrom、valueFrom无法更新环境变量，需要重启Pod，如果引用的key不存在，会忽略掉无效的key
    subPath也是无法热更新的
    ConfigMap&Secret最好不要太大，官方推荐不要大于1M，与etcd有关系
三、 ConfigMap&Secret不可变
在k8s的1.18版本后新增加了一个不可变的ConfigMap&Secret设置参数immutable: true来进行设置。1.18为Alpha版本需要在特性进行手动开启参数为--feature-gates="ImmutableEphemeralVolumes: true"，1.19以后默认为true无需手动设置。

Kubernetes 特性 不可变更的 Secret 和 ConfigMap提供了一种将各个 Secret 和 ConfigMap 设置为不可变更的选项。对于大量使用 ConfigMap 的 集群（至少有数万个各不相同的 ConfigMap 给 Pod 挂载）而言，禁止更改 ConfigMap 的数据有以下好处：

    保护应用，使之免受意外（不想要的）更新所带来的负面影响。
    通过大幅降低对 kube-apiserver 的压力提升集群性能，这是因为系统会关闭 对已标记为不可变更的 ConfigMap 的监视操作。

一旦某 ConfigMap 被标记为不可变更，则 无法 逆转这一变化，也无法更改 data 或 binaryData 字段的内容。你只能删除并重建 ConfigMap。

示例如下：

apiVersion: v1
kind: ConfigMap
metadata:
  ...
data:
  ...
immutable: true
四、对比
最后我们来对比下 Secret 和 ConfigMap这两种资源对象的异同点：

    相同点
    key/value的形式
    属于某个特定的命名空间
    可以导出到环境变量
    可以通过目录/文件形式挂载
    通过 volume 挂载的配置信息均可热更新
    不同点
    Secret 可以被 ServerAccount 关联
    Secret 可以存储 docker register 的鉴权信息，用在 ImagePullSecret 参数中，用于拉取私有仓库的镜像
    Secret 支持 Base64 加密
    Secret 分为 多种类型，而 Configmap 不区分类型