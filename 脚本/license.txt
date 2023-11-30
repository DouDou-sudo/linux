#!/bin/bash
Pid=$(ps -ef | grep ksl-daemon | grep -v "grep" | awk '{print $2}')
kill -9 $Pid
if [ $? -eq 0 ];then
  ln -sf /home/.kms/license /usr/share/ks-license/license
  echo "success">>/var/log/license.log
else
  echo "kill ksl-daemon failed">>/var/log/license.log
  exit 1
fi
