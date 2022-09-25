ConfigMap

官方文档：https://kubernetes.io/zh/docs/concepts/configuration/configmap/
1.1 ConfigMap介绍

在微服务架构中，大多数的服务的配置文件都是与服务本身分开的，由统一的配置中心进行管理，服务启动后会到配置中心读取自己的配置文件之后启动服务。Kubernetes中也提供了一个配置文件的api就是configmap。

ConfigMap 将您的环境配置信息和容器镜像解耦，便于应用配置的修改。

注意：ConfigMap 并不提供保密或者加密功能。 如果你想存储的数据是机密的，请使用Secret， 或者使用其他第三方工具来保证你的数据的私密性，而不是用 ConfigMap。

ConfigMap 在设计上不是用来保存大量数据的。在 ConfigMap 中保存的数据不可超过 1 MiB。如果你需要保存超出此尺寸限制的数据，你可能希望考虑挂载存储卷或者使用独立的数据库或者文件服务。

一、 ConfigMap的创建

创建方式

    可以编写yaml文件，使用kubectl命令指定yaml文件进行创建
    可以直接使用kubectl create cm命令指定文件或文件夹进行创建

kubectl create cm参数

    --from-file：以本地文件的方式创建cm的配置文件，可以指定一个文件或多个文件或一个文件夹。格式--from-file=定义到mc中的key名称=文件位置，如果不写名称直接指定文件默认会把文件名称作为key名称。不能与--from-env-file一起用。
    --from-env-file：以本地文件的方式创建cm的变量，可以指定多个。
    --from-literal：创建变量的cm，格式--from-literal=变量名称=变量值,可以写多个
1.以yaml文件的方式

    apiVersion: v1        #api版本
    kind: ConfigMap       #资源类型
    metadata:             #元数据定义
      name: test          #名称
      namespace: default  #命名空间
    data:                 #具体配置定义，配置定义的类型有俩种
      player_initial_lives: "3"   #这种以容器环境变量的方式使用
      game.properties: |          #这种以映射到容器中为配置文件的方式使用
        enemy.types=aliens,monsters
        player.maximum-lives=5

2.使用命令指定文件或文件夹的方式创建configmap的配置文件
#查看文件夹下的文件

    [root@k8s test]# ls configmap/
    nginx.config  redis.config

    指定文件的方式创建
    单个文件

    kubectl create configmap test --from-file=redis=configmap/redis.config  
    #验证
    [root@k8s test]# kubectl get cm test -oyaml
    apiVersion: v1
    data:
      redis: |
        redis true
        password 123456
    kind: ConfigMap
    metadata:
      creationTimestamp: "2022-01-06T02:57:23Z"
      name: test
      namespace: default
      resourceVersion: "57247020"
      selfLink: /api/v1/namespaces/default/configmaps/test
      uid: 0a6ab496-f932-4a36-be25-b82aeb20cd84

多个文件

   kubectl create configmap test --from-file=redis=configmap/redis.config --from-file=nginx=configmap/nginx.config

验证

   [root@k8s test]# kubectl get cm test -oyaml
   apiVersion: v1
    data:
      nginx: |
        nginx true
        password 123456
      redis: |
        redis true
        password 123456
    kind: ConfigMap
    metadata:
      creationTimestamp: "2022-01-06T02:58:35Z"
      name: test
      namespace: default
      resourceVersion: "57247237"
      selfLink: /api/v1/namespaces/default/configmaps/test
      uid: 605febe6-b50e-45d5-b8ea-fee042279320

指定文件夹的方式创建

    kubectl create configmap test --from-file=configmap/
    
#验证

    [root@k8s test]# kubectl get cm test -oyaml
    apiVersion: v1
    data:      #会把目录下的所有文件定义到mc中
      nginx.config: |    #key名称为文件名称，值为文件内容
        nginx true
        password 123456
      redis.config: |
        redis true
        password 123456
    kind: ConfigMap
    metadata:
      creationTimestamp: "2022-01-06T03:05:20Z"
      name: test
      namespace: default
      resourceVersion: "57248456"
      selfLink: /api/v1/namespaces/default/configmaps/test
      uid: b68d1351-b70c-436b-ad5c-9f7610fce0bd

