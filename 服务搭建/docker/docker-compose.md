一、安装docker-compose
下载

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
