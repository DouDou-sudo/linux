所谓定制镜像，那一定是以一个镜像为基础，在其上进行定制
基础镜像是必须指定的。而 FROM 就是指定 基础镜像，因此一个 Dockerfile 中 FROM 是必备的指令，并且必须是第一条指令

除了选择现有镜像为基础镜像外，Docker 还存在一个特殊的镜像，名为 scratch。这个镜像是虚拟的概念，并不实际存在，它表示一个空白的镜像。如果你以 scratch 为基础镜像的话，意味着你不以任何镜像为基础，接下来所写的指令将作为镜像第一层开始存在。
对于 Linux 下静态编译的程序来说，并不需要有操作系统提供运行时支持，所需的一切库都已经在可执行文件里了

Dockerfile 中每一个指令都会建立一层

    FROM debian:stretch

    RUN apt-get update
    RUN apt-get install -y gcc libc6-dev make wget
    RUN wget -O redis.tar.gz "http://download.redis.io/releases/redis-5.0.3.tar.gz"
    RUN mkdir -p /usr/src/redis
    RUN tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1
    RUN make -C /usr/src/redis
    RUN make -C /usr/src/redis install
如上这种，创建了7层镜像，构建的镜像就会太臃肿，层数太多，现在最大层数不能超过127层
对上面的dockerfile进行优化

    FROM debian:stretch

    RUN set -x; buildDeps='gcc libc6-dev make wget' \
        && apt-get update \
        && apt-get install -y $buildDeps \
        && wget -O redis.tar.gz "http://download.redis.io/releases/redis-5.0.3.tar.gz" \
        && mkdir -p /usr/src/redis \
        && tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
        && make -C /usr/src/redis \
        && make -C /usr/src/redis install \
        && rm -rf /var/lib/apt/lists/* \
        && rm redis.tar.gz \
        && rm -r /usr/src/redis \
        && apt-get purge -y --auto-remove $buildDeps
镜像构建上下文（Context）
首先我们要理解 docker build 的工作原理。Docker 在运行时分为 Docker 引擎（也就是服务端守护进程）和客户端工具。Docker 的引擎提供了一组 REST API，被称为 Docker Remote API (opens new window)，而如 docker 命令这样的客户端工具，则是通过这组 API 与 Docker 引擎交互，从而完成各种功能。因此，虽然表面上我们好像是在本机执行各种 docker 功能，但实际上，一切都是使用的远程调用形式在服务端（Docker 引擎）完成。也因为这种 C/S 设计，让我们操作远程服务器的 Docker 引擎变得轻而易举。

镜像构建时有时会使用COPY、ADD等指令，而 docker build 命令构建镜像，其实并非在本地构建，而是在服务端，也就是 Docker 引擎中构建的
当构建的时候，用户会指定构建镜像上下文的路径，docker build 命令得知这个路径后，会将路径下的所有内容打包，然后上传给 Docker 引擎。这样 Docker 引擎收到这个上下文包后，展开就会获得构建镜像所需的一切文件
COPY 这类指令中的源文件的路径都是相对路径。这也是初学者经常会问的为什么 COPY ../package.json /app 或者 COPY /opt/xxxx /app 无法工作的原因，因为这些路径已经超出了上下文的范围



Dockerfile 指令详解

COPY 复制文件
    格式：

        COPY [--chown=<user>:<group>] <源路径>... <目标路径>
        COPY [--chown=<user>:<group>] ["<源路径1>",... "<目标路径>"]
    COPY 指令将从构建上下文目录中 <源路径> 的文件/目录复制到新的一层的镜像内的 <目标路径> 位置。比如：

        COPY package.json /usr/src/app/
    <源路径> 可以是多个，甚至可以是通配符，其通配符规则要满足 Go 的 filepath.Match (opens new window)规则，如：

        COPY hom* /mydir/
        COPY hom?.txt /mydir/
    <目标路径> 可以是容器内的绝对路径，也可以是相对于工作目录的相对路径（工作目录可以用 WORKDIR 指令来指定）。目标路径不需要事先创建，如果目录不存在会在复制文件前先行创建缺失目录。

    此外，还需要注意一点，使用 COPY 指令，源文件的各种元数据都会保留。比如读、写、执行权限、文件变更时间等。这个特性对于镜像定制很有用。特别是构建相关文件都在使用 Git 进行管理的时候。
    在使用该指令的时候还可以加上 --chown=<user>:<group> 选项来改变文件的所属用户及所属组。

        COPY --chown=55:mygroup files* /mydir/
        COPY --chown=bin files* /mydir/
        COPY --chown=1 files* /mydir/
        COPY --chown=10:11 files* /mydir/
    如果源路径为文件夹，复制的时候不是直接复制该文件夹，而是将文件夹中的内容复制到目标路径。
当copy多个文件时，dockerfile方法如下：

    COPY count.sh  mariadb-5.5.68-1.el7.x86_64.rpm mariadb-libs-5.5.68-1.el7.x86_64.rpm /root
    报错:
    Step 2/6 : COPY count.sh  mariadb-5.5.68-1.el7.x86_64.rpm mariadb-libs-5.5.68-1.el7.x86_64.rpm /root
    When using COPY with more than one source file, the destination must be a directory and end with a /
    当copy多个文件时，目标必须是一个目录并且以/结尾
    如下就可以正常拷贝多个文件：
    COPY count.sh  mariadb-5.5.68-1.el7.x86_64.rpm mariadb-libs-5.5.68-1.el7.x86_64.rpm /root/
    进入容器查看
    [root@k8snode2 ms]# docker exec -it count bash
    [root@9c4eaaa334e1 ~]# cd /root/
    [root@9c4eaaa334e1 ~]# ls
    anaconda-ks.cfg  count.sh  mariadb-5.5.68-1.el7.x86_64.rpm  mariadb-libs-5.5.68-1.el7.x86_64.rpm

    eg：
    1、当dockerfile文件如下时，构建images会报错，如果想把一个文件拷贝到目录下，必须在目录后面加上/
    COPY app.py /code
    WORKDIR /code
    Cannot mkdir: /code is not a directory
    2、当dockerfike文件没有指定WORKDIR时, COPY 会拷贝到/目录下
    FROM centos:7
    COPY pip.conf  .

ADD指令 更高级的复制文件
    ADD 指令和 COPY 的格式和性质基本一致。但是在 COPY 基础上增加了一些功能。

    比如 <源路径> 可以是一个 URL，这种情况下，Docker 引擎会试图去下载这个链接的文件放到 <目标路径> 去。下载后的文件权限自动设置为 600，如果这并不是想要的权限，那么还需要增加额外的一层 RUN 进行权限调整，另外，如果下载的是个压缩包，需要解压缩，也一样还需要额外的一层 RUN 指令进行解压缩。所以不如直接使用 RUN 指令，然后使用 wget 或者 curl 工具下载，处理权限、解压缩、然后清理无用文件更合理。因此，这个功能其实并不实用，而且不推荐使用。一般在需要自动解压缩的时候使用ADD

    如果 <源路径> 为一个 tar 压缩文件的话，压缩格式为 gzip, bzip2 以及 xz 的情况下，ADD 指令将会自动解压缩这个压缩文件到 <目标路径> 去。

    在某些情况下，这个自动解压缩的功能非常有用，比如官方镜像 ubuntu 中：


    FROM scratch
    ADD ubuntu-xenial-core-cloudimg-amd64-root.tar.gz /
    ...
    但在某些情况下，如果我们真的是希望复制个压缩文件进去，而不解压缩，这时就不可以使用 ADD 命令了。
    在使用该指令的时候还可以加上 --chown=<user>:<group> 选项来改变文件的所属用户及所属组。

        ADD --chown=55:mygroup files* /mydir/
        ADD --chown=bin files* /mydir/
        ADD --chown=1 files* /mydir/
        ADD --chown=10:11 files* /mydir/

    eg：
    1、当dockerfile文件如下时，构建images会报错，如果想把一个文件拷贝到目录下，必须在目录后面加上/
    ADD app.py /code
    WORKDIR /code
    Cannot mkdir: /code is not a directory
    2、当dockerfike文件没有指定WORKDIR时, ADD 会拷贝到/目录下
    FROM centos:7
    ADD pip.conf  .

    
CMD 容器启动命令
    CMD 指令的格式和 RUN 相似，也是两种格式：

        shell 格式：CMD <命令>
        exec 格式：CMD ["可执行文件", "参数1", "参数2"...]
        参数列表格式：CMD ["参数1", "参数2"...]。在指定了 ENTRYPOINT 指令后，用 CMD 指定具体的参数。
    在指令格式上，一般推荐使用 exec 格式，这类格式在解析时会被解析为 JSON 数组，因此一定要使用双引号 "，而不要使用单引号。

    如果使用 shell 格式的话，实际的命令会被包装为 sh -c 的参数的形式进行执行。比如：

        CMD echo $HOME
    在实际执行中，会将其变更为：

        CMD [ "sh", "-c", "echo $HOME" ]
    将 CMD 写为：


    CMD service nginx start
    然后发现容器执行后就立即退出了。甚至在容器内去使用 systemctl 命令结果却发现根本执行不了
    对于容器而言，其启动程序就是容器应用进程，容器就是为了主进程而存在的，主进程退出，容器就失去了存在的意义，从而退出，其它辅助进程不是它需要关心的东西。

    而使用 service nginx start 命令，则是希望 upstart 来以后台守护进程形式启动 nginx 服务。而刚才说了 CMD service nginx start 会被理解为 CMD [ "sh", "-c", "service nginx start"]，因此主进程实际上是 sh。那么当 service nginx start 命令结束后，sh 也就结束了，sh 作为主进程退出了，自然就会令容器退出。

    正确的做法是直接执行 nginx 可执行文件，并且要求以前台形式运行。比如：

        CMD ["nginx", "-g", "daemon off;"]
ENTRYPOINT 入口点
    ENTRYPOINT 的格式和 RUN 指令格式一样，分为 exec 格式和 shell 格式。
    当指定了 ENTRYPOINT 后，CMD 的含义就发生了改变，不再是直接的运行其命令，而是将 CMD 的内容作为参数传给 ENTRYPOINT 指令，换句话说实际执行时，将变为：

        <ENTRYPOINT> "<CMD>"
    eg：
    ENTRYPOINT ["/root/count.sh"]
    当使用某个脚本作为启动命令时，这个脚本必须具有+x权限，否则会报错： 权限被拒绝
    docker: Error response from daemon: OCI runtime create failed: runc create failed: unable to start container process: exec: "/root/count.sh": permission denied: unknown.


ENV 设置环境变量
    格式有两种：

        ENV <key> <value>
        ENV <key1>=<value1> <key2>=<value2>...
    这个指令很简单，就是设置环境变量而已，无论是后面的其它指令，如 RUN，还是运行时的应用，都可以直接使用这里定义的环境变量。

        ENV VERSION=1.0 DEBUG=on \
            NAME="Happy Feet"
    定义了环境变量，那么在后续的指令中，就可以使用这个环境变量。比如在官方 node 镜像 Dockerfile 中，就有类似这样的代码：


        ENV NODE_VERSION 7.2.0

        RUN curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.xz" \
        && curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
        && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc \
        && grep " node-v$NODE_VERSION-linux-x64.tar.xz\$" SHASUMS256.txt | sha256sum -c - \
        && tar -xJf "node-v$NODE_VERSION-linux-x64.tar.xz" -C /usr/local --strip-components=1 \
        && rm "node-v$NODE_VERSION-linux-x64.tar.xz" SHASUMS256.txt.asc SHASUMS256.txt \
        && ln -s /usr/local/bin/node /usr/local/bin/nodejs
    在这里先定义了环境变量 NODE_VERSION，其后的 RUN 这层里，多次使用 $NODE_VERSION 来进行操作定制。可以看到，将来升级镜像构建版本的时候，只需要更新 7.2.0 即可，Dockerfile 构建维护变得更轻松了。

    下列指令可以支持环境变量展开： ADD、COPY、ENV、EXPOSE、FROM、LABEL、USER、WORKDIR、VOLUME、STOPSIGNAL、ONBUILD、RUN。

    可以从这个指令列表里感觉到，环境变量可以使用的地方很多，很强大。通过环境变量，我们可以让一份 Dockerfile 制作更多的镜像，只需使用不同的环境变量即可。

    eg：添加多个env，一般情况下添加环境变量为:ENV MYSQL_USER root不加"=",当一次添加多个env时，使用 \换行并在环境变量之间添加"="
    ENV MYSQL_USER="root" \
        MYSQL_IP=""   \
        MYSQL_PORT="3306" \
        MYSQL_PASSWD=""
    进入容器查看环境变量是否生效
    [root@k8snode2 ms]# docker exec -it count bash
    [root@9c4eaaa334e1 ~]# echo $MYSQL_USER
    root
    [root@9c4eaaa334e1 ~]# echo $MYSQL_PORT
    3306



ARG 构建参数
    构建参数和 ENV 的效果一样，都是设置环境变量。所不同的是，ARG 所设置的构建环境的环境变量，在将来容器运行时是不会存在这些环境变量的。但是不要因此就使用 ARG 保存密码之类的信息，因为 docker history 还是可以看到所有值的。

    Dockerfile 中的 ARG 指令是定义参数名称，以及定义其默认值。该默认值可以在构建命令 docker build 中用 --build-arg <参数名>=<值> 来覆盖。
    ARG 指令有生效范围，如果在 FROM 指令之前指定，那么只能用于 FROM 指令中。


        ARG DOCKER_USERNAME=library

        FROM ${DOCKER_USERNAME}/alpine

        RUN set -x ; echo ${DOCKER_USERNAME}
    使用上述 Dockerfile 会发现无法输出 ${DOCKER_USERNAME} 变量的值，要想正常输出，你必须在 FROM 之后再次指定 ARG


        # 只在 FROM 中生效
        ARG DOCKER_USERNAME=library

        FROM ${DOCKER_USERNAME}/alpine

        # 要想在 FROM 之后使用，必须再次指定
        ARG DOCKER_USERNAME=library

        RUN set -x ; echo ${DOCKER_USERNAME}

VOLUME 定义匿名卷
    格式为：

        VOLUME ["<路径1>", "<路径2>"...]
        VOLUME <路径>
    在 Dockerfile 中，我们可以事先指定某些目录挂载为匿名卷，这样在运行时如果用户不指定挂载，其应用也可以正常运行，不会向容器存储层写入大量数据。


    VOLUME /data
    这里的 /data 目录就会在容器运行时自动挂载为匿名卷，任何向 /data 中写入的信息都不会记录进容器存储层，从而保证了容器存储层的无状态化。当然，运行容器时可以覆盖这个挂载设置。比如：


    $ docker run -d -v mydata:/data xxxx
    在这行命令中，就使用了 mydata 这个命名卷挂载到了 /data 这个位置，替代了 Dockerfile 中定义的匿名卷的挂载配置。

EXPOSE 声明端口
    格式为 EXPOSE <端口1> [<端口2>...]。

    EXPOSE 指令是声明容器运行时提供服务的端口，这只是一个声明，在容器运行时并不会因为这个声明应用就会开启这个端口的服务。在 Dockerfile 中写入这样的声明有两个好处，一个是帮助镜像使用者理解这个镜像服务的守护端口，以方便配置映射；另一个用处则是在运行时使用随机端口映射时，也就是 docker run -P 时，会自动随机映射 EXPOSE 的端口。

    要将 EXPOSE 和在运行时使用 -p <宿主端口>:<容器端口> 区分开来。-p，是映射宿主端口和容器端口，换句话说，就是将容器的对应端口服务公开给外界访问，而 EXPOSE 仅仅是声明容器打算使用什么端口而已，并不会自动在宿主进行端口映射。

WORKDIR 指定工作目录

    格式为 WORKDIR <工作目录路径>。

    使用 WORKDIR 指令可以来指定工作目录（或者称为当前目录），以后各层的当前目录就被改为指定的目录，如该目录不存在，WORKDIR 会帮你建立目录。
    如果需要改变以后各层的工作目录的位置，那么应该使用 WORKDIR 指令。

        WORKDIR /app
        RUN echo "hello" > world.txt
    如果你的 WORKDIR 指令使用的相对路径，那么所切换的路径与之前的 WORKDIR 有关：

        WORKDIR /a
        WORKDIR b
        WORKDIR c

    RUN pwd
    RUN pwd 的工作目录为 /a/b/c。

USER 指定当前用户

        格式：USER <用户名>[:<用户组>]

    USER 指令和 WORKDIR 相似，都是改变环境状态并影响以后的层。WORKDIR 是改变工作目录，USER 则是改变之后层的执行 RUN, CMD 以及 ENTRYPOINT 这类命令的身份。

    注意，USER 只是帮助你切换到指定用户而已，这个用户必须是事先建立好的，否则无法切换。

        RUN groupadd -r redis && useradd -r -g redis redis
        USER redis
        RUN [ "redis-server" ]
    如果以 root 执行的脚本，在执行期间希望改变身份，比如希望以某个已经建立好的用户来运行某个服务进程，不要使用 su 或者 sudo，这些都需要比较麻烦的配置，而且在 TTY 缺失的环境下经常出错。建议使用 gosu (opens new window)。

    # 建立 redis 用户，并使用 gosu 换另一个用户执行命令
    RUN groupadd -r redis && useradd -r -g redis redis
    # 下载 gosu
    RUN wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/1.12/gosu-amd64" \
        && chmod +x /usr/local/bin/gosu \
        && gosu nobody true
    # 设置 CMD，并以另外的用户执行
    CMD [ "exec", "gosu", "redis", "redis-server" ]

HEALTHCHECK 健康检查
    格式：

        HEALTHCHECK [选项] CMD <命令>：设置检查容器健康状况的命令
        HEALTHCHECK NONE：如果基础镜像有健康检查指令，使用这行可以屏蔽掉其健康检查指令
    HEALTHCHECK 指令是告诉 Docker 应该如何进行判断容器的状态是否正常，这是 Docker 1.12 引入的新指令。
    HEALTHCHECK 支持下列选项：

        --interval=<间隔>：两次健康检查的间隔，默认为 30 秒；
        --timeout=<时长>：健康检查命令运行超时时间，如果超过这个时间，本次健康检查就被视为失败，默认 30 秒；
        --retries=<次数>：当连续失败指定次数后，则将容器状态视为 unhealthy，默认 3 次。
    和 CMD, ENTRYPOINT 一样，HEALTHCHECK 只可以出现一次，如果写了多个，只有最后一个生效
    在 HEALTHCHECK [选项] CMD 后面的命令，格式和 ENTRYPOINT 一样，分为 shell 格式，和 exec 格式。命令的返回值决定了该次健康检查的成功与否：0：成功；1：失败；2：保留，不要使用这个值。。
    假设我们有个镜像是个最简单的 Web 服务，我们希望增加健康检查来判断其 Web 服务是否在正常工作，我们可以用 curl 来帮助判断，其 Dockerfile 的 HEALTHCHECK 可以这么写：

        FROM nginx
        RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
        HEALTHCHECK --interval=5s --timeout=3s \
        CMD curl -fs http://localhost/ || exit 1
    这里我们设置了每 5 秒检查一次（这里为了试验所以间隔非常短，实际应该相对较长），如果健康检查命令超过 3 秒没响应就视为失败，并且使用 curl -fs http://localhost/ || exit 1 作为健康检查命令。

    使用 docker build 来构建这个镜像：

        $ docker build -t myweb:v1 .
    构建好了后，我们启动一个容器：

        $ docker run -d --name web -p 80:80 myweb:v1
    当运行该镜像后，可以通过 docker container ls 看到最初的状态为 (health: starting)：

        $ docker container ls
        CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                            PORTS               NAMES
        03e28eb00bd0        myweb:v1            "nginx -g 'daemon off"   3 seconds ago       Up 2 seconds (health: starting)   80/tcp, 443/tcp     web
    在等待几秒钟后，再次 docker container ls，就会看到健康状态变化为了 (healthy)：

        $ docker container ls
        CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                    PORTS               NAMES
        03e28eb00bd0        myweb:v1            "nginx -g 'daemon off"   18 seconds ago      Up 16 seconds (healthy)   80/tcp, 443/tcp     web
    如果健康检查连续失败超过了重试次数，状态就会变为 (unhealthy)。

    为了帮助排障，健康检查命令的输出（包括 stdout 以及 stderr）都会被存储于健康状态里，可以用 docker inspect 来查看。

        $ docker inspect --format '{{json .State.Health}}' web | python -m json.tool
        {
            "FailingStreak": 0,
            "Log": [
                {
                    "End": "2016-11-25T14:35:37.940957051Z",
                    "ExitCode": 0,
                    "Output": "<!DOCTYPE html>\n<html>\n<head>\n<title>Welcome to nginx!</title>\n<style>\n    body {\n        width: 35em;\n        margin: 0 auto;\n        font-family: Tahoma, Verdana, Arial, sans-serif;\n    }\n</style>\n</head>\n<body>\n<h1>Welcome to nginx!</h1>\n<p>If you see this page, the nginx web server is successfully installed and\nworking. Further configuration is required.</p>\n\n<p>For online documentation and support please refer to\n<a href=\"http://nginx.org/\">nginx.org</a>.<br/>\nCommercial support is available at\n<a href=\"http://nginx.com/\">nginx.com</a>.</p>\n\n<p><em>Thank you for using nginx.</em></p>\n</body>\n</html>\n",
                    "Start": "2016-11-25T14:35:37.780192565Z"
                }
            ],
            "Status": "healthy"
        }

SHELL 指令
    格式：SHELL ["executable", "parameters"]

    SHELL 指令可以指定 RUN ENTRYPOINT CMD 指令的 shell环境，Linux 中默认为 ["/bin/sh", "-c"]
    两个 RUN 运行同一命令，第二个 RUN 运行的命令会打印出每条命令并当遇到错误时退出。

    当 ENTRYPOINT CMD 以 shell 格式指定时，SHELL 指令所指定的 shell 也会成为这两个指令的 shell


        SHELL ["/bin/sh", "-cex"]

        # /bin/sh -cex "nginx"
        ENTRYPOINT nginx

        SHELL ["/bin/sh", "-cex"]

        # /bin/sh -cex "nginx"
        CMD nginx

ONBUILD
    格式：ONBUILD <其它指令>。

    ONBUILD 是一个特殊的指令，它后面跟的是其它指令，比如 RUN, COPY 等，而这些指令，在当前镜像构建时并不会被执行。只有当以当前镜像为基础镜像，去构建下一级镜像的时候才会被执行。

    Dockerfile 中的其它指令都是为了定制当前镜像而准备的，唯有 ONBUILD 是为了帮助别人定制自己而准备的。
    我们用 ONBUILD 重新写一下基础镜像的 Dockerfile:

        FROM node:slim
        RUN mkdir /app
        WORKDIR /app
        ONBUILD COPY ./package.json /app
        ONBUILD RUN [ "npm", "install" ]
        ONBUILD COPY . /app/
        CMD [ "npm", "start" ]
    这次我们回到原始的 Dockerfile，但是这次将项目相关的指令加上 ONBUILD，这样在构建基础镜像的时候，这三行并不会被执行。然后各个项目的 Dockerfile 就变成了简单地：

        FROM my-node
    是的，只有这么一行。当在各个项目目录中，用这个只有一行的 Dockerfile 构建镜像时，之前基础镜像的那三行 ONBUILD 就会开始执行，成功的将当前项目的代码复制进镜像、并且针对本项目执行 npm install，生成应用镜像。

docker build时会把当前目录下的文件和目录全部拷贝到服务端

    [root@k8snode2 dockerfile]# ls
    app  copy  dockerfile  hello1  hello2  hello3  hello.tar
    [root@k8snode2 dockerfile]# du -sh 
    233M	.

    开始构建
    [root@k8snode2 dockerfile]# docker build . 
    Sending build context to Docker daemon  243.9MB
    Step 1/4 : FROM centos:latest
    ---> 5d0da3dc9764
    Step 2/4 : WORKDIR /root/
    ---> Using cache
    ---> 6ba106174374
    Step 3/4 : COPY app .
    ---> Using cache
    ---> 4989ab273da2
    Step 4/4 : CMD ["./app"]
    ---> Using cache
    ---> 4a7c708a1930
    Successfully built 4a7c708a1930


dockerfile文件示例:

    [root@k8snode2 copy]# cat dockerfile.copy 
    FROM alpine:latest
    RUN apk --no-cache add ca-certificates
    WORKDIR /root/
    COPY app .
    CMD ["./app"]
在/目录下构建，指定dockerfile文件，会将context下的文件和目录都上传到服务端，因为在context找不到app文件，所以构建失败

    [root@k8snode2 ~]# docker build -f dockerfile/copy/dockerfile.copy  .
    Sending build context to Docker daemon  362.5MB
    Step 1/5 : FROM alpine:latest
    ---> c059bfaa849c
    Step 2/5 : RUN apk --no-cache add ca-certificates
    ---> Using cache
    ---> a89cd93d3a4f
    Step 3/5 : WORKDIR /root/
    ---> Using cache
    ---> c12a350e5018
    Step 4/5 : COPY app .
    COPY failed: stat /var/lib/docker/tmp/docker-builder430395421/app: no such file or directory
    [root@k8snode2 ~]# du -sh
    346M	.

使用dockerfile构建镜像时，最好是新建一个目录，将所需文件全部拷贝到新建目录下，然后再开始构建

docker build -t go/helloworld:1 -f Dockerfile.one .

-t: 指定构建后镜像的name和TAG
-f: 指定dockerfile文件

.dockerignore实践
[root@k8snode2 ms]# ls
bin  count.sh  dockerfile  lib64  mariadb-5.5.68-1.el7.x86_64.rpm  mariadb-libs-5.5.68-1.el7.x86_64.rpm  usr  yum
[root@k8snode2 ms]# cat .dockerignore 
bin
lib64
usr
yum
当docker build时不会上传bin,usr,lib64目录，不会上传yum文件