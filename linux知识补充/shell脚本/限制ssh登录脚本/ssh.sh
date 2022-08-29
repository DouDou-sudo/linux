#!/bin/bash
#a.txt
awk 'NR>2{print $1}' /etc/hosts > a.txt
#for ip in `cat a.txt`
#do
#echo "sshd:$line:allow">>/etc/hosts.allow
#done
cat a.txt | while read line
do
echo "sshd:$line:allow">>/etc/hosts.allow
done
echo "sshd:all" >> /etc/hosts.deny
systemctl restart sshd

#while read line
#do
#	echo "sshd:$line:allow">>/etc/hosts.allow
#done < a.txt
