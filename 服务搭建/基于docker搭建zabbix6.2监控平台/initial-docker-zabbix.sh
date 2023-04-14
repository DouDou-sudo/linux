#!/bin/bash
echo "-------------安装docker和docker-compose--------------------"
yum install -y docker-ce-20.10.10
systemctl enable --now docker
cat >>/etc/docker/daemon.json <<EOF
{
  "registry-mirrors": ["https://v16stybc.mirror.aliyuncs.com"],
  "exec-opts": ["native.cgroupdriver=systemd"]
}
EOF
systemctl daemon-reload
systemctl restart docker
curl -L https://get.daocloud.io/docker/compose/releases/download/v2.9.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

echo "-------------创建目录--------------------"
mkdir -p /docker/images
mkdir -p /docker/yaml
echo "-------------拷贝镜像和yaml文件--------------------"
cp -r images/* /docker/images
cp -r yaml/* /docker/yaml
echo "-------------导入镜像--------------------"
for i in `ls /docker/images`;do  docker load -i /docker/images/$i;done
echo "------------创建持久化存储卷--------------------"
docker volume create mysql_data
docker volume create mysql_logs
docker volume create mysql_conf
docker volume create zabbix_server
docker volume create zabbix_agent

echo "-------------启动docker-compose--------------------"
docker-compose -f /docker/yaml/docker-compose-zabbix.yml -d 

echo "请使用浏览器访问主机ip地址"
echo "login:Admin"
echo "passwd:zabbix"