3.使用命令指定文件的方式创建环境变量的configmap

#创建变量文件

    [root@k8s configmap]# cat env 
    name=zz
    name1=zhnagzhuo
    name2=zhang
    name3=zhuo
    #创建cm
    kubectl create cm test --from-env-file=env
    #验证

    [root@k8s configmap]# kubectl get cm test -oyaml
    apiVersion: v1
    data:
      name: zz
      name1: zhnagzhuo
      name2: zhang
      name3: zhuo
    kind: ConfigMap
    metadata:
      creationTimestamp: "2022-01-06T03:13:47Z"
      name: test
      namespace: default
      resourceVersion: "57249972"
      selfLink: /api/v1/namespaces/default/configmaps/test
      uid: 12f35011-0593-47de-a1ff-67484917c437

二、 ConfigMap的使用

首先创建一个configmap，示例如下

    [root@k8smaster configmap]# cat test-configmap.yml 
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: test
    data:
      nginx.config:
        nginx true
        password 123456
      redis.config:
        redis true
        password 123456
      test1: zz
      test2: bxw
    [root@k8smaster configmap]# kubectl get cm test -oyaml
    apiVersion: v1
    data:
      nginx.config: nginx true password 123456
      redis.config: redis true password 123456
      test1: zz
      test2: bxw
    kind: ConfigMap
    metadata:
      annotations:
        kubectl.kubernetes.io/last-applied-configuration: |
          {"apiVersion":"v1","data":{"nginx.config":"nginx true password 123456","redis.config":"redis true password 123456","test1":"zz"
    ,"test2":"bxw"},"kind":"ConfigMap","metadata":{"annotations":{},"name":"test","namespace":"default"}}  creationTimestamp: "2022-09-22T11:48:37Z"
      name: test
      namespace: default
      resourceVersion: "272694"
      selfLink: /api/v1/namespaces/default/configmaps/test
      uid: 93a33cd7-074e-4150-abb0-f78abf772221
    [root@k8smaster configmap]# cat test-configmap.yml 
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: test
    data:
      nginx.config:
        nginx true
        password 123456
      redis.config:
        redis true
        password 123456
      test1: zz
      test2: bxw

1.以环境变量的方式使用ConfigMap

Pod使用cm需要与cm在同一个命名空间，否则是无法调用的，会提示找不到。

手动指定环境变量，单个引用

    [root@k8smaster configmap]# cat configmap-deploy.yml 
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx-configmap
    spec:
      selector:
        matchLabels:
            app: nginx-config
      template:
        metadata:
          labels:
            app: nginx-config
        spec:
          containers:
          - image: nginx:1.8
            name: nginx
            env:
            - name: test1
              valueFrom:
                configMapKeyRef:
                  name: test
                  key: test1
            - name: test2
              valueFrom:
                configMapKeyRef:
                  name: test
                  key: test2
#验证

    [root@k8smaster configmap]# kubectl exec  nginx-configmap-f47cd58bd-ncjzn -- env | grep -e test1 -e test2
    test1=zz
    test2=bxw

定义一次定义多个环境变量

这种定义会轮询cm中所有的数据，把key作为变量的名称。

    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx-1
      namespace: default
    spec:
      selector:
        matchLabels:
          app: nginx-1
      template:
        metadata:
          labels:
            app: nginx-1
        spec:
          containers:
          - image: 10.122.6.81:5000/image/nginx:v1
            name: nginx
            envFrom:      #轮询定义环境变量
            - configMapRef:  #可以写多个，来源cm
                name: test   #cm名称
#验证

    kubectl exec -it nginx-1-5f4f8cc5b4-9s9gc -- env | grep -e test1 -e test2 -e 
    test1=zz
    test2=zhanghuo

