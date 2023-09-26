一、kubernetes的服务调用

服务的访问分为俩种形式分别是服务之间的调用(南北流量)，用户流量的访问(东西流量)，k8s中提供的南北流量的解决方法是使用ingress，东西流量的解决方案是service。

二、Service资源介绍
Service主要用于Pod之间的通信，由于Pod是一种临时资源可能随时会被调度重建，重建后Pod的IP地址也会进行变化，由于Pod的IP地址不确定性，我们无法使用Pod的IP地址来进行服务的访问，所以k8s中加入了一个service资源用来解决Pod的访问。Service一般会通过选择器选择一个或一组Pod，之后通过iptables或者ipvs的方式进行代理，service的请求会被转发到自己所代理的Pod。service资源创建后只要不进行修改他的IP地址就不会变化相对来说他的IP地址是固定的，k8s中还引用了dns组件用来解析service资源的名称得到他的IP地址，所以集群中访问service，可以直接通过service的名称就是访问service了。

2.1 集群内部访问service
创建一个service
[root@k8smaster service]# cat hostname-svc.yml 
apiVersion: v1
kind: Service
metadata:
  name: hostname
spec:
  selector:
    app: hostnames
  ports:
  - name: default
    protocol: TCP
    port: 80
    targetPort: 9376
创建service匹配的pod
[root@k8smaster service]# cat hostname.yml 
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hostname
spec:
  selector:
    matchLabels:
      app: hostnames
  replicas: 3
  template:
    metadata:
      labels:
        app: hostnames
    spec:
      containers:
      - name: hostnames
        image: mirrorgooglecontainers/serve_hostname
        ports:
        - containerPort: 9376
          protocol: TCP
查看创建的svc
[root@k8smaster service]# kubectl get svc
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
hostname     ClusterIP   10.1.115.75    <none>        80/TCP    48m
查看创建的pod
[root@k8smaster service]# kubectl get pods
NAME                               READY   STATUS    RESTARTS   AGE
hostname-69cd56f6d5-cgnrx          1/1     Running   0          37m
hostname-69cd56f6d5-rlpng          1/1     Running   0          39m
hostname-69cd56f6d5-vjd95          1/1     Running   0          37m
查看生成的endpoints
[root@k8smaster service]# kubectl get ep
NAME         ENDPOINTS                                            AGE
hostname     10.244.1.54:9376,10.244.1.55:9376,10.244.2.53:9376   48m
此时在宿主机上就可以访问cluster-ip了（此处为10.1.115.75）
[root@k8smaster service]# curl 10.1.115.75
hostname-69cd56f6d5-vjd95
[root@k8smaster service]# curl 10.1.115.75
hostname-69cd56f6d5-cgnrx
[root@k8smaster service]# curl 10.1.115.75
hostname-69cd56f6d5-rlpng

2.1 如何通过service的名称访问service
由于service资源是有命名空间隔离性的，所以不同命名空间可以创建相同名称的service。

如果是相同命名空间的Pod调用service可以直接使用service名称即可,如下
[root@k8s ~]# kubectl get svc nginx-1
NAME      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
nginx-1   ClusterIP   172.20.173.132   <none>        80/TCP    18h
#进入pod访问
/ # ping nginx-1
PING nginx-1 (172.20.173.132): 56 data bytes

如果Pod调用不同命名空间的service，如下：
[root@k8s ~]# kubectl get svc -n engage 
NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                          AGE
mysql                 NodePort    172.20.119.154   <none>        3306:41344/TCP                   119d
#进入Pod访问
/ # ping mysql.engage
PING mysql.engage (172.20.119.154): 56 data bytes

kubernetes中一个service的完整域名
#进入Pod使用nslookup命令检查
/ # nslookup nginx-1
Name:   nginx-1.default.svc.cluster.local
Address: 172.20.173.132

这里可以看到一个nginx-1名称的service的完整域名为nginx-1.default.svc.cluster.local具体含义如下：
service名称.命名空间.svc.集群部署时定义的Domain名称

service名称解析的相关内容

    k8s想要解析service的名称需要部署coredns，部署完成coredns插件后会生成一个kube-dns的service，他的地址需要手动配置正常情况下是service地址池的第二个地址
    kubelet中需要配置clusterDNS的参数为kube-dns的IP地址，clusterDomain根据安装时的规划配置：默认为cluster.local。
    之后创建Pod内部会自动把Pod的dns指向kube-dns的IP地址，这样就可以让Pod解析service的名称了

