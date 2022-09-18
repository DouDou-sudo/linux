一、ConfigMap

官方文档：https://kubernetes.io/zh/docs/concepts/configuration/configmap/
1.1 ConfigMap介绍

在微服务架构中，大多数的服务的配置文件都是与服务本身分开的，由统一的配置中心进行管理，服务启动后会到配置中心读取自己的配置文件之后启动服务。Kubernetes中也提供了一个配置文件的api就是configmap。

ConfigMap 将您的环境配置信息和容器镜像解耦，便于应用配置的修改。

注意：ConfigMap 并不提供保密或者加密功能。 如果你想存储的数据是机密的，请使用Secret， 或者使用其他第三方工具来保证你的数据的私密性，而不是用 ConfigMap。

ConfigMap 在设计上不是用来保存大量数据的。在 ConfigMap 中保存的数据不可超过 1 MiB。如果你需要保存超出此尺寸限制的数据，你可能希望考虑挂载存储卷或者使用独立的数据库或者文件服务。

1.2 ConfigMap的创建

创建方式

    可以编写yaml文件，使用kubectl命令指定yaml文件进行创建
    可以直接使用kubectl create cm命令指定文件或文件夹进行创建

kubectl create cm参数

    --from-file：以本地文件的方式创建cm的配置文件，可以指定一个文件或多个文件或一个文件夹。格式--from-file=定义到mc中的key名称=文件位置，如果不写名称直接指定文件默认会把文件名称作为key名称。不能与--from-env-file一起用。
    --from-env-file：以本地文件的方式创建cm的变量，可以指定多个。
    --from-literal：创建变量的cm，格式--from-literal=变量名称=变量值,可以写多个
1.示例yaml文件

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

2.以文件或文件夹的方式创建configmap的配置文件
#查看文件夹下的文件
[root@k8s test]# ls configmap/
nginx.config  redis.config

以文件的方式创建
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
#多个文件
kubectl create configmap test --from-file=redis=configmap/redis.config --from-file=nginx=configmap/nginx.config
#验证
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

以文件夹的方式创建

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

3.以文件的方式创建环境变量的configmap

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

1.3 ConfigMap的使用

首先创建一个configmap，示例如下

apiVersion: v1
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: test
data:
  nginx.config: |
    nginx true
    password 123456
  redis.config: |
    redis true
    password 123456
  test1: zz
  test2: zhanghuo

1.以环境变量的方式使用ConfigMap

Pod使用cm需要与cm在同一个命名空间，否则是无法调用的，会提示找不到。

手动指定环境变量，单个引用

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
        env:             #变量定义
        - name: test1    #定义到容器中的变量名称
          valueFrom:     #定义变量来自哪里
            configMapKeyRef:  #定义cm
              name: test      #cm名称
              key: test1      #cm中key名称
        - name: test2
          valueFrom:
            configMapKeyRef:
              name: test
              key: test2
#验证
kubectl exec -it nginx-1-6dc7b8574-8862z -- env | grep -e test1 -e test2
test1=zz
test2=zhanghuo

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
        volumeMounts:       #容器挂载定义
          - name: nginx-conf  #定义的volumes的名称
            mountPath: /etc/config  #挂载点，目录
      volumes:          #挂载定义
        - name: nginx-conf  #挂载名称
          configMap:        #挂载的资源
            name: test      #cm名称
#验证
kubectl exec -it nginx-1-777cb9459f-bd7t9 -- ls /etc/config
nginx.config  redis.config  test1  test2

#挂载的文件修改cm内容会自动刷新，同步会有一定的时间
kubectl edit cm test
  nginx.config: |
    nginx true-true
    password 12345
#验证
kubectl exec -it nginx-1-777cb9459f-bd7t9 -- cat /etc/config/nginx.config
nginx true-true
password 12345

挂载configmap中其中的一个配置到文件夹

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
        volumeMounts:
          - name: nginx-conf
            mountPath: /etc/config 
      volumes:
        - name: nginx-conf
          configMap:
            name: test 
            items:  #具体定义要挂载的内容
              - key: redis.config  #cm的key名称
                path: redis.conf   #挂载后的文件名称
                mode: 0000         #文件权限，要比defaultMode优先级高
            defaultMode: 0666      #文件权限作用与整个cm
#验证
kubectl exec -it nginx-1-6b7bcc77cc-n2f7k -- ls /etc/config
redis.conf
root@nginx-1-74bb6cfdfd-fqw6q:/# cd /etc/config/
root@nginx-1-74bb6cfdfd-fqw6q:/etc/config# ls -l
total 0
lrwxrwxrwx 1 root root 17 Jan  6 05:40 redis.conf -> ..data/redis.conf
root@nginx-1-74bb6cfdfd-fqw6q:/etc/config# ls -l ..data/redis.conf 
-rw-rw-rw- 1 root root 27 Jan  6 05:40 ..data/redis.conf

挂载一个配置文件到文件，避免目录覆盖

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
        volumeMounts:
          - name: test
            mountPath: /etc/nginx/nginx.conf  #挂载路径写全部，到挂载的文件
            subPath: nginx.conf  #挂载为那个文件
      volumes:          
        - name: test  
          configMap:   
            name: nginx-conf
            items:
              - key: nginx.conf
                path: nginx.conf
