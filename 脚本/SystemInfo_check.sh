#!/bin/bash
Log_path="/var/log/Sys_check.log"

Max_mem=100
Max_cpu=5
Max_file=200
Max_load=1
Max_idle=10
Max_io=100

disk_info=$(lsblk | awk '{print $1}' | grep -ev "NAME|sr0" | grep -v - | grep -v "─")

declare -i i=0
function my_cmd(){
	echo $i
#30秒内平均io队列非空的时间占比
        sys_io=`iostat -x 1 30 | grep ${disk_info[i]} | awk '{sum+=$14} END {print "",sum/NR}'`

        overstep_io="`date +"%F %T"` warning: sys_io overstep! io = $sys_io %"
#如果$sys_io大于$Max_io为真，则执行echo，否则无任何输出
        if [ `echo "$sys_io > $Max_io" | bc` -eq 1 ];then
                echo $overstep_io >> $Log_path
        fi
	let ++i
}

tmp_fifofile="/tmp/$$.fifo"
mkfifo $tmp_fifofile      # 新建一个fifo类型的文件
exec 6<>$tmp_fifofile     # 将fd6指向fifo类型
rm $tmp_fifofile    #删也可以

thread_num=100  # 最大可同时执行线程数量
job_num=${#disk_info[@]}   # 任务总数

#根据线程总数量设置令牌个数
for ((i=0;i<${thread_num};i++));do
    echo
done >&6

for ((i=0;i<${job_num};i++));do # 任务数量
    # 一个read -u6命令执行一次，就从fd6中减去一个回车符，然后向下执行，
    # fd6中没有回车符的时候，就停在这了，从而实现了线程数量控制
    read -u6

    #可以把具体的需要执行的命令封装成一个函数
    {
        my_cmd
    } &

    echo >&6 # 当进程结束以后，再向fd6中加上一个回车符，即补上了read -u6减去的那个
done

wait
exec 6>&- # 关闭fd6

sys_mem=`vmstat -S m -s | grep "used memory" | awk '{print $1}'`
sys_cpu=`top -b -n 1 | head -n 10 | awk '{print $9}' | tail -n 3 | head -n 1`
#sys_file=`lsof -n | awk '{print $2}' | uniq -c | sort -nr | head -n 1 | awk '{print $1}'`
sys_core=`lscpu | awk '/^CPU\(s\)/{print $2}'`
sys_load=`uptime | awk -F "," '{print $4}'`
#计算$sys_load除以$sys_core的值（保留两位小数），赋值给sys_overload
sys_overload=`echo "scale=2; $sys_load / $sys_core" | bc`
sys_idle=`sar -u 1 1 | awk '{print $8}' | tail -n 1`
sys_lock=`vmstat 1 30 | awk '{print $2}'| tail -n 20 | awk '{sum+=$1} END {print "", sum/NR}'`
sys_swap=`vmstat 1 30 | awk '{print $7 $8}'| tail -n 20 | awk '{sum+=$1} END {print "", sum/NR}'`

overstep_mem="`date "+%Y-%m-%d %H:%M:%S"` warning : sys_mem overstep! used_memory = $sys_mem M"
overstep_cpu="`date "+%Y-%m-%d %H:%M:%S"` warning : sys_cpu overstep! process_cpu = $sys_cpu %"
overstep_file="`date "+%Y-%m-%d %H:%M:%S"` warning: sys_file overstep! file_open = $sys_file"
overstep_load="`date "+%Y-%m-%d %H:%M:%S"` warning: sys_load overstep! load_average = $sys_overload"
overstep_idle="`date "+%Y-%m-%d %H:%M:%S"` warning: sys_idle overstep! idle = $sys_idle %"
overstep_lock="`date "+%Y-%m-%d %H:%M:%S"` warning: sys_lock!"
overstep_swap="`date "+%Y-%m-%d %H:%M:%S"` warning: swap over!"

if [ $sys_mem -gt $Max_mem ];then
	echo $overstep_mem >> $Log_path 
fi

if [ `echo "$sys_cpu > $Max_cpu" | bc` -eq 1 ];then
	echo $overstep_cpu >> $Log_path
fi

if [ $sys_file -gt $Max_file ];then
	echo $overstep_file >> $Log_path
fi

if [ `echo "$sys_overload > $Max_load" | bc` -eq 1 ];then
	echo $overstep_load >> $Log_path
fi

if [ `echo "$sys_idle > $Max_idle" | bc` -eq 1 ];then
	echo null >/dev/null
else
	echo $overstep_idle >> $Log_path
fi

if [ `echo " $sys_swap > 0" | bc` -eq 1 ];then
	echo $overstep_swap >> $Log_path
fi

if [ `echo " $sys_lock > 0" | bc` -eq 1 ];then
	echo $overstep_lock >> $Log_path
fi
