一、介绍
在 Kubernetes 中，你可以为 Pod 里的容器定义一个健康检查“探针”（Probe）。这样，kubelet 就会根据这个 Probe 的返回值决定这个容器的状态，而不是直接以容器镜像是否运行（来自 Docker 返回的信息）作为依据。这种机制，是生产环境中保证应用健康存活的重要手段。

二、 三种探针类型
探针分为三种分别是StartupProbe、LivenessProbe、ReadinessProbe.
1、StartupProbe：k8s的1.16版本后新加的探测方式，用于判断容器内应用程序是否已经启动。如果配置了startupProbe，就会先禁止其他的探测，直到它成功为止，成功后将不在进行探测。注意在k8s1.17及以下版本属于新特性默认是没有开启的需要手动配置开启--feature-gates=StartupProbe=true
2、LivenessProbe：用于探测容器是否运行，如果探测失败，kubelet会根据配置的重启策略进行相应的处理。若没有配置该探针，默认就是success。
3、ReadinessProbe：一般用于探测容器内的程序是否健康，它的返回值如果为success，那么久代表这个容器已经完成启动，并且程序已经是可以接受流量的状态。

三、 探测方式
1、ExecAction：在容器内执行一个命令，如果返回值为0，则认为容器健康。
livenessProbe:  # 可选，健康检查
      exec:        # 执行容器命令检测方式
        command: 
        - cat
        - /health
2、TCPSocketAction：通过TCP连接检查容器内的端口是否是通的，如果是通的就认为容器健康。
livenessProbe:  # 可选，健康检查
      tcpSocket:    # 端口检测方式
        port: 80
3、HTTPGetAction：通过应用程序暴露的API地址来检查程序是否是正常的，如果状态码为200~400之间，则认为容器健康。(推荐使用api接口来来探测容器是否运行正常，需要开发提供相应接口)
readinessProbe:    # 可选，健康检查
      httpGet:        # httpGet检测方式
        path: /_health # 检查路径
        port: 8080
        httpHeaders:   # 检查的请求头
        - name: end-user
          value: Jason
四、 探针检查参数配置
1、initialDelaySeconds：容器启动后要等待多少秒后存活和就绪探测器才被初始化，默认是 0 秒，最小值是 0。请根据pod启动时间进行设置。
2、periodSeconds：执行探测的时间间隔（单位是秒）。默认是 10 秒。最小值是 1。不易设置过短也不宜过长，请根据实际情况设置。
3、timeoutSeconds：探测的超时后等待多少秒。默认值是 1 秒。最小值是 1。在1.20之前的版本这个参数是可以进行设置其他值在1.20之后官方修改为1。可以在你可以在所有的 kubelet 上禁用 ExecProbeTimeout。恢复之前的模式。
3、successThreshold：探测器在失败后，被视为成功的最小连续成功数。默认值是 1。 存活和启动探测的这个值必须是 1。最小值是 1。
4、failureThreshold：当探测失败时，Kubernetes 的重试次数。 存活探测情况下的放弃就意味着重新启动容器。 就绪探测情况下的放弃 Pod 会被打上未就绪的标签。默认值是 3。最小值是 1。不要设置的次数过多如果pod出问题会导致检测错误的时间变的过长。

一个探针的探测失败时间=initialDelaySeconds+((periodSeconds+timeoutSeconds)*failureThreshold)
为什么要使用StartupProbe？
如果没有startupProbe，对于启动慢的程序来说必定需要调整等待时间参数，这会导致pod从启动到就绪的时间变长，在更新服务时会导致服务就绪时间变长从而有可能导致服务不可用。如果配置StartupProbe，可以把StartupProbe时间设置较长，其余探针设置的时间短一些，由于StartupProbe执行成功后变不会在执行所以会极大的降低pod的就绪时间。在滚动升级时非常好用。
示例：
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

[root@k8smaster build-yaml]# cat test-livenes.yml
apiVersion: v1
kind: Pod
metadata:
  labels:
    test: liveness
  name: test-liveness-exec
spec:
  containers:
  - name: liveness
    image: busybox
    args:
    - /bin/sh
    - -c
    - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 5
      periodSeconds: 5
在这个 Pod 中，我们定义了一个有趣的容器。它在启动之后做的第一件事，就是在 /tmp 目录下创建了一个 healthy 文件，以此作为自己已经正常运行的标志。而 30 s 过后，它会把这个文件删除掉。
与此同时，我们定义了一个这样的 livenessProbe（健康检查）。它的类型是 exec，这意味着，它会在容器启动后，在容器里面执行一条我们指定的命令，比如：“cat /tmp/healthy”。这时，如果这个文件存在，这条命令的返回值就是 0，Pod 就会认为这个容器不仅已经启动，而且是健康的。这个健康检查，在容器启动 5 s 后开始执行（initialDelaySeconds: 5），每 5 s 执行一次（periodSeconds: 5）。
创建这个pod
[root@k8smaster build-yaml]# kubectl create -f test-livenes.yml 
pod/test-liveness-exec created

查看pod状态
[root@k8smaster build-yaml]# kubectl get pods
NAME                                READY   STATUS    RESTARTS   AGE
test-liveness-exec                  1/1     Running   0          6s

查看该pod的events
[root@k8smaster build-yaml]# kubectl describe pod test-liveness-exec
Events:
  Type     Reason     Age                        From               Message
  ----     ------     ----                       ----               -------
  Normal   Scheduled  4m22s                      default-scheduler  Successfully assigned default/test-liveness-exec to k8snode2
  Normal   Pulled     47s (x3 over 3m47s)        kubelet, k8snode2  Successfully pulled image "busybox"
  Normal   Created    47s (x3 over 3m47s)        kubelet, k8snode2  Created container liveness
  Normal   Started    46s (x3 over 3m46s)        kubelet, k8snode2  Started container liveness
  Warning  Unhealthy  5s (x9 over 3m15s)         kubelet, k8snode2  Liveness probe failed: cat: can't open '/tmp/healthy': No such file or directory
  Normal   Killing    5s (x3 over 3m5s)          kubelet, k8snode2  Container liveness failed liveness probe, will be restarted
  Normal   Pulling    <invalid> (x4 over 3m50s)  kubelet, k8snode2  Pulling image "busybox"
此时查看pod虽然是running状态，但是已经restart了3次
[root@k8smaster build-yaml]# kubectl get pods
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-6f859b4555-s26gr   1/1     Running   1          17h
nginx-deployment-6f859b4555-zcj9t   1/1     Running   1          17h
test-liveness-exec                  1/1     Running   3          4m35s
test-projected-volume               1/1     Running   0          43m
这个功能就是 Kubernetes 里的 Pod 恢复机制，也叫 restartPolicy。它是 Pod 的 Spec 部分的一个标准字段（pod.spec.restartPolicy），默认值是 Always，即：任何时候这个容器发生了异常，它一定会被重新创建。
但一定要强调的是，Pod 的恢复过程，永远都是发生在当前节点上，而不会跑到别的节点上去。事实上，一旦一个 Pod 与一个节点（Node）绑定，除非这个绑定发生了变化（pod.spec.node 字段被修改），否则它永远都不会离开这个节点。这也就意味着，如果这个宿主机宕机了，这个 Pod 也不会主动迁移到其他节点上去。
pod和deployment的区别：如果pod所在node异常，那么该pod不会被迁移到其他节点；而deployment会由K8S负责调度保障业务正常运行。
还可以通过设置 restartPolicy，改变 Pod 的恢复策略。除了 Always，它还有 OnFailure 和 Never 两种情况：
对象是pod中的容器，不是pod本身
Always：在任何情况下，只要容器不在运行状态，就自动重启容器；
OnFailure: 只在容器 异常时才自动重启容器；  对于container exit！=0时才会重启
Never: 从来不重启容器。
1、只要 Pod 的 restartPolicy 指定的策略允许重启异常的容器（比如：Always），那么这个 Pod 就会保持 Running 状态，并进行容器重启。否则，Pod 就会进入 Failed 状态 。
2、对于包含多个容器的 Pod，只有它里面所有的容器都进入异常状态后，Pod 才会进入 Failed 状态。在此之前，Pod 都是 Running 状态。此时，Pod 的 READY 字段会显示正常容器的个数，比如：

livenessprobe定义http或者tcp请求的方法
...
livenessProbe:
     httpGet:
       path: /healthz
       port: 8080
       httpHeaders:
       - name: X-Custom-Header
         value: Awesome
       initialDelaySeconds: 3
       periodSeconds: 3

    ...
    livenessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 15
      periodSeconds: 20
liveness: 失败情况下，容器会重启。 readness: 失败情况下，仅仅Service不可访问。

