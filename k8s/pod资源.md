一、 pod介绍
Pod是Kubernetes中最小的单元，它由一组、一个或多个容器组成，每个Pod还包含了一个Pause容器，Pause容器是Pod的父容器，主要负责僵尸进程的回收管理，通过Pause容器可以使同一个Pod里面的多个容器共享存储、网络、PID、IPC等。

Pod在生产环境中基本不单独使用，一般都是配合控制器来使用

二、pod示例文件
apiVersion: v1  # 必选，API的版本号
kind: Pod       # 必选，类型Pod
metadata:       # 必选，元数据
  name: pod     # 必选，符合RFC 1035规范的Pod名称
  namespace: default # 可选，Pod所在的命名空间，不指定默认为default，可以使用-n 指定namespace 
  labels:       # 可选，标签选择器，一般用于过滤和区分Pod
    app: pod    #可以写多个
  annotations:  #可选，注释列表，可以写多个
    app: nginx
spec:           # 必选，用于定义容器的详细信息
  initContainers:  #初始化容器，在容器启动之前执行的一些初始化操作
  - name: init-echo
    command:
    - sh
    - -c
    - echo "I am InitContainer for init some configuration"
    image: 10.122.6.81:5000/image/centos:v1
    imagePullPolicy: IfNotPresent
  containers:    # 必选，容器列表
  - name: nginx  # 必选，符合RFC 1035规范的容器名称
    image: 10.122.6.81:5000/image/nginx:v1     # 必选，容器所用的镜像的地址
    imagePullPolicy: IfNotPresent     # 可选，镜像拉取策略, IfNotPresent: 如果宿主机有这个镜像，那就不需要拉取了. Always: 总是拉取, Never: 不管是否存储都不拉去，默认为Always
    workingDir: /usr/share/nginx/html       # 可选，容器的工作目录
    ports:  # 可选，容器需要暴露的端口号列表
    - name: http    # 端口名称
      containerPort: 80     # 端口号
      protocol: TCP # 端口协议，默认TCP
    envFrom:      #轮询定义环境变量
    - configMapRef:  #可以写多个，来源cm
        name: test   #cm名称
    env:    #可选，环境变量配置列表
    - name: TZ      # 变量名
      value: Asia/Shanghai # 变量的值
    - name: LANG
      value: en_US.utf8
    - name: test1    #变量名称
      valueFrom:     #定义变量来自哪里
        configMapKeyRef:  #定义cm
          name: test      #cm名称
          key: test1      #cm中key名称
    resources:      # 可选，资源限制和资源请求限制
      limits:       # 最大限制设置
        cpu: 200m  #1核CPU=1000m
        memory: 100Mi
      requests:     # 启动所需的资源
        cpu: 200m
        memory: 100Mi
    startupProbe: # 探测容器是否启动
      exec:        # 执行容器命令检测方式
        command: 
        - ls
        - /tmp/zhangzhuo/1/
      initialDelaySeconds: 30   # 初始化时间
      timeoutSeconds: 1     # 超时时间
      periodSeconds: 3      # 检测间隔
      successThreshold: 1 # 检查成功为1次表示就绪
      failureThreshold: 3 # 检测失败3次表示未就绪
    readinessProbe: # 可选，就绪探针。
      httpGet:      # httpGet检测方式
        path: / # 检查路径
        port: 80        # 监控端口
      initialDelaySeconds: 3
      timeoutSeconds: 1 
      periodSeconds: 3  
      successThreshold: 1 
      failureThreshold: 3 
    livenessProbe:  # 可选，存活探针
      tcpSocket:    # 端口检测方式
        port: 80
      initialDelaySeconds: 60       
      timeoutSeconds: 2     
      periodSeconds: 5      
      successThreshold: 1 
      failureThreshold: 3 
    lifecycle:    #回调配置
      postStart:  #容器启动后执行  
        exec:     #执行命令
          command:
          - sh
          - -c
          - 'mkdir /tmp/zhangzhuo/1 -p'
      preStop:    #容器终止前执行
        httpGet:  #执行http请求    
              path: /
              port: 80

  restartPolicy三种: 
  Always   # 可选，默认为Always，容器故障或者没有启动成功，那就自动重启该容器，
  Onfailure: 容器以不为0的状态终止，自动重启该容器,
  Never:无论何种状态，都不会重启，默认为Always
  nodeSelector:  #可选，指定Node节点
    type: normal #填写node标签
  terminationGracePeriodSeconds: 30 #Pod删除停止时最大容忍时间，默认30s
  dnsPolicy: ClusterFirstWithHostNet  #Pod的dns策略,Default,None,ClusterFirst,ClusterFirstWithHostNet
  hostNetwork: true #使用宿主机网络，默认不使用

