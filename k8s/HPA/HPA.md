Horizontal Pod Autoscaler(HAP),是一个Pod 水平自动扩缩（Horizontal Pod Autoscaler） 可以基于 CPU 利用率自动扩缩 ReplicationController、Deployment、ReplicaSet 和 StatefulSet 中的 Pod 数量。 除了 CPU 利用率，也可以基于其他应程序提供的自定义度量指标来执行自动扩缩。 Pod 自动扩缩不适用于无法扩缩的对象，比如 DaemonSet。

Pod 水平自动扩缩特性由 Kubernetes API 资源和控制器实现。资源决定了控制器的行为。 控制器会周期性地调整副本控制器或 Deployment 中的副本数量，以使得类似 Pod 平均 CPU 利用率、平均内存利用率这类观测到的度量值与用户所设定的目标值匹配。

一、HAP自动扩缩容

HAP的api版本有三个分别是：

    HPA v1为稳定版自动水平伸缩，只支持CPU指标，需要安装metrics-server
    V2为beta版本，分为v2beta1(支持CPU、内存和自定义指标)
    v2beta2(支持CPU、内存、自定义指标Custom和额外指标ExternalMetrics)

1.1 使用v1版本的cpu利用率进行扩缩容

前置条件

    只能对副本控制器使用(deployment ,replicaset,replicationcontroller,statefulset)
    需要安装metrics-server采集pod监控指标,集群可以使用kubectl top命令
    Pod资源定义中必须启用资源限制resources.limits
    控制器选择器中选择器必须匹配到的pod只能一个，不能是多个，不然会导致HPA控制器获取不到监控指标，HAP无法工作。

示例文件
apiVersion: autoscaling/v1      #必须，api版本
kind: HorizontalPodAutoscaler   #必须，资源类型
metadata:                       #必须，元数据定义
  name: nginx-1                 #必须，名称
  namespace: default            #可选，命名空间
spec:                           #实际内容定义              
  maxReplicas: 10               #最大pod数
  minReplicas: 1                #最小pod数           
  scaleTargetRef:               #选择匹配的控制器设置
    apiVersion: apps/v1         #控制器的api版本
    kind: Deployment            #控制器类型
    name: nginx-1               #控制器名称
  targetCPUUtilizationPercentage: 10  #扩缩容条件CPU负载超过的值百分比
创建命令

一般HAP的v1版本，创建的参数较少，可以直接使用命令创建。

kubectl autoscale 控制器类型 控制器名称 --min=最小pod数 --max=最大pod数 --cpu-percent=CPU负载阈值百分比

#示例
kubectl autoscale deployment nginx-1 --min=1 --max=10 --cpu-percent=10

验证

[root@km1-81 test]# kubectl get hpa 
NAME      REFERENCE            TARGETS      MINPODS      MAXPODS   REPLICAS   AGE
nginx-1   Deployment/nginx-1   0%/10%       1            10        1          8m9s
名称       选择的控制器           当前CPU负载    最小pod数     最大pod数  当前运行pod数

#如果负载值一直为空，请检查控制器是否匹配多个pod。

[root@k8smaster DaemonSet]# kubectl autoscale deployment nginx-deployment --min=1 --max=10 --cpu-percent=10
horizontalpodautoscaler.autoscaling/nginx-deployment autoscaled
[root@k8smaster DaemonSet]# kubectl get hpa
NAME               REFERENCE                     TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
nginx-deployment   Deployment/nginx-deployment   0%/10%   1         10        0          5s

注意:有些场景是不适用CPU扩缩容副本
前端服务--》后端服务--》数据库
大多数微服务架构都是类似的结构图，都是前端接收请求，之后交予后端进行处理后写入数据库。如果数据库的负载过高导致后端服务的CPU负载增高，我们在后端使用了hpa这会导致后端增加副本数，但是根本原因并不是后端服务的问题，增加副本数反而会导致数据库负载继续增加，又可能导致数据库最后宕机。

所以HPA的v1版本并不可以随便使用，要根据实际情况来进行使用。

