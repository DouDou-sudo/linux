一、安装docker-compose
1、rpm安装
$ yum install -y docker-compose-plugin
 验证程序是否正常
$ docker compose version
Docker Compose version v2.6.0
2、二进制安装
    curl -L https://get.daocloud.io/docker/compose/releases/download/v2.9.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
添加执行权限

    chmod +x /usr/local/bin/docker-compose

二、用python来建立一个能够记录页面访问此时的web网站

新建一个目录
web应用

    [root@zabbix-server docker-compose]# cat app.py 
    #!/usr/bin/python
    from flask import Flask
    from redis import Redis

    app = Flask(__name__)
    redis = Redis(host='redis', port=6379)

    @app.route('/')
    def hello():
        count = redis.incr('hits')
        return 'Hello World! 该页面已被访问 {} 次。\n'.format(count)

    if __name__ == "__main__":
        app.run(host="0.0.0.0", debug=True)

dockerfile 

    [root@zabbix-server docker-compose]# cat dockerfile
    FROM python:3.7
    ADD . /code
    WORKDIR /code
    RUN pip install -i https://pypi.tuna.tsinghua.edu.cn/simple redis flask
    CMD ["python","app.py"]


docker-compose.yml

    [root@zabbix-server docker-compose]# cat docker-compose.yml 
    version: '3'
    services:

    web:
        build: .
        ports:
        - "5000:5000"

    redis:
        image: "redis:alpine"

运行compose项目

    [root@zabbix-server docker-compose]# docker-compose  up
    [root@zabbix-server docker-compose]# docker-compose  up -d
    [+] Building 17.6s (9/9) FINISHED                                                                                                                                                                                
    => [internal] load build definition from dockerfile                                                                                                                                                        0.0s
    => => transferring dockerfile: 237B                                                                                                                                                                        0.0s
    => [internal] load .dockerignore                                                                                                                                                                           0.0s
    => => transferring context: 2B                                                                                                                                                                             0.0s
    => [internal] load metadata for docker.io/library/python:3.7                                                                                                                                               0.0s
    => [internal] load build context                                                                                                                                                                           0.0s
    => => transferring context: 1.03kB                                                                                                                                                                         0.0s
    => CACHED [1/4] FROM docker.io/library/python:3.7                                                                                                                                                          0.0s
    => [2/4] ADD . /code                                                                                                                                                                                       0.0s
    => [3/4] WORKDIR /code                                                                                                                                                                                     0.0s
    => [4/4] RUN pip install -i https://pypi.tuna.tsinghua.edu.cn/simple redis flask                                                                                                                          17.3s
    => exporting to image                                                                                                                                                                                      0.2s
    => => exporting layers                                                                                                                                                                                     0.2s
    => => writing image sha256:3fd45ee46c57b2698e99c5c7d80ddfbb5e92df93584f2dc73d7b0c8fde2faf25                                                                                                                0.0s 
    => => naming to docker.io/library/docker-compose-web                                                                                                                                                       0.0s 
                                                                                                                                                                                                                    
    Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them                                                                                                             
    [+] Running 3/1
    ⠿ Network docker-compose_default    Created                                                                                                                                                                0.1s
    ⠿ Container docker-compose-web-1    Created                                                                                                                                                                0.1s
    ⠿ Container docker-compose-redis-1  Created                                                                                                                                                                0.1s
    Attaching to docker-compose-redis-1, docker-compose-web-1
    docker-compose-redis-1  | 1:C 18 Nov 2022 08:01:24.478 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
    docker-compose-redis-1  | 1:C 18 Nov 2022 08:01:24.478 # Redis version=6.2.6, bits=64, commit=00000000, modified=0, pid=1, just started
    docker-compose-redis-1  | 1:C 18 Nov 2022 08:01:24.478 # Warning: no config file specified, using the default config. In order to specify a config file use redis-server /path/to/redis.conf
    docker-compose-redis-1  | 1:M 18 Nov 2022 08:01:24.480 * monotonic clock: POSIX clock_gettime
    docker-compose-redis-1  | 1:M 18 Nov 2022 08:01:24.481 * Running mode=standalone, port=6379.
    docker-compose-redis-1  | 1:M 18 Nov 2022 08:01:24.481 # WARNING: The TCP backlog setting of 511 cannot be enforced because /proc/sys/net/core/somaxconn is set to the lower value of 128.
    docker-compose-redis-1  | 1:M 18 Nov 2022 08:01:24.481 # Server initialized
    docker-compose-redis-1  | 1:M 18 Nov 2022 08:01:24.481 # WARNING overcommit_memory is set to 0! Background save may fail under low memory condition. To fix this issue add 'vm.overcommit_memory = 1' to /etc/sys
    ctl.conf and then reboot or run the command 'sysctl vm.overcommit_memory=1' for this to take effect.docker-compose-redis-1  | 1:M 18 Nov 2022 08:01:24.481 * Ready to accept connections
    docker-compose-web-1    |  * Serving Flask app 'app'
    docker-compose-web-1    |  * Debug mode: on
    docker-compose-web-1    | WARNING: This is a development server. Do not use it in a production deployment. Use a production WSGI server instead.
    docker-compose-web-1    |  * Running on all addresses (0.0.0.0)
    docker-compose-web-1    |  * Running on http://127.0.0.1:5000
    docker-compose-web-1    |  * Running on http://172.18.0.2:5000
    docker-compose-web-1    | Press CTRL+C to quit
    docker-compose-web-1    |  * Restarting with stat
    docker-compose-web-1    |  * Debugger is active!
    docker-compose-web-1    |  * Debugger PIN: 845-485-003

    docker-compose-web-1    | 192.168.189.1 - - [18/Nov/2022 08:05:48] "GET / HTTP/1.1" 200 -
    docker-compose-web-1    | 192.168.189.1 - - [18/Nov/2022 08:05:48] "GET /favicon.ico HTTP/1.1" 404 -
    docker-compose-web-1    | 192.168.189.1 - - [18/Nov/2022 08:06:12] "GET / HTTP/1.1" 200 -
    docker-compose-web-1    | 192.168.189.1 - - [18/Nov/2022 08:06:13] "GET / HTTP/1.1" 200 -


