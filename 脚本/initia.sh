#!/bin/bash
read -p "请输入主机名:" hostName
hostnamectl set-hostname $hostName


echo "---------设置防火墙----------------------"
firewall(){
  setenforce 0
  sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
  systemctl disable --now firewalld
}

read -p "是否关闭防火墙:" Fw
case $Fw in
YES|yes|y|Y)
  firewall
  ;;
NO|no|n|N)
  echo "不修改防火墙默认设置"
  ;;
*)
  echo "请输入YES|yes|y|Y|NO|no|n|N";;
esac

echo "-------------设置swap--------------------"
swp(){
  swapoff -a
  sed -i.bak '/swap/s/^/#/g' /etc/fstab
}
read -p "是否禁用swap:" Sp
case $Sp in
YES|yes|y|Y)
  swp
  ;;
NO|no|n|N)
  echo "不禁用swap"
  ;;
*)
  echo "请输入YES|yes|y|Y|NO|no|n|N";;
esac


echo "--------------设置dns-------------------"
read -p "请输入DNS地址:" inps

ipct=`echo $inps|awk -F "." '{print NF}'`

if [ $ipct = 4 ]
then
count=0
  for ((i=1; i<=4; i++))
  do
    ad="$"${i}
    ipctsep=`echo $inps | awk -F "." '{print '${ad}'}'`
    if [[ $ipctsep =~ [0-255] ]];then
       count=`expr $count + 1 `
    fi
  done
  if [ $count -eq 4 ];then
    echo "nameserver $inps" >>/etc/resolv.conf
  else
    echo "ip不合规"
  fi
else
  echo "ip不合规"
fi

echo "--------------配置yum源-------------------"
mkdir /etc/yum.repos.d/bak
mv /etc/yum.repos.d/* /etc/yum.repos.d/bak/ 
curl -o /etc/yum.repos.d/Centos-7.repo   https://mirrors.aliyun.com/repo/Centos-7.repo
curl -o /etc/yum.repos.d/epel-7.repo   https://mirrors.aliyun.com/repo/epel-7.repo
curl -o /etc/yum.repos.d/docker-ce.repo   https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum clean all && yum makecache

echo "---------------------------------"
yum install -y yum-utils device-mapper-persistent-data lvm2 wegt net-tools vim sshpass

echo "--------------设置时区并设置时间同步-------------------"
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 
timedatectl set-timezone 'Asia/Shanghai'
sed -i '/server/s/^/#/g' /etc/chrony.conf
echo "server cn.pool.ntp.org" >> /etc/chrony.conf
systemctl enable --now chronyd
systemctl restart chronyd
hwclock -w

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
