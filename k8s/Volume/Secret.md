
Secret，Secret用来保存敏感信息，例如密码、OAuth 令牌和 ssh key 等等，将这些信息放在 Secret 中比放在 Pod 的定义中或者 Docker 镜像中要更加安全和灵活。

Secret 主要使用的有以下三种类型：

Opaque：base64 编码格式的 Secret，用来存储密码、密钥等；但数据也可以通过base64 –decode解码得到原始数据，所有加密性很弱。
kubernetes.io/dockerconfigjson：用来存储私有docker registry的认证信息。
kubernetes.io/service-account-token：用于被 ServiceAccount ServiceAccount 创建时 Kubernetes 会默认创建一个对应的 Secret 对象。Pod 如果使用了 ServiceAccount，对应的 Secret 会自动挂载到 Pod 目录 /run/secrets/kubernetes.io/serviceaccount 中。
bootstrap.kubernetes.io/token：用于节点接入集群的校验的 Secret


一、 Secret创建
1、 Opaque类型的Secret创建方式
1.1、使用yaml方式创建
[root@k8smaster secret]# echo -n "admin" | base64
YWRtaW4=
[root@k8smaster secret]# echo -n "root.2020" | base64
cm9vdC4yMDIw
[root@k8smaster secret]# cat secret-demo.yaml 
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: cm9vdC4yMDIw

[root@k8smaster secret]# kubectl get secret
NAME                  TYPE                                  DATA   AGE
default-token-b45m8   kubernetes.io/service-account-token   3      23d
mysecret              Opaque                                2      25s

[root@k8smaster secret]# kubectl describe secrets mysecret 
Name:         mysecret
Namespace:    default
Labels:       <none>
Annotations:  
Type:         Opaque

Data
====
password:  9 bytes
username:  5 bytes
[root@k8smaster secret]# kubectl get secrets mysecret -oyaml
apiVersion: v1
data:
  password: cm9vdC4yMDIw
  username: YWRtaW4=
kind: Secret
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"password":"cm9vdC4yMDIw","username":"YWRtaW4="},"kind":"Secret","metadata":{"annotations":{},"name"
:"mysecret","namespace":"default"},"type":"Opaque"}  creationTimestamp: "2022-09-22T19:52:49Z"
  name: mysecret
  namespace: default
  resourceVersion: "315227"
  selfLink: /api/v1/namespaces/default/secrets/mysecret
  uid: cffcd52a-4cdd-4453-900b-4869ecb9eae5
type: Opaque

1.2、 使用kubectl create secret generic命令指定文件或文件夹创建
文件里为明文，容器里也为明文，只是在secret里面使用密文
[root@k8smaster secret]# cat file/password 
kylin.2020
[root@k8smaster secret]# cat file/re.ip 
192.168.189.143:8443

指定目录创建secret
[root@k8smaster secret]# kubectl create secret generic test --from-file=file/
secret/test created
[root@k8smaster secret]# kubectl get secrets 
NAME                  TYPE                                  DATA   AGE
default-token-b45m8   kubernetes.io/service-account-token   3      23d
mysecret              Opaque                                2      63m
test                  Opaque                                2      8s
[root@k8smaster secret]# kubectl get secrets test  -oyaml
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

[root@k8smaster secret]# kubectl describe secrets test
Name:         test
Namespace:    default
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
password:  11 bytes
re.ip:     21 bytes

2、 创建用于认证镜像仓库的Secret(kubernetes.io/dockerconfigjson)
一般由kubectl create secret docker-registry命令创建。

kubectl create secret docker-registry docker-registry --docker-server=10.122.6.79:9000 --docker-username=admin --docker-password=123456 --docker-email=k8s@zhang.com

#参数说明
--docker-server  #仓库地址
--docker-username #用户名
--docker-password #密码
--docker-email    #邮箱

3、 使用Secret管理tls证书(kubernetes.io/tls)
一般由kubectl create secret tls命令创建。

#生成证书文件
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=www.zhangzhuo.com"

kubectl create secret tls tls --key=tls.key --cert=tls.crt 

#参数
--key #私钥
--cert #公钥

二、 创建好Secret对象后使用以下方式来使用它

2.1 使用Opaque类型的secret
2.1.1 #文件挂载到目录使用
        volumeMounts:
          - name: nginx-conf
            mountPath: /etc/config
      volumes:
        - name: nginx-conf
          secret:
            secretName: test
            items:
              - key: redis.config
                path: redis.conf
                mode: 0000
            defaultMode: 0666
2.1.2 #文件挂载到文件使用
        volumeMounts:
          - name: test
            mountPath: /etc/nginx/nginx.conf
            subPath: nginx.conf
      volumes:          #挂载定义
        - name: test  #挂载名称
          secret:
            secretName: nginx-conf
            items:
              - key: nginx.conf
                path: nginx.conf
2.1.3 #引用变量
        envFrom:
        - secretRef:
            name: test
        env:
        - name: zhangzhuo
          valueFrom:
            secretKeyRef:
              name: test      #cm名称
              key: zhangzhuo      #cm中key名称

volume挂载示例：
创建一个secret，
[root@k8smaster secret]# cat secret-demo.yaml 
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=        #base64编码
  password: cm9vdC4yMDIw    #echo -n "xx" | base64
[root@k8smaster secret]# cat secret-demo.yaml 
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: cm9vdC4yMDIw
使用上面创建的这个secret
[root@k8smaster secret]# cat secret-deploy.yml 
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-secret
spec:
  selector:
    matchLabels:
      app: nginx-secret
  template:
    metadata:
      labels:
        app: nginx-secret
    spec:
      containers:
      - image: nginx:1.8
        name: nginx
        volumeMounts:
        - name: secrets
          mountPath: /etc/secrets
      volumes:
      - name: secrets
        secret:
          secretName: mysecret
进入pod查看
root@nginx-secret-7d79c48444-7p2vl:/etc/secrets# cat password 
root.2020
root@nginx-secret-7d79c48444-7p2vl:/etc/secrets# cat username 
admin

2.2 使用kubernetes.io/dockerconfigjson类型的secret
如何使用，需要在Pod定义中使用imagePullSecrets配置，可以写多个他会自动匹配。

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
      imagePullSecrets:       #配置认证仓库的Secret
      - name: docker-registry #Secret名称，可以写多个
      containers:
      - image: 10.122.6.81:5000/image/nginx:v1
        name: nginx

2.3 使用tls类型的secret
一般由kubectl create secret tls命令创建。

#生成证书文件
一般的tls的secret由k8s中的ingress资源使用

apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: nginx-https-test
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
  rules:
    - host: www.zhangzhuo.com
      http:
        paths:
        - backend:
            serviceName: nginx
            servicePort: 80
  tls:          #tls证书配置
  - hosts: 
    - www.zhangzhuo.com
    - www.zhangzhuo1.com
    secretName: tls
    secretName: tls1