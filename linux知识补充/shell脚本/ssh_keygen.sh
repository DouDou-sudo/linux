#!/bin/bash
echo "敲两下空格"
ssh-keygen -t rsa -f id_rsa
cd /root/.ssh
（1）cp id_rsa.pub authorized_keys
（2）ssh-copy-id root@localhost
ssh-add  id_rsa
（1）read -p "输入另一台ha地址：" node2IP
scp -p authorized_keys id_rsa root@$node2IP:/root/.ssh
（2） for i in `seq 1 3`;do scp -rp .ssh/ 192.168.189.12$i:/home/cephadm/; done






