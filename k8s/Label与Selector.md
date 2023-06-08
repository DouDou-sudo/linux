Label：对k8s中各种资源进行分类、分组，添加一个具有特别属性的一个标签。

Selector：通过一个过滤的语法进行查找到对应标签的资源
一、Label & Selector
当Kubernetes对系统的任何API对象如Pod和节点进行“分组”时，会对其添加Label（key=value形式的“键-值对”）用以精准地选择对应的API对象。而Selector（标签选择器）则是针对匹配对象的查询方法。注：键-值对就是key-value pair。

例如，常用的标签tier可用于区分容器的属性，如frontend、backend；或者一个release_track用于区分容器的环境，如canary、production等。
1.1 Label的介绍
label可以给k8s中大多数资源进行标签的定义，主要作用为用于指定对用户有意义且相关的对象的标识属性，但并不对资源工作参数任何影响。

标签是键值对。有效的标签键有两个段：可选的前缀和名称，用斜杠（/）分隔。 名称段是必需的，必须小于等于 63 个字符，以字母数字字符（[a-z0-9A-Z]）开头和结尾， 带有破折号（-），下划线（_），点（ .）和之间的字母数字。

kubernetes.io/ 和 k8s.io/ 前缀是为 Kubernetes 核心组件保留的。
1.2 标签的使用

添加`label`

大多数kubernetes中的资源都是可以进行打标签的，命令是kubectl label。
可以给指定资源添加多个标签
#给指定的资源添加label
kubectl label deployments.apps nginx-1 app=nginx 

#给某一类型的所有资源打标签
kubectl label deployments.apps app=test --all
删除label

#删除指定资源的label
kubectl label deployments.apps nginx-1 app-

#删除某一类型所有资源的标签
kubectl label deployments.apps app- --all
修改标签

kubectl label deployments.apps nginx-1 app=nginx-1 --overwrite

查看资源标签

kubectl get deployments.apps --show-labels

1.3 利用label查找资源
kubectl get可以打印资源列表，并且可以使用-l参数利用label进行资源筛选
```
    ==,=:等于
    !=:不等于
    in:包含
    notin:不包含
```
#一个条件

kubectl get po -l app=nginx
kubectl get po -l app!=nginx

#多个label的and关系
kubectl get po -l app=nginx,app=nginx-1

#多个label的in关系
kubectl get po -l 'app in(nginx,nginx-1)'

#多个label的notin关系
kubectl get po -l 'app notin(nginx,nginx-1)'

查看labels
[root@k8smaster DaemonSet]# kubectl get deployments.apps nginx-deployment  --show-labels  
NAME               READY   UP-TO-DATE   AVAILABLE   AGE   LABELS
nginx-deployment   3/3     3            3           19d   app=nginx-1,da=test
查找app=nginx的pod，此查找不是严格匹配
[root@k8smaster DaemonSet]# kubectl get pods -l app=nginx
NAME                               READY   STATUS    RESTARTS   AGE
nginx-deployment-bb4d88ddf-fk28t   1/1     Running   0          77m
nginx-deployment-bb4d88ddf-xp8v2   1/1     Running   0          77m
nginx-deployment-bb4d88ddf-xt255   1/1     Running   0          78m
