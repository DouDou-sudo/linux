#!/bin/bash
nodes=`grep "clusters_node" ../vars/main.yml | grep -v "^#" | awk -F"[" '{print $2}' | awk -F"]" '{print $1}'`
echo $nodes
for node in nodes
do
    #node_array+=(node)
    #echo ${node_array[@]}
    echo $node
done