2.2 service常见的类型

    ClusterIP：在集群内部使用，也是默认值。
    ExternalName：通过返回定义的CNAME别名。
    NodePort：在所有安装了kube-proxy的节点上打开一个端口，此端口可以代理至后端Pod，然后集群外部可以使用节点的IP地址和NodePort的端口号访问到集群Pod的服务。NodePort端口范围默认是30000-32767。
    LoadBalancer：使用云提供商的负载均衡器公开服务。

三、service使用
3.1 service文件示例
apiVersion: v1     #必须，api版本
kind: Service      #必须，资源类型
metadata:          #必须，元数据定义     
  name: nginx-1    #必须，名称
  namespace: default  #可选命名空间
spec:              #具体定义信息
  ports:           #service端口定义
  - name: 80-80    #端口名称
    port: 80       #svc自己的端口
    protocol: TCP  #网络协议，UDP TCP SCTP，默认为TCP
    targetPort: 80 #后端应用端口
  selector:        #选择器，用来选择Pod
    app: nginx-1   #Pod的label
  type: ClusterIP  #svc类型，默认为ClusterIP

3.3 使用service代理k8s外部服务
如果需要使用service代理外部服务，需要手动创建service与endpoints示例如下。只要endpoints的名称与service的名称一致他们就会自动建立连接。
service示例yaml文件

apiVersion: v1   
kind: Service      
metadata:             
  name: minio   
  namespace: default  
spec:             
  ports:           
  - name: http    
    port: 9000      
    protocol: TCP  
    targetPort: 9000 
  type: ClusterIP

endpoints示例yaml文件

apiVersion: v1    #必须，api版本
kind: Endpoints   #资源类型
metadata:         #元数据定义
  name: minio     #名称，必须与service名称一致
  namespace: default
subsets:          #后端服务定义
- addresses:      #后端服务的IP定义,列表可以写多个
  - ip: 10.28.88.9 
  - ip: 10.28.88.10
  - ip: 10.28.88.11
  - ip: 10.28.88.12
  ports:          #后端服务的端口定义
  - name: http    #名称必须与service的端口定义一致
    port: 9000    #端口
    protocol: TCP #协议

创建后验证
#验证svc
[root@k8s test]# kubectl get svc minio
NAME    TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
minio   ClusterIP   172.20.226.170   <none>        9000/TCP   14s
#验证ep
[root@k8s test]# kubectl get ep minio
NAME    ENDPOINTS                                                        AGE
minio   10.28.88.10:9000,10.28.88.11:9000,10.28.88.12:9000 + 1 more...   2s
#访问测试
[root@graylog ~]# curl http://172.20.226.170:9000 
<?xml version="1.0" encoding="UTF-8"?>

3.4 使用service反代外部域名

示例文件

apiVersion: v1     #必须，api版本
kind: Service      #必须，资源类型
metadata:          #必须，元数据定义     
  name: dep   #必须，名称
  namespace: default  #可选命名空间
spec:              #具体定义信息
  type: ExternalName
  externalName: www.dev.pcep.cloud

创建测试

[root@km1-81 test]# kubectl get svc dep
NAME   TYPE           CLUSTER-IP   EXTERNAL-IP          PORT(S)   AGE
dep    ExternalName   <none>       www.dev.pcep.cloud   <none>    5s
#测试
/ # curl www.dev.pcep.cloud
{"message":"no route and no API found with those values"}

3.5 NodePort
[root@k8smaster service]# cat node-svc.yml 
apiVersion: v1
kind: Service
metadata:
  name: my-nginx
  labels:
    run: my-nginx
spec:
  type: NodePort
  ports:
  - port: 8080
    nodePort: 30000
    targetPort: 80
    protocol: TCP
    name: http
  - port: 443
    protocol: TCP
    name: https
  selector:
    run: my-nginx

[root@k8smaster service]# kubectl get svc
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)                        AGEd
my-nginx     NodePort    10.1.246.166   <none>        8080:30000/TCP,443:31446/TCP   5m7s

nodeport的端口号可以指定也可以不指定，指定的话范围必须在30000-32767内

pod 80--》svc 8080 --》node 30000