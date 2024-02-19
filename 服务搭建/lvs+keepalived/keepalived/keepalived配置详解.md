```
[root@master ~]# cat /etc/keepalived/keepalived.conf
! Configuration File for keepalived
global_defs {
   notification_email {         #设置报警邮箱地址，需要开启sendmail服务，可以设置多个
     acassen@firewall.loc
     failover@firewall.loc
     sysadmin@firewall.loc
   }
   notification_email_from Alexandre.Cassen@firewall.loc
   router_id master             #router_id自定义，每个节点不一样
#   router_id node1
#   router_id node2
#   smtp_server 192.168.200.1   #设置邮件的smtp server地址
#   smtp_connect_timeout 30     #设置连接smtp server的超时时间
#   vrrp_skip_check_adv_addr
#   vrrp_strict                 
#如果配置了vrrp_strict参数或npaccept参数，每次发生主备切换都会在主节点上上添加如下防火墙规则，拒绝所有连接，导致不能访问后端服务。使用yum安装keepalived配置文件中默认有该选项，注释掉该选项
#DROP       all  --  anywhere             anywhere             match-set keepalived dst
#   vrrp_garp_interval 0
#   vrrp_gna_interval 0
}

vrrp_script check_haproxy {                   #健康检测配置名
    script "/etc/keepalived/check_haproxy.sh" #指定检测脚本路径，脚本必须有权限执行
    interval 1                                #指定检测间隔时间,1s
    weight -30                                #当检测失败时，修改优先级priority+weight
    fall 3
    rise 2
    timeout 2
}


vrrp_instance VI_1 {      #定义一个虚拟路由
    state MASTER          #当前节点在虚拟路由器上的角色[MASTER|BACKUP]
#   state BACKUP
    interface ens33       #虚拟路由器使用的物理接口
    virtual_router_id 51  #虚拟routeer_id[0-255]，一个集群内的每个节点保持一致
    priority 100          #节点优先级[1-254]，每个节点不能一样，一样可能会导致脑裂
#   priority 99
#   priority 98
#   nopreempt             #不抢占模式，如果有节点的优先级高于当前master的优先级，不会抢占，默认为抢占
#   mcast_src_ip          #设置组播包发送的地址，如果不设置，使用interface所对应的ip
    advert_int 1          #组播发送时间间隔，包含主机优先级，心跳等，默认为1s
    authentication {      #认证配置
        auth_type PASS    #认证方式,[PASS|AH]
        auth_pass 1111    #密码修改为其他，不要使用默认密码
    }
    virtual_ipaddress {    #VIP地址，可以有多个，可以指定绑定的网卡和别名
        192.168.189.199 dev ens33 label ens33:1
    }
    track_interface {     #监控网口,如果网口出现问题节点会进入FAULT(故障)状态
        ens33
    }
    track_script {        #指定健康检测配置
        check_haproxy
    }
#    notify_master`  #当节点为Master时要执行的脚本，这个脚本可以是一个状态报警脚本，也可以是一个服务管理脚本。Keepalived允许脚本传入参数，因此灵活性很强。
#    notify_backup   #当节点为BACKUP时要执行的脚本，同理，这个脚本可以是一个状态报警脚本，也可以是一个服务管理脚本。
#    notify_fault    #当节点为FAULT时要执行的脚本，脚本功能与前两个类似。notify_stop：指定当Keepalived程序终止时需要执行的脚本。
#    notify_stop     #当节点为STOP时要执行的脚本
}
```
以下脚本需要根据实际需要更改端口名，脚本权限要注意更改
```
[root@node2 ~]# ll -a /etc/keepalived/check_haproxy.sh 
-rwxr-xr-x 1 root root 398 Jan 25 10:18 /etc/keepalived/check_haproxy.sh
[root@master ~]# cat /etc/keepalived/check_haproxy.sh
#!/bin/bash
counter=$(netstat -tnlp | grep -w 8080 | wc -l)
if [ "$counter" == "0" ];then
#  systemctl restart haproxy.service              #无状态服务，可以直接设置对应服务开机自启，不使用keepalived对服务进行管理
  counter=$(netstat -tnlp | grep -w 8080 | wc -l)
  if [ "$counter" != "0" ];then
    exit 0
  fi
  exit 1
fi
```
