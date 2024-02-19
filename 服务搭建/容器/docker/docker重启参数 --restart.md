--restart=always参数能够使我们在重启docker时，自动启动相关容器。
Docker容器的重启策略如下：
no，默认策略，在容器退出时不重启容器
on-failure，在容器非正常退出时（退出状态非0），才会重启容器
on-failure:3，在容器非正常退出时重启容器，最多重启3次
always，在容器退出时总是重启容器
unless-stopped，在容器退出时总是重启容器，但是不考虑在Docker守护进程启动时就已经停止了的容器

一、启动时添加参数

1、使用-restart=always参数启动容器

    [root@k8snode2 dockerfile]# docker run -itd --restart=always --name always centos:latest 
    509b38839b9881e0223aeb3c628f42332387ba1842aae5b3cc9bcc9b196fcd62
2、不加–restart=always参数启动容器

    [root@k8snode2 dockerfile]# docker run -itd  --name default centos:latest 
    08323389b4c27ebed24facde3c6db2408019838b0668e8a9e51c81e60d56949a
查看container

    [root@k8snode2 dockerfile]# docker ps
    CONTAINER ID        IMAGE                                               COMMAND                  CREATED             STATUS              PORTS               NAMES
    08323389b4c2        centos:latest                                       "/bin/bash"              3 seconds ago       Up 2 seconds                            default
    509b38839b98        centos:latest                                       "/bin/bash"              28 seconds ago      Up 27 seconds                           always
重启docker，查看container

    [root@k8snode2 dockerfile]# docker ps
    509b38839b98        centos:latest                                       "/bin/bash"              53 seconds ago      Up 4 seconds                            always

二、启动后修改
 在启动时如果没有添加这个参数怎么办呢，比如1a7a3b5112fd这个容器在启动的时候是没有添加–restart=always参数的，针对这种情况我们可以使用命令进行修改。docker container update --restart=always 容器名字

1、运行一个容器

    [root@k8snode2 dockerfile]# docker run -itd centos:latest 
    a84a249a58a739c69ab1505dd86d6ecd5e3a77c50ae51318b991c5f796fe2a2b
    [root@k8snode2 dockerfile]# docker ps
    CONTAINER ID        IMAGE                                               COMMAND                  CREATED              STATUS              PORTS               NAMES
    a84a249a58a7        centos:latest                                       "/bin/bash"              About a minute ago   Up About a minute                       dazzling_sutherland
2、启动后使用命令添加

    [root@k8snode2 dockerfile]# docker container update --restart=always a84a249a58a7
    a84a249a58a7
重启docker，查看container

    [root@k8snode2 dockerfile]# docker ps
    a84a249a58a7        centos:latest                                       "/bin/bash"              4 minutes ago       Up 4 seconds                            dazzling_sutherland