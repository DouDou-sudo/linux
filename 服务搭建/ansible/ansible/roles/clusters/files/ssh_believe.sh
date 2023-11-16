#!/bin/bash
#安装及检测sshpass
rpm -q sshpass &>/dev/null || yum -y install sshpass
#生成公私钥对
[ -f /root/.ssh/id_rsa ] || ssh-keygen -f /root/.ssh/id_rsa -P ''
#从hosts文件中读取ip地址，然后配置互信
pass=root.2020
#hosts=("192.168.189.160" "192.168.189.100")
#for host in ${hosts[@]}
#awk '{print $0}' hosts  |while read host
awk 'NR>2{print $1}' /etc/hosts | while read host
do
	sshpass -p $pass ssh-copy-id -o stricthostkeychecking=no $host
done