重启一个终端，测试

    [root@zabbix-server docker-compose]# curl 127.0.0.1:5000
    Hello World! 该页面已被访问 10 次。
查看启动的容器

[root@zabbix-server docker-compose]# docker ps
CONTAINER ID   IMAGE                                             COMMAND                  CREATED          STATUS          PORTS                                             NAMES
af924406ab84   docker-compose-web                                "python app.py"          22 seconds ago   Up 21 seconds   0.0.0.0:5000->5000/tcp, :::5000->5000/tcp         docker-compose-web-1
f1ad9354e29a   redis:alpine                                      "docker-entrypoint.s…"   22 seconds ago   Up 21 seconds   6379/tcp                                          docker-compose-redis-1

Docker Compose使用基于yml格式的模板文件来定义应用，默认的模板文件名称为docker-compose.yml。
我们先来了解一下模板的格式 ，如下是一个模板示例，它定义了两个应用服务。
version: "3.9"
services:
  myapp:
    build: ./myapp
    ports:
      - "80:80"
    container_name: myapp
    depends_on:
      - database
    networks:
      - mynet
  database:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: 123456
    ports:
      - "3306:3306"
    container_name: db
    volumes:
      - "mydata:/var/lib/mysql"
    networks:
      - mynet
networks:
  mynet:
    driver: bridge

volumes:
  mydata:
下面我们根据模板的各个功能参数来进行拆分讲解：
version 

定义了Compose文件使用的API版本（并非Compose产品自身版本），一般建议使用最新的版本，此处使用3.9。

services

用于定义不同的应用服务，此处定义了两个服务：web和database。Docker Compose会将每个服务部署在各自的容器中。

build

指定Dockerfile所在文件夹的路径（可支持绝对路径或相对路径），Compose 会使用它来构建镜像，该镜像会被用于启动该服务的容器。

ports

指定Docker容器的端口映射到主机，该功能与docker run  启动容器的 -p 选项一样，格式为“主机端口：容器端口”。如果只指定了容器端口，那么主机端口将使用随机端口。

container_name

用于指定容器名称，如果不指定的话，则默认会使用“目录名+服务名+序号”这样的格式来表示。

depends_on

指定服务间的依赖关系。如本例所示，由于myapp依赖于database服务，所以Docker Compose会先启动database服务对应的容器，然后再进行myapp的启动。

networks

一级的networks参数用于指定创建服务所使用的网络，此处指定名称为mynet，Compose会在该名称前面再加上目录名称作为创建的网络名称，格式为“目录名_mynet”。网络驱动使用bridge，表明这是一个bridge网络。

services 项中的服务使用networks用于指定加入的网络，此处都使用mynet网络。

image

指定服务使用的镜像。本示例中的database服务使用mysql:5.7的镜像，Compose会自动下载该镜像用于启动容器。

environment

设置环境变量，该变量会在容器启动时自动加载。如果只指定了变量名称，则会自动获取Compose主机上对应变量的值，这种方式可用来防止重要数据的泄漏。

示例：

environment:
  api_key:
volumes
指定容器挂载的数据卷，可以使用目录挂载或volumes挂载。

