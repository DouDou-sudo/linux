docker提供了三种持久化数据的方式：
1、bind mount: 把宿主机的某个文件或目录挂载到指定目录或文件下，非docker进程可以修改这些数据
2、volumes: 实际上，是在 Docker 的/var/lib/docker/volumes/文件夹内创建一个相同名字的文件夹来保存数据。因为这个文件夹在 Docker 管控范围里，Docker 可以根据挂载的设定来控制容器对 Volume 的读写权限。
3、tmpfs mount: 存于内存中，并不是持久化到磁盘。在容器的生命周期中，它能被容器用来存放非持久化的状态或敏感信息


一、volumes
如果没有显式创建，一个卷会在最开始挂载时被创建。当容器停止时，卷仍然存在。多个容器可以通过read-write或read-only的方式使用同一个卷。
注意:

    只有在显式删除时，卷才会被删除。​如果将一个空卷挂载到容器中一个存有文件或目录的目录中，这些文件或目录会被拷贝到空卷中；如果将一个非空卷挂载到容器中一个存有文件或目录的目录中，这些文件或目录会被隐藏。​

使用

    创建：​​docker volume create​​ 卷名(创建后可以在/var/lib/docker/volumes/目录下查找到对应目录)

    删除某个卷：​​docker volume rm 卷名​​

    删除所有未使用的卷：​​docker volume prune​​

    列出所有卷：​​docker volume ls​​

    查看某个卷的信息：​​docker volume inspect 卷名​​
场景

    多个运行容器间共享数据

    当Docker主机不确保具有给定的目录或文件时。卷可以将容器运行时与Docker主机的配置解耦合

    备份、恢复、或将数据从一个Docker主机迁移到另一个Docker主机时
eg 

    创建volume
    [root@k8snode2 volumes]# docker volume create nginx-volume
    查看/var/lib/docker/volumes目录
    [root@k8snode2 volumes]# ls /var/lib/docker/volumes/
    metadata.db   nginx-volume    backingFsBlockDev
    通过inspect查看volume详情
    [root@k8snode2 ~]# docker volume inspect nginx-volume 
    '[
        {
            "CreatedAt": "2022-10-27T23:15:29+08:00",
            "Driver": "local",
            "Labels": {},
            "Mountpoint": "/var/lib/docker/volumes/nginx-volume/_data",
            "Name": "nginx-volume",
            "Options": {},
            "Scope": "local"
        }
    ]

    创建容器并挂载volume
    [root@k8snode2 volumes]# docker run -d --mount type=volume,src=nginx-volume,target=/usr/share/nginx/html -p 82:80  --name nginx-volume nginx:latest 
    05f737bee1fc4155568d494f93b2f683b7890473ccf0dd5a0c6cac84d2c185a5
    [root@k8snode2 volumes]# docker ps | grep nginx
    05f737bee1fc        nginx:latest                                        "/docker-entrypoint.…"   9 seconds ago       Up 8 seconds        0.0.0.0:82->80/tcp   nginx-volume
    通过inspect查看挂载
    [root@k8snode2 ~]# docker inspect nginx-volume
            "Mounts": [
                {
                    "Type": "volume",
                    "Name": "nginx-volume",
                    "Source": "/var/lib/docker/volumes/nginx-volume/_data",
                    "Destination": "/usr/share/nginx/html",
                    "Driver": "local",
                    "Mode": "z",
                    "RW": true,
                    "Propagation": ""
                }
            ],
    进入刚才创建的容器
    [root@k8snode2 ~]# docker exec -it nginx-volume bash
    通过df -h是查看不到volume挂载的
    root@05f737bee1fc:/# df -h   
    Filesystem               Size  Used Avail Use% Mounted on
    overlay                   17G  5.6G   12G  33% /
    tmpfs                     64M     0   64M   0% /dev
    tmpfs                    910M     0  910M   0% /sys/fs/cgroup
    shm                       64M     0   64M   0% /dev/shm
    /dev/mapper/centos-root   17G  5.6G   12G  33% /etc/hosts
    tmpfs                    910M     0  910M   0% /proc/asound
    tmpfs                    910M     0  910M   0% /proc/acpi
    tmpfs                    910M     0  910M   0% /proc/scsi
    tmpfs                    910M     0  910M   0% /sys/firmware
    查看挂载目录下的文件，挂载的volume是个空卷，但是此时下面存在文件，印证了上面的注意，当空卷挂载到容器中时一个有文件或者目录的目录中，这些文件或者目录会被拷贝到空卷中
    root@05f737bee1fc:/# cd /usr/share/nginx/html/
    root@05f737bee1fc:/usr/share/nginx/html# ls
    50x.html  index.html
    在卷中查看拷贝过来的文件
    [root@k8snode2 _data]# pwd
    /var/lib/docker/volumes/nginx-volume/_data
    [root@k8snode2 _data]# ls
    50x.html  index.html
    删除容器
    [root@k8snode2 _data]# docker stop nginx-volume 
    nginx-volume
    [root@k8snode2 _data]# docker rm nginx-volume
    nginx-volume
    volume不会被删除
    [root@k8snode2 _data]# pwd
    /var/lib/docker/volumes/nginx-volume/_data
    [root@k8snode2 _data]# ls
    50x.html  index.html
    将volume更改为只读
    [root@k8snode2 ~]# docker run -d --mount type=volume,src=nginx-volume,target=/usr/share/nginx/html,ro -p 82:80  --name nginx-volume nginx:latest 
    [root@k8snode2 _data]# docker inspect nginx-volume 
        "Mounts": [
            {
                "Type": "volume",
                "Name": "nginx-volume",
                "Source": "/var/lib/docker/volumes/nginx-volume/_data",
                "Destination": "/usr/share/nginx/html",
                "Driver": "local",
                "Mode": "z",
                "RW": false,                ##此处为false，为只读模式
                "Propagation": ""
            }
        ],
    进入容器尝试在挂载目录下创建文件
    root@99d30e601b28:/usr/share/nginx/html# touch 1
    touch: cannot touch '1': Read-only file system

