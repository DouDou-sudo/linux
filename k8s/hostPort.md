hostnetwork，pod使用的是宿主机的网络地址，即pod ip就是宿主机的ip，pod直接使用宿主机的网络根namespace
hostport就是暴露pod所在节点ip+port给外部访问，使用DNAT机制将hostport指定的端口映射到容器的端口之上

1、hostnetwork:
[root@k8smaster hostport]# cat hostport.yml 
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx-hostport
  name: nginx-hostport
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-hostport
  template:
    metadata:
      labels:
        app: nginx-hostport
    spec:
      containers:
      - image: nginx:1.8
        name: nginx
      #  ports:
      #  - containerPort: 80
      #    hostPort: 8000
      #    name: http
      #    protocol: TCP
      hostNetwork: true
[root@k8smaster hostport]# kubectl get pods nginx-hostport-56b4fb8bcf-54c7h -owide
NAME                              READY   STATUS    RESTARTS   AGE    IP                NODE       NOMINATED NODE   READINESS GATES
nginx-hostport-56b4fb8bcf-54c7h   1/1     Running   0          118s   192.168.189.201   k8snode1   <none>           <none>
访问测试:
[root@k8smaster hostport]# curl 192.168.189.201
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
进入pod查看网络，看到的就是宿主机的网络
[root@k8smaster hostport]# kubectl exec -it nginx-hostport-56b4fb8bcf-54c7h -- /bin/bash
root@k8snode1:/# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0c:29:a7:24:f5 brd ff:ff:ff:ff:ff:ff
    inet 192.168.189.201/24 brd 192.168.189.255 scope global noprefixroute ens33
       valid_lft forever preferred_lft forever
    inet6 fe80::ac64:b26:e745:a829/64 scope link noprefixroute 
       valid_lft forever preferred_lft forever
3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default 
    link/ether 02:42:92:0e:bc:6a brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
       valid_lft forever preferred_lft forever
4: flannel.1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1450 qdisc noqueue state UNKNOWN group default 
    link/ether 3e:01:fe:74:63:6c brd ff:ff:ff:ff:ff:ff
    inet 10.244.1.0/32 scope global flannel.1
       valid_lft forever preferred_lft forever
    inet6 fe80::3c01:feff:fe74:636c/64 scope link 
       valid_lft forever preferred_lft forever

2、hostport:
[root@k8smaster hostport]# cat hostport.yml 
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx-hostport
  name: nginx-hostport
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-hostport
  template:
    metadata:
      labels:
        app: nginx-hostport
    spec:
      containers:
      - image: nginx:1.8
        name: nginx
        ports:
        - containerPort: 80
          hostPort: 8000
          name: http
          protocol: TCP

[root@k8smaster hostport]# kubectl get pod nginx-hostport-65ddc56dd5-qzrdg  -owide
NAME                              READY   STATUS    RESTARTS   AGE     IP             NODE       NOMINATED NODE   READINESS GATES
nginx-hostport-65ddc56dd5-qzrdg   1/1     Running   0          2m34s   10.244.1.132   k8snode1   <none>           <none>
访问测试
通过pod ip+containerport访问
[root@k8smaster hostport]# curl 10.244.1.132
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
通过node节点ip+hostport访问
[root@k8smaster hostport]# curl 192.168.189.201:8000
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