示例：目录挂载

    volumes:
      - "/data/mysql:/var/lib/mysql"
示例：volumes挂载
volumes:
  - mydata:/var/lib/mysql
  
volumes:
  mydata:
除了以上这些常用的参数外，Docker Compose还有不少其他的参数可以使用，限于篇幅原因在此不做过多介绍，有兴趣的读者可自己查看相关资料。
三. 管理应用
Docker Compose在新的版本中，使用的命令格式与其他docker 命令类似，为“ docker compose + 选项 + [参数]”。原有的“docker-compose 命令不再使用。

1. 查看帮助

执行"docker compose --help" 可查看命令的详细介绍

$ docker compose --help
......
Commands:
  build       Build or rebuild services
  convert     Converts the compose file to platform's canonical format
  cp          Copy files/folders between a service container and the local filesystem
  create      Creates containers for a service.
  down        Stop and remove containers, networks
  events      Receive real time events from containers.
  exec        Execute a command in a running container.
  images      List images used by the created containers
  kill        Force stop service containers.
  logs        View output from containers
......

2. 构建镜像

使用docker compose build 命令可用于构建项目中镜像，我们可以使用此命令先提前准备好相关的镜像。否则 ，相关镜像将在启动服务时进行构建。

$ docker compose build
[+] Building 3.3s (2/3)                                                                                                                                      
 => [internal] load build definition from Dockerfile                                  0.0s
 => => transferring dockerfile: 92B                                                   0.0s
 => [internal] load .dockerignore                                                     0.0s
 => => transferring context: 2B                                                       0.0s
 => [internal] load metadata for docker.io/library/python:3.7-alpine  
 ......
构建完成后，我们可以看到镜像已经生成，镜像名称格式为“目录名+服务名”。
$ docker compose images
Container           Repository          Tag                 Image Id            Size
myapp               test_myapp          latest              fda100c4b2ae        51.6MB

3. 服务创建并启动

使用docker compose up 命令可启动docker-compose.yml 模板文件中定义的服务，通常会在命令后面加上 -d ，用于后台启动。该命令会构建所需的镜像、创建网络和存储卷、然后进行容器的启动等操作。

由于前面已提前构建好myapp的镜像，故这里只需要拉取database的镜像，并创建相关的容器资源即可。

$ docker compose up -d
[+] Running 12/12
 ⠿ database Pulled                                                                                                                                     98.0s
   ⠿ 72a69066d2fe Pull complete                                                                                                                        35.4s
   ⠿ 93619dbc5b36 Pull complete                                                                                                                        35.5s
   ⠿ 99da31dd6142 Pull complete                                                                                                                        35.7s
   ⠿ 626033c43d70 Pull complete                                                                                                                        35.9s
   ⠿ 37d5d7efb64e Pull complete                                                                                                                        35.9s
   ⠿ ac563158d721 Pull complete                                                                                                                        36.8s
   ⠿ d2ba16033dad Pull complete                                                                                                                        36.9s
   ⠿ 0ceb82207cd7 Pull complete                                                                                                                        36.9s
   ⠿ 37f2405cae96 Pull complete                                                                                                                        80.7s
   ⠿ e2482e017e53 Pull complete                                                                                                                        80.7s
   ⠿ 70deed891d42 Pull complete                                                                                                                        80.8s
[+] Running 4/4
 ⠿ Network test_mynet    Created                                                                                                                        0.2s
 ⠿ Volume "test_mydata"  Created                                                                                                                        0.0s
 ⠿ Container db          Started                                                                                                                        0.5s
 ⠿ Container myapp       Started
查看容器状态，可看到服务容器已正常启动。

docker compose ps
NAME                COMMAND                  SERVICE             STATUS              PORTS
db                  "docker-entrypoint.s…"   database            running             0.0.0.0:3306->3306/tcp, :::3306->3306/tcp
myapp               "python /opt/myapp/m…"   myapp               running             0.0.0.0:80->80/tcp, :::80->80/tcp

4.  停止/启动服务
使用docker compose stop/start 命令可对服务相关的容器进行停止/启动。

停止服务

$ docker compose stop 
[+] Running 2/2
 ⠿ Container myapp  Stopped                    10.2s
 ⠿ Container db     Stopped
启动服务
$ docker compose start
[+] Running 2/2
 ⠿ Container db     Started                     0.4s
 ⠿ Container myapp  Started

5. 服务下线

使用docker compose down 会将服务相关的容器与网络进行停止并清理删除，只保留镜像、存储卷等资源。

$ docker compose down
[+] Running 3/3
 ⠿ Container myapp     Removed                             10.2s
 ⠿ Container db        Removed                             1.9s
 ⠿ Network test_mynet  Removed