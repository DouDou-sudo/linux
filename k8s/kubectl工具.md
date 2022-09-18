一、kubectl备忘录
1.1 kubectl的bash自动补全

在 bash中设置当前 shell的自动补全，需先安装 bash-completion包否则会报错

#方法1：临时生效退出重新登录会失效
source <(kubectl completion bash)
#方法2：永久生效写入到/etc/profile文件中,这样会导致所有用户都会生效，如果只给一个用户生效请写到~/.bashrc文件中
echo 'source <(kubectl completion bash)' >>/etc/profile

您还可以为 kubectl 使用一个速记别名，该别名也可以与 completion 一起使用

cat ~/.bashrc
source <(kubectl completion bash)
alias k=kubectl    #别名
complete -F __start_kubectl k #设置completion

#执行命令如下
k get pod 
1.2 kubectl上下文和配置

kubectl客户端工具与k8s集群通信，需要一个上下文配置文件，之后需要把这个配置文件拷贝到执行kubectl命令的用户家目录 ~/.kube/config或者，把配置文件写入到环境变量 KUBECONFIG中。
示例如下：

#拷贝文件
[root@km1-81 ~]# ll ~/.kube/config 
-rw------- 1 root root 6255 Dec 14 10:54 /root/.kube/config
#设置变量，变量可以设置多个配置文件使用:隔开
export KUBECONFIG=~/.kube/config:~/.kube/config2

设置 kubectl 与哪个 Kubernetes 集群进行通信并修改配置信息。 上下文的切换一般在多kubernetes集群中使用同一个 kubectl时，需要使用上下文切换所操作的集群。单个集群无需操作。
#显示kubeconfig上下文配置
kubectl config view
#显示上下文列表
kubectl config get-contexts
#展示当前所处的上下文
kubectl config current-context 
#设置默认上下文
kubectl config use-context my-cluster-name
1.3 kubectl查看资源配置文件的说明文档

kubernetes中有各种资源，创建这些资源需要创建yaml文件，不同的kubernetes版本，资源配置可能不太一样，所以需要查看官方的一些文档。

示例如下:

#查看Pod文档
kubectl explain pod.spec

注意：yaml配置文件以空格来区分上下级，所以查看某个配置时，需要用 .来区别上下级。
1.4 kubectl创建资源

kubectl创建资源的命令有俩个分别是 apply与 create他们的区别分别是

    apply:一般使用-f参数指定yaml文件创建资源，如果这个资源已经存在，他会检查已有资源跟yaml文件中定义是否一致，不一致会进行更新，一致不对已有资源进行更改，没有则创建。
    create:可以直接使用命令创建资源(不推荐),也可以使用-f指定yaml文件，如果资源存在则不创建，不存在则创建。
create的一些实用操作：

生成资源模板文件，并不创建资源
#创建deploy
kubectl create deployment nginx-2 --image=10.122.6.81:5000/image/nginx:v1 --dry-run -oyaml
#创建svc
kubectl create service clusterip my-cs --tcp=5678:8080 --dry-run -oyaml

--dir-run   #不做创建操作
-oyaml      #输出为yaml文件

#如果不会写创建资源命令可以使用-h查看帮助，如
kubectl create deploy -h
kubectl create svc -h
....

1.5 kubectl删除资源

kubectl可以使用 delete删除资源，可以使用 -f参数指定yaml文件删除资源

kubectl delete deploy nginx

kubectl delete -f nginx.yaml

1.6 查看与查找资源

kubectl查看资源的命令有俩个分别是 get、describe、diff区别为：

    get:一般为查看资源的列表，也可以使用-o参数把查找的资源输出为别的文件类型。
    describe:查看资源的详细信息，可以检查资源的运行事件帮助排错。
    diff：可以用来比较资源差异性，可以使用-f指定yaml文件。

示例：
#查找资源
kubectl get deployments.apps

#查找一个资源并且输出为yaml文件
kubectl get deployments.apps  nginx -oyaml -A
#查看pv并且以容量排序
kubectl get pv --sort-by=spec.capacity.storage
#列出事件（Events），按时间戳排序
kubectl get events --sort-by=.metadata.creationTimestamp -A

-owide  #以列表方式查看更详细的信息
-oyaml  #以yaml文件格式输出
-A      #所有命名空间
--sort-by  #排序，以资源中定义的某些信息排序
--show-labels #显示资源的labels
-l          #显示匹配标签的资源

#查看一个资源的详细信息
kubectl describe nodes 10.122.6.75

1.7 与Pod进行交互

查看容器日志 log

#追踪容器日志
kubectl logs nginx -c nginx  -f

#容器执行命令
kubectl exec nginx-1-84c69559fd-cxmhr -c centos -- ls

#进入Pod终端
kubectl exec -it nginx-1-84c69559fd-cxmhr -c centos bash

#在本地打开端口并转发到pod上，需要有socat命令
kubectl port-forward --address 0.0.0.0 nginx-1-84c69559fd-cxmhr 80:80

#显示给定 Pod 和其中容器的监控数据
kubectl top pod nginx --containers

1.8 格式化输出

要以特定格式将详细信息输出到终端窗口，将 -o（或者 --output）参数添加到支持的 kubectl 命令中。
输出格式	描述
-o=custom-columns=<spec>	使用逗号分隔的自定义列来打印表格
-o=custom-columns-file=<filename>	使用<filename> 文件中的自定义列模板打印表格
-o=json	输出 JSON 格式的 API 对象
-o=jsonpath=<template>	打印 jsonpath表达式中定义的字段
-o=jsonpath-file=<filename>	打印在<filename> 文件中定义的 jsonpath 表达式所指定的字段。
-o=name	仅打印资源名称而不打印其他内容
-o=wide	以纯文本格式输出额外信息，对于 Pod 来说，输出中包含了节点名称
-o=yaml	输出 YAML 格式的 API 对象

示例

#输出集群中运行着的所有镜像
kubectl get pods -A -o=custom-columns='DATA:spec.containers[*].image'

二、kubectl扩展脚本安装