#!/bin/bash
read -p "请输入kms-server端ip地址:" serverip
read -p "请输入ip地址范围:" IP
#read -p "请输入服务号:" serverUid
read -p "请输入客户端root密码:" pass

for i in `seq $IP`
do
  ip=192.168.84.$i
  if [ $ip == $serverip ];then
     continue
  fi
  ping -c 1 $ip 2>&1 >>/dev/null
  if [ $? -eq 0 ];then
    echo "192.168.84.$i" >> /kms/hosts
  fi
done

#指定被控制端的密码
#pass=Kylinsec.123
#检查sshpass并确认安装，确认配置好yum源
rpm -q sshpass &>/dev/null || yum -y install sshpass
#生成公私钥对
[ -f /root/.ssh/id_rsa ] || ssh-keygen -f /root/.ssh/id_rsa -P ''
#从hosts文件中读取ip地址，然后配置互信
awk '{print $0}' /kms/hosts  |while read line
do
	sshpass -p $pass ssh-copy-id -o stricthostkeychecking=no $line
done


ansible all -i /kms/hosts -m shell -a 'ksl-register -r -d 192.168.84.6  --service-uid TX3L-17189 -s'
ansible all -i /kms/hosts -m file -a 'path=/home/.kms state=directory'
ansible all -i /kms/hosts -m shell -a 'cp -rp /usr/share/ks-license/license /home/.kms/license'
ansible all -i /kms/hosts -m copy -a 'src=/kms/license.sh dest=/home/.kms/ mode=777 backup=yes'
#ansible all -i /kms/hosts -m shell -a 'echo "sleep 30" >>/etc/rc.local'
#ansible all -i /kms/hosts -m shell -a 'echo "sh /home/.kms/license.sh" >>/etc/rc.local'
