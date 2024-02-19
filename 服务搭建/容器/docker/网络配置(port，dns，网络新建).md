查看映射端口配置

    [root@k8snode2 ~]# docker port nginx 
    80/tcp -> 0.0.0.0:80
    [root@k8snode2 ~]# docker port nginx 80
    0.0.0.0:80
端口映射可以指定-p port:port,-p IP::port,-p IP:host_port:container_port来指定运行访问容器的主机上的ip，端口
   
    docker run -itd -p 192.168.189.200:8081:80 centos:7
如果希望永久绑定到某个固定的ip地址，可以在/etc/docker/daemon.json中添加以下内容

    {
    "ip": "1.2.3.4"
    }
0.0.0.0意味着接受主机所有接口的流量，指定某一个ip，例如1.2.3.4,只接受ip地址为1.2.3.4接口上的流量

新建网络

    [root@k8snode2 ~]# docker network create -d bridge mynet
    d0cc9e06f473ab370d4b9ccf4b8b3e0595de3ee6c26ef2ecedb161ac94a3779e
查看网络，默认有三个，bridge,host,none

    [root@k8snode2 ~]# docker network ls
    NETWORK ID          NAME                DRIVER              SCOPE
    5288124ebe27        bridge              bridge              local
    6659c8aa768d        host                host                local
    d0cc9e06f473        mynet               bridge              local
    ae1adbdfff1a        none                null                local
运行两个容器，并连接到新建的网络

    [root@k8snode2 ~]# docker run -it --rm --name busybox1 --network mynet busybox:latest sh
    [root@k8snode2 ~]# docker run -it --rm --name busybox2 --network mynet busybox:latest sh
进入容器，ping另外一个连接到mynet网络的容器

    容器1
    / # ip a
    1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue qlen 1000
        link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
        inet 127.0.0.1/8 scope host lo
        valid_lft forever preferred_lft forever
    23: eth0@if24: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue 
        link/ether 02:42:ac:12:00:02 brd ff:ff:ff:ff:ff:ff
        inet 172.18.0.2/16 brd 172.18.255.255 scope global eth0
        valid_lft forever preferred_lft forever
    / # ping 172.18.0.3
    PING 172.18.0.3 (172.18.0.3): 56 data bytes
    64 bytes from 172.18.0.3: seq=0 ttl=64 time=0.064 ms
    容器2
    / # ip a
    1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue qlen 1000
        link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
        inet 127.0.0.1/8 scope host lo
        valid_lft forever preferred_lft forever
    25: eth0@if26: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue 
        link/ether 02:42:ac:12:00:03 brd ff:ff:ff:ff:ff:ff
        inet 172.18.0.3/16 brd 172.18.255.255 scope global eth0
        valid_lft forever preferred_lft forever
    / # ping 172.18.0.2
    PING 172.18.0.2 (172.18.0.2): 56 data bytes
    64 bytes from 172.18.0.2: seq=0 ttl=64 time=0.078 ms
    ping默认bridge网络下的容器，ping不通
    / # ping 172.17.0.3
    PING 172.17.0.3 (172.17.0.3): 56 data bytes
高级网络配置
当docker启动时，会自动在主机上创建一个docker0的虚拟网桥，就是一个linux的bridge，同时docker随机分配一个本地未占用的私有网段给docker接口，比如:172.17.42.1/16。
此后启动容器就会自动分配一个这个网段的地址，当创建容器时会自动创建一个veth pair 接口（当数据包发送到一个接口时，另外一个接口也可以收到相同的数据包）。这个接口一段在容器给内，另一端被挂载在docker0网桥上。

如何查看veth pair 接口的对应关系呢？
    
    创建一个容器
    [root@7e100666cec3 /]# ip a
    1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
        link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
        inet 127.0.0.1/8 scope host lo
        valid_lft forever preferred_lft forever
    44: eth0@if45: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
        link/ether 02:42:ac:11:00:05 brd ff:ff:ff:ff:ff:ff link-netnsid 0
        inet 172.17.0.5/16 brd 172.17.255.255 scope global eth0
        valid_lft forever preferred_lft forever
    在host上执行ip a查看
    45: veth614a001@if44: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue master docker0 state UP group default 
        link/ether 3e:d8:c9:cf:b1:9a brd ff:ff:ff:ff:ff:ff link-netnsid 3
        inet6 fe80::3cd8:c9ff:fecf:b19a/64 scope link 
        valid_lft forever preferred_lft forever
一对veth pair 接口的对应关系一般为docker容器eth0网口ifxx对应host的ifxx+1

