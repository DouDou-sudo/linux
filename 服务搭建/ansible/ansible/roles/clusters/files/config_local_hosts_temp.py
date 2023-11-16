#!/usr/bin/python
# -*- coding: utf8 -*-
import commands
cmd = '''grep "clusters_node" ../vars/main.yml | grep -v "^#"| awk -F':' '{print $2}'  '''
#cmd = '''grep "clusters_node" ../vars/main.yml | grep -v "^#" |  awk -F"[" '{print $2}' | awk -F"]" '{print$1}' '''
print(cmd)
datalist1=[]
rc,datalist2 = commands.getstatusoutput(cmd)
datalist=datalist1 + datalist2
print(datalist)
for data in datalist:
  print(data)
#data.strip('\"')
#with open("../vars/main.yml") as f:
#  text = f.readlines()
#  for line in text:
#    if "clusters_node" in line:
#      print(line).strip('\n').

