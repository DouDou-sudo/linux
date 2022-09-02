PV描述的是持久化存储数据卷，这个 API 对象主要定义的是一个持久化存储在宿主机上的目录，比如一个 NFS 的挂载目录。
PVC描述的是pod所希望的持久化存储的属性，比如，volume存储的大小，可读权限等
定义一个NFS类型的PV，
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs
spec:
  storageClassName: manual
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteMany
  nfs:
    server: 10.244.1.4
    path: "/"
声明一个1G大小的PVC

apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: manual
  resources:
    requests:
      storage: 1Gi
用户创建的 PVC 要真正被容器使用起来，就必须先和某个符合条件的 PV 进行绑定。这里要检查的条件，包括两部分：
1、第一个条件，当然是 PV 和 PVC 的 spec 字段。比如，PV 的存储（storage）大小，就必须满足 PVC 的要求。
2、而第二个条件，则是 PV 和 PVC 的 storageClassName 字段必须一样。这个机制我会在本篇文章的最后一部分专门介绍
PV和PVC是一对一的关系
pod使用上面声明的PVC

apiVersion: v1
kind: Pod
metadata:
  labels:
    role: web-frontend
spec:
  containers:
  - name: web
    image: nginx
    ports:
      - name: web
        containerPort: 80
    volumeMounts:
        - name: nfs
          mountPath: "/usr/share/nginx/html"
  volumes:
  - name: nfs
    persistentVolumeClaim:
      claimName: nfs

pod在volume字段声明自己要使用的PVC的名字，等这个pod被创建之后，kubelet就会把这个PVC所对应的PV，挂载在这个pod容器内的目录上
PVC和PV的设计，其实和“面向对象”的思想完全一致，PVC是持久化存储的“接口”，提供对某种持久化存储的描述；持久化存储的具体实现部分由PV负责完成

在k8s中，存在一个专门处理持久化存储的控制器，Volume Controller
这个Contruller维护着多个控制循环，其中一个循环，扮演的就是撮合 PV 和 PVC 的“红娘”的角色。它的名字叫作 PersistentVolumeController
PersistentVolumeController 会不断地查看当前每一个 PVC，是不是已经处于 Bound（已绑定）状态。如果不是，那它就会遍历所有的、可用的 PV，并尝试将其与这个“单身”的 PVC 进行绑定。这样，Kubernetes 就可以保证用户提交的每一个 PVC，只要有合适的 PV 出现，它就能够很快进入绑定状态
所谓将一个 PV 与 PVC 进行“绑定”，其实就是将这个 PV 对象的名字，填在了 PVC 对象的 spec.volumeName 字段上。
[root@k8smaster volume]# kubectl get pv
NAME      CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                STORAGECLASS   REASON   AGE
pv-test   3Gi        RWX            Recycle          Bound    default/www-web1-0   manual                  2d1
9h[root@k8smaster volume]# kubectl get pv
NAME      CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                STORAGECLASS   REASON   AGE
pv-test   3Gi        RWX            Recycle          Bound    default/www-web1-0   manual                  2d19h
[root@k8smaster volume]# kubectl get pvc
NAME         STATUS    VOLUME    CAPACITY   ACCESS MODES   STORAGECLASS   AGE
www-web1-0   Bound     pv-test   1Gi        RWX            manual         2d19h
在这个www-web1-0的pvc的yaml配置信息里可以看到volumename的字段就是pv-test
[root@k8smaster volume]# kubectl get pvc www-web1-0  -o yaml | grep -i  volumename
  volumeName: pv-test
持久化 Volume 的实现，往往依赖于一个远程存储服务，比如：远程文件存储（比如，NFS、GlusterFS）、远程块存储（比如，公有云提供的远程磁盘）等等。
而 Kubernetes 需要做的工作，就是使用这些存储服务，来为容器准备一个持久化的宿主机目录，以供将来进行绑定挂载时使用。而所谓“持久化”，指的是容器在这个目录里写入的文件，都会保存在远程存储中，从而使得这个目录具备了“持久性”。
这个准备“持久化”宿主机目录的过程，我们可以形象地称为“两阶段处理”
1、“第一阶段”（Attach），
由volume controller负责维护，控制循环的名字叫AttachDetachController。作用就是不断的检查每一个pod对应的pv，和这个pod所在宿主机之间的挂载情况，从而决定是否对这个PV进行Attach或Dettach操作。
volume控制器是kube-controller-manager的一部分，所以一定是运行在master节点上。
Kubernetes 提供的可用参数是 nodeName，即宿主机的名字
2、而对于“第二阶段”（Mount）；
第二阶段的mount或者umount操作，必须发生在pod对应的宿主机上，所以它是kubelet组件的一部分，这个控制循环的名字，叫VolumeManagerReconciler，运行起来后是一个独立于kubelet的goroutime
Kubernetes 提供的可用参数是 dir，即 Volume 的宿主机目录

删除一个PV时，k8s需要umount和dettch两个阶段，执行反向操作即可

二、StorageClass
在大规模集群中PVC数量很多，如果手动创建PV，耗时耗力
k8s提供了一套可以自动创建PV的机制，即：Dynamic Provisioning。
人工管理 PV 的方式就叫作 Static Provisioning
Dynamic Provisioning 机制工作的核心，在于一个名叫 StorageClass 的 API 对象。
StorageClass 对象的作用，其实就是创建 PV 的模板。
第一，PV 的属性。比如，存储类型、Volume 的大小等等。
第二，创建这种 PV 需要用到的存储插件。比如，Ceph 等等。
有了着两个就可以根据用户提交的PVC找到对应的StorageClass，然后调用storageclass声明的存储插件，创建PV
使用ROOK存储服务的话，storageclass使用如下的yml文件定义

apiVersion: ceph.rook.io/v1beta1
kind: Pool
metadata:
  name: replicapool
  namespace: rook-ceph
spec:
  replicated:
    size: 3
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: block-service
provisioner: ceph.rook.io/block
parameters:
  pool: replicapool
  #The value of "clusterNamespace" MUST be the same as the one in which your rook cluster exist
  clusterNamespace: rook-ceph
定义了一个block-service的storageclass，存储插件是rook-ceph
创建这个storageclas
$ kubectl create -f sc.yaml
开发者，只需要在PVC指定storageclass名字即可

apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: claim1
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: block-service
  resources:
    requests:
      storage: 30Gi
在这个PVC里添加了一个storageclassname的字段:block-service
当我们通过 kubectl create 创建上述 PVC 对象之后，Kubernetes 就会调用ROOK-ceph 的 API，创建出一块Persistent Disk。然后，再使用这个 Persistent Disk 的信息，自动创建出一个对应的 PV 对象。

有了 Dynamic Provisioning 机制，运维人员只需要在 Kubernetes 集群里创建出数量有限的 StorageClass 对象就可以了。这就好比，运维人员在 Kubernetes 集群里创建出了各种各样的 PV 模板。这时候，当开发人员提交了包含 StorageClass 字段的 PVC 之后，Kubernetes 就会根据这个 StorageClass 创建出对应的 PV。
StorageClass 并不是专门为了 Dynamic Provisioning 而设计的。
storageClassName并不只是为了动态创建PV而使用的，当资源中并不存在该storageclass对象，但是PV和PVC又指定了相同的storageClassName名字意味者我想将两者进行绑定，为我所控
如果你的集群已经开启了名叫 DefaultStorageClass 的 Admission Plugin，它就会为 PVC 和 PV 自动添加一个默认的 StorageClass；否则，PVC 的 storageClassName 的值就是“”，这也意味着它只能够跟 storageClassName 也是“”的 PV 进行绑定。

PVC 描述的，是 Pod 想要使用的持久化存储的属性，比如存储的大小、读写权限等。
PV 描述的，则是一个具体的 Volume 的属性，比如 Volume 的类型、挂载目录、远程存储服务器地址等。
而 StorageClass 的作用，则是充当 PV 的模板。并且，只有同属于一个 StorageClass 的 PV 和 PVC，才可以绑定在一起。
StorageClass 的另一个重要作用，是指定 PV 的 Provisioner（存储插件）。这时候，如果你的存储插件支持 Dynamic Provisioning 的话，Kubernetes 就可以自动为你创建 PV 了。

总结：
用户提交请求创建pod，Kubernetes发现这个pod声明使用了PVC，那就靠PersistentVolumeController帮它找一个PV配对。
没有现成的PV，就去找对应的StorageClass，帮它新创建一个PV，然后和PVC完成绑定。
新创建的PV，还只是一个API 对象，需要经过“两阶段处理”变成宿主机上的“持久化 Volume”才真正有用：
第一阶段由运行在master上的AttachDetachController负责，为这个PV完成 Attach 操作，为宿主机挂载远程磁盘；
第二阶段是运行在每个节点上kubelet组件的内部，把第一步attach的远程磁盘 mount 到宿主机目录。这个控制循环叫VolumeManagerReconciler，运行在独立的Goroutine，不会阻塞kubelet主循环。
完成这两步，PV对应的“持久化 Volume”就准备好了，POD可以正常启动，将“持久化 Volume”挂载在容器内指定的路径。