二、bind mount
主机中的文件或目录通过全路径被引用。在使用绑定挂载时，这些目录或文件不一定要已经存在。

注意：

    ​如果使用这种方式将一个目录挂载到容器中一个存有文件或目录的目录中，这些文件或目录会被隐藏；如果主机中的文件或目录不存在，当使用​​--mount​​挂载时，Docker会报错，当使用​​-v​​或​​--volume​​时，会在主机上创建目录​
场景

    主机与容器共享配置文件（Docker默认情况下通过这种方式为容器提供DNS解析，通过将/etc/resolv.conf挂载到容器中）

    共享源代码或build artifacts（比如将Maven的target/目录挂载到容器中，每次在Docker主机中build Maven工程时，容器能够访问到那些rebuilt artifacts）

    当 docker主机中的文件或目录结构和容器需要的一致时

eg

    创建容器，挂载一个文件，源文件为绝对路径，必须指定挂载到容器里的文件名，
    [root@k8snode2 opt]# docker run -d --name nginx-bind -v /opt/abc:/usr/share/nginx/html/a:ro nginx:latest 
    挂载多个文件
    [root@k8snode2 ms]# docker run -it --name test --rm -v /root/ms/mariadb-5.5.68-1.el7.x86_64.rpm:/root/m1.rpm -v /root/ms/mariadb-libs-5.5.68-1.el7.x86_64.rpm:/root/2.rpm centos:latest
    创建容器，使用bind挂载，主机中的/opt/nginx-bind目录不存在它会自动创建

    [root@k8snode2 ~]# docker run -d --name nginx-bind -v /opt/nginx-bind:/usr/share/nginx/html -p 83:80 nginx:latest 

    进入容器查看挂载目录下的文件，发现原有文件被覆盖

    [root@k8snode2 nginx-bind]# docker exec -it nginx-bind  bash
    root@32e094973506:/# cd /usr/share/nginx/html/
    root@32e094973506:/usr/share/nginx/html# ls
    root@32e094973506:/usr/share/nginx/html# 

    在宿主机目录下创建文件，进入容器可以正常看到
    [root@k8snode2 nginx-bind]# pwd
    /opt/nginx-bind
    [root@k8snode2 nginx-bind]# touch 1
    root@32e094973506:/usr/share/nginx/html# ls
    1
    使用--mount挂载bind，宿主机目录存在，可以正常创建
    [root@k8snode2 _data]# docker run -d --name nginx-bind --mount type=bind,src=/opt/nginx-bind,target=/usr/share/nginx/html -p 83:80 nginx:latest
    使用--mount挂载bind，如何所指定宿主机目录不存在，则会报错
    [root@k8snode2 _data]# docker run -d --name nginx-bind --mount type=bind,src=/opt/nginx-test,target=/usr/share/nginx/html -p 83:80 nginx:latest 
    docker: Error response from daemon: invalid mount config for type "bind": bind source path does not exist: /opt/nginx-test.
    See 'docker run --help'.

    以只读模式映射bind mount，也可以不使用-v，使用--volume(docker run -d --name nginx-bind --mount type=bind,src=/opt/nginx-bind,target=/usr/share/nginx/html,ro -p 83:80 nginx:latest)
    [root@k8snode2 _data]# docker run -d --name nginx-bind -v /opt/nginx-bind/:/usr/share/nginx/html:ro -p 83:80 nginx:latest 
    root@1e5a4560c37f:/usr/share/nginx/html# touch 2
    touch: cannot touch '2': Read-only file system
    使用inspect查看容器
    [root@k8snode2 _data]# docker inspect nginx-bind 
            "Mounts": [
            {
                "Type": "bind",
                "Source": "/opt/nginx-bind",
                "Destination": "/usr/share/nginx/html",
                "Mode": "ro",
                "RW": false,
                "Propagation": "rprivate"
            }
        ],

三、tmpfs
相对于volumes和bind mount，tmpfs mount是临时的，只在主机内存中持久化。当容器停止，tmpfs mount会被移除。对于临时存放敏感文件很有用

不同于volumes和bind mount，多个容器无法共享tmpfs mount

场景

场景

    最好的使用场景是你既不想将数据存于主机，又不想存于容器中时。这可以是出于安全的考虑，或当应用需要写大量非持久性的状态数据时为了保护容器的性能

eg

    创建容器，使用tmpfs挂载
    [root@k8snode2 _data]# docker run --mount type=tmpfs,target=/usr/share/nginx/html -p 84:80 -d --name nginx-tmpfs nginx:latest
    使用inspect查看容器
    [root@k8snode2 _data]# docker inspect nginx-tmpfs 
            "Mounts": [
                {
                    "Type": "tmpfs",
                    "Source": "",
                    "Destination": "/usr/share/nginx/html",
                    "Mode": "",
                    "RW": true,
                    "Propagation": ""
                }
            ],