在 docker run 的时候通过 --net 参数来指定容器的网络配置，有4个可选值：

    --net=bridge 这个是默认值，连接到默认的网桥。
    --net=host 告诉 Docker 不要将容器网络放到隔离的命名空间中，即不要容器化容器内的网络。此时容器使用本地主机的网络，它拥有完全的本地主机接口访问权限。容器进程可以跟主机其它 root 进程一样可以打开低范围的端口，可以访问本地网络服务比如 D-bus，还可以让容器做一些影响整个主机系统的事情，比如重启主机。因此使用这个选项的时候要非常小心。如果进一步的使用 --privileged=true，容器会被允许直接配置主机的网络堆栈。
    --net=container:NAME_or_ID 让 Docker 将新建容器的进程放到一个已存在容器的网络栈中，新容器进程有自己的文件系统、进程列表和资源限制，但会和已存在的容器共享 IP 地址和端口等网络资源，两者进程可以直接通过 lo 环回接口通信。
    --net=none 让 Docker 将新容器放到隔离的网络栈中，但是不进行网络配置。之后，用户可以自己进行配置。

--net=container:NAME_or_ID网络配置:

    [root@k8snode2 ~]# docker ps 
    CONTAINER ID        IMAGE                                               COMMAND                  CREATED             STATUS              PORTS               NAMES
    7e100666cec3        centos:latest                                       "/bin/bash"              4 hours ago         Up 4 hours                              test
    新建一个容器，网络加入test容器的网络配置
    [root@k8snode2 ~]# docker run -itd --rm --net=container:test --name ces centos:latest
    [root@k8snode2 ~]# docker exec -it ces bash
    进入容器查看test和ces容器的网络配置一致
    [root@7e100666cec3 /]# ip a
    1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
        link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
        inet 127.0.0.1/8 scope host lo
        valid_lft forever preferred_lft forever
    44: eth0@if45: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
        link/ether 02:42:ac:11:00:05 brd ff:ff:ff:ff:ff:ff link-netnsid 0
        inet 172.17.0.5/16 brd 172.17.255.255 scope global eth0
        valid_lft forever preferred_lft forever
    [root@k8snode2 ~]# docker exec -it test  bash
    [root@7e100666cec3 /]# ip a
    1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
        link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
        inet 127.0.0.1/8 scope host lo
        valid_lft forever preferred_lft forever
    44: eth0@if45: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
        link/ether 02:42:ac:11:00:05 brd ff:ff:ff:ff:ff:ff link-netnsid 0
        inet 172.17.0.5/16 brd 172.17.255.255 scope global eth0
        valid_lft forever preferred_lft forever

如何自定义配置容器的主机名和 DNS ?秘诀就是 Docker 利用虚拟文件来挂载容器的 3 个相关配置文件。
在容器中使用 mount 命令可以看到挂载信息：

    / # mount
    /dev/mapper/centos-root on /etc/resolv.conf type xfs (rw,relatime,attr2,inode64,noquota)
    /dev/mapper/centos-root on /etc/hostname type xfs (rw,relatime,attr2,inode64,noquota)
    /dev/mapper/centos-root on /etc/hosts type xfs (rw,relatime,attr2,inode64,noquota)
当宿主机的dns更新时，所有docker容器的dns通过/etc/resolv.conf文件立刻得到更新
配置全部容器的dns，也可以在/etc/docker/daemin.json文件中增加以下内容来设置

    [root@k8snode2 ~]# cat /etc/docker/daemon.json 
    {
    "registry-mirrors": ["https://u96790vf.mirror.aliyuncs.com"],
    "exec-opts": ["native.cgroupdriver=systemd"],                   #注意每一个kv对后面的,号
    "dns": [
        "223.5.5.5"
    ]
    }
    此操作需要重启docker服务才可以生效
    重新创建一个容器，进入容器查看/etc/resolv.conf文件
    [root@k8snode2 ~]# docker run -it --rm --name test-dns centos:latest
    [root@aecfb557545c /]# cat /etc/resolv.conf 
    search localdomain
    nameserver 223.5.5.5
如果需要手动指定某一个容器给的配置，可以在docker run时添加参数：
-h HOSTNAME或者--hostname=HOSTNAME设定容器的主机名，会写入容器内的/etc/hostname和/etc/hosts文件中，但是在容器外看不到，docker ps看不到，也不会在其他容器的/etc/hosts中看到
--dns=IP_ADDRESS添加dns地址到容器的/etc/resolv.conf中
没有指定参数则使用docker的默认配置

    [root@k8snode2 ~]# docker run -it --rm -h dns --name test-dns --dns=1.1.2.3 centos:latest 
    [root@dns /]# hostname
    dns
    [root@dns /]# cat /etc/hosts
    127.0.0.1	localhost
    ::1	localhost ip6-localhost ip6-loopback
    fe00::0	ip6-localnet
    ff00::0	ip6-mcastprefix
    ff02::1	ip6-allnodes
    ff02::2	ip6-allrouters
    172.17.0.4	dns
    [root@dns /]# cat /etc/hostname 
    dns
    [root@dns /]# cat /etc/resolv.conf  
    search localdomain
    nameserver 1.1.2.3
在容器外docker ps 看不到dns的主机名

    [root@k8snode2 ~]# docker ps | grep "dns"
    461467cda29c        centos:latest                                       "/bin/bash"              37 seconds ago      Up 37 seconds                           test-dns