三、 初始化容器
每个Pod中可以包含多个容器， 应用运行在这些容器里面，同时 Pod 也可以有一个或多个先于应用容器启动的 Init 容器。

可以定义多个init容器

Init 容器与普通的容器非常像，除了如下两点：

    它们总是运行到完成然后退出。
    每个都必须在下一个启动之前成功完成。
    pause --> InitContainer --> yml文件中定义containers

如果 Pod 的 Init 容器失败，kubelet 会不断地重启该 Init 容器直到该容器成功为止。 然而，如果 Pod 对应的 restartPolicy 值为 "Never"，并且 Pod 的 Init 容器失败， 则 Kubernetes 会将整个 Pod 状态设置为失败。

initContainers:  #初始化容器，在容器启动之前执行的一些初始化操作
  - name: init-run-sleep 
    command: ['bash','sleep','15']   #初始化容器启动后执行的命令分别有agrs，command
    image: 10.122.6.81:5000/image/centos:v1
    imagePullPolicy: IfNotPresent
  - name: init-test
    command:
    - sh
    - -c
    - 'ping -c1 test'
    image: 10.122.6.81:5000/image/centos:v1
    imagePullPolicy: IfNotPresent

command、args两项实现覆盖Dockerfile中ENTRYPOINT的功能

    具体的command命令代替ENTRYPOINT的命令行，args代表集体的参数。
    如果command和args均没有写，那么用Dockerfile的配置。
    如果command写了，但args没有写，那么Dockerfile默认的配置会被忽略，执行输入的command（不带任何参数，当然command中可自带参数）。
    如果command没写，但args写了，那么Dockerfile中配置的ENTRYPOINT的命令行会被执行，并且将args中填写的参数追加到ENTRYPOINT中。
    如果command和args都写了，那么Dockerfile的配置被忽略，执行command并追加上args参数。

四、pod退出流程
用户执行删除操作后pod所做的事情
pod退出流程
用户执行删除操作-->pod变成terminated-->endpoints会删除该pod的ip地址-->如果pod还有进程在运行，会默认预留一定的时间，默认为30秒-->执行prestop-->执行终止容器的信号
最后两个操作必须在预留时间内执行完成如果没有执行完成k8s会直接杀死容器
注意：在上面的流程中除去执行preStop外，其余必定都会执行。
4.1 容器回调
有两个回调暴露给容器：
PostStart

这个回调在容器被创建之后立即被执行。 但是，不能保证回调会在容器入口点（ENTRYPOINT）之前执行。 没有参数传递给处理程序。基本很少使用。
PreStop

在容器因 API 请求或者管理事件（诸如存活态探针、启动探针失败、资源抢占、资源竞争等） 而被终止之前，此回调会被调用。 在用来停止容器的 TERM 信号被发出之前，回调必须执行结束。 Pod 的终止宽限周期在 PreStop 回调被执行之前即开始计数，所以无论回调函数的执行结果如何，容器最终都会在 Pod 的终止宽限期内被终止。
  terminationGracePeriodSeconds:容器的宽限时间
回调方式
    Exec - 在容器执行给定的命令或脚本，重点使用。
    HTTP - 对容器上的特定端点执行 HTTP 请求。
配置示例：
lifecycle:    #回调配置
      postStart:  #容器启动后执行  
        exec:     #执行命令
          command:
          - sh
          - -c
          - 'mkdir /data/ '
      preStop:    #容器终止前执行
        httpGet:  #执行http请求    
              path: /
              port: 80


