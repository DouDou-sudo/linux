save保存的是镜像，export保存的是容器
load用来载入镜像包，import用来载入容器包，但两者都会恢复为镜像
load不能对载入的镜像重命名，而import可以为镜像指定新名称
load不能载入容器包，import能载入镜像包，但不能使用
export导入的镜像看不到镜像的构建history，

保存镜像

    $ docker save NAME:TAG  -o filename
    $ file filename
    filename: POSIX tar archive
这里的 filename 可以为任意名称甚至任意后缀名，但文件的本质都是归档文件
注意：如果同名则会覆盖（没有警告）
若使用 gzip 压缩：

    $ docker save alpine | gzip > alpine-latest.tar.gz
将 alpine-latest.tar.gz 文件复制到了到了另一个机器上，可以用下面这个命令加载镜像：

    $ docker load -i alpine-latest.tar.gz
    Loaded image: alpine:latest
如果我们结合这两个命令以及 ssh 甚至 pv 的话，利用 Linux 强大的管道，我们可以写一个命令完成从一个机器将镜像迁移到另一个机器，并且带进度条的功能：

    docker save <镜像名> | bzip2 | pv | ssh <用户名>@<主机名> 'cat | docker load'

导出和导入容器

运行一个容器
    [root@k8snode2 dockerfile]# docker run -d centos:hello /bin/sh -c "while true;do echo hello world;sleep 2;done"
    60cc0a4b838e3b10f16990e5eff91feb311b9bbfe4c600b2a8518b25ee280ddf
    [root@k8snode2 dockerfile]# docker ps
    CONTAINER ID        IMAGE                                               COMMAND                  CREATED             STATUS              PORTS               NAMES
    60cc0a4b838e        centos:hello                                        "/bin/sh -c 'while t…"   2 seconds ago       Up 2 seconds                            condescending_hy
    patia
导出容器
    [root@k8snode2 dockerfile]# docker export 60cc0a4b838e >hello.tar
    [root@k8snode2 dockerfile]# ls
    hello.tar
    或者使用这种方式
    $ docker export -o mysql-`date +%Y%m%d`.tar a404c6c174a2
将tar包传到其他服务器上
    [root@k8snode2 dockerfile]# scp -rp hello.tar 192.168.189.201:/root/
导入容器并指定image名称
    [root@k8snode1 ~]# docker import hello.tar centos:test1
    sha256:46c42a0d51329c77211819b6ba489a563edb9c73a911c3806ea607ac376a16b2
    [root@k8snode1 ~]# docker images centos
    REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
    centos              test1               46c42a0d5132        7 seconds ago       233MB
    或者使用这种方式导入，并指定image名称
    [root@k8snode1 ~]# cat hello.tar | docker import - centos:test
    sha256:01b4d32cf4351dc08a5e78215397eb8bbee3ad0655793402dd96ed784370fefd
    [root@k8snode1 ~]# docker images centos
    REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
    centos              test                01b4d32cf435        5 seconds ago       233MB