2.以配置文件的形式挂载ConfigMap

挂载configmap到文件夹中

这种挂载方式会把configmap的所有内容以文件的方式挂载到所定义的挂载点，挂载点如果有内容会进行覆盖。文件名称会定义为cm中的key名称，内容为key的value。

    [root@k8smaster configmap]# cat configmap-deploy.yml 
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx-configmap
    spec:
      selector:
        matchLabels:
            app: nginx-config
      template:
        metadata:
          labels:
            app: nginx-config
        spec:
          containers:
          - image: nginx:1.8
            name: nginx
            volumeMounts:
              - name: nginx-conf
                mountPath: /etc/config
          volumes:
            - name: nginx-conf
              configMap:
                name: test
#验证

    [root@k8smaster configmap]# kubectl exec -it nginx-configmap-57d694b45-6mr6n -- ls /etc/config
    nginx.config  redis.config  test1  test2
    [root@k8smaster configmap]# kubectl exec -it nginx-configmap-57d694b45-6mr6n -- cat /etc/config/nginx.config
    nginx true password 123456
    [root@k8smaster configmap]# kubectl exec -it nginx-configmap-57d694b45-6mr6n -- cat /etc/config/redis.confi
    redis true password 123456
修改test-configmap.yml的password为12345重新apply
#挂载的文件修改cm内容会自动刷新，同步会有一定的时间
或者直接编辑

    kubectl edit cm test
      nginx.config: |
        nginx true-true
        password 12345
#验证

    kubectl exec -it nginx-1-777cb9459f-bd7t9 -- cat /etc/config/nginx.config
    nginx true-true
    password 12345

挂载configmap中其中的一个配置到文件夹

    [root@k8smaster configmap]# cat configmap-deploy.yml 
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx-configmap
    spec:
      selector:
        matchLabels:
            app: nginx-config
      template:
        metadata:
          labels:
            app: nginx-config
        spec:
          containers:
          - image: nginx:1.8
            name: nginx
            volumeMounts:
              - name: nginx-conf
                mountPath: /etc/config
          volumes:
            - name: nginx-conf
              configMap:
                name: test
                items:
                  - key: redis.config
                    path: redis.conf
                    mode: 0777
                defaultMode: 0666
#验证

    [root@k8smaster configmap]# kubectl exec -it nginx-configmap-74df796678-cm7cc -- ls -la  /etc/config
    total 0
    drwxrwxrwx 3 root root 77 Sep 22 19:43 .
    drwxr-xr-x 1 root root 20 Sep 22 19:43 ..
    drwxr-xr-x 2 root root 24 Sep 22 19:43 ..2022_09_22_19_43_48.007213378
    lrwxrwxrwx 1 root root 31 Sep 22 19:43 ..data -> ..2022_09_22_19_43_48.007213378
    lrwxrwxrwx 1 root root 17 Sep 22 19:43 redis.conf -> ..data/redis.conf

    root@nginx-configmap-5767df6b97-jpn49:/etc/config# ls -la ..data/redis.conf 
    ---------- 1 root root 25 Sep 22 19:42 ..data/redis.conf

挂载一个配置文件到文件，避免目录覆盖

    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx-configmap
    spec:
      selector:
        matchLabels:
            app: nginx-config
      template:
        metadata:
          labels:
            app: nginx-config
        spec:
          containers:
          - image: nginx:1.8
            name: nginx
            volumeMounts:
              - name: nginx-conf
                mountPath: /etc/config/nginx.conf
          volumes:
            - name: nginx-conf
              configMap:
                name: test
                items:
                  - key: nginx.config
                    path: nginx.conf
                    mode: 0777
                defaultMode: 0666

    root@nginx-configmap-667b46c999-f7mhq:/etc/config# ls -la
    total 0
    drwxr-xr-x 3 root root 24 Sep 22 19:51 .
    drwxr-xr-x 1 root root 20 Sep 22 19:51 ..
    drwxrwxrwx 2 root root  6 Sep 22 19:51 nginx.conf
