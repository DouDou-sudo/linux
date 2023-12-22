### 一、问题发现
客户告知一台服务器于12月19日异常，可以ping通，ssh不通，20日强制重启恢复
查看messages日志
`系统于20日11:57:27按电源键强制重启,12:04:56操作系统开始加载，12:05:14 rsyslogd服务拉起`

    Dec 20 11:57:27 localhost systemd-logind: Power key pressed.
    Dec 20 12:04:56 localhost kernel: Initializing cgroup subsys cpuset
    Dec 20 12:04:56 localhost kernel: Initializing cgroup subsys cpu
    Dec 20 12:04:56 localhost kernel: Initializing cgroup subsys cpuacct
    ...
    Dec 20 12:04:56 localhost kernel: Kernel command line: BOOT_IMAGE=/vmlinuz-3.10.0-957.ky3.kb10.x86_64 root=/dev/mapper/unikylin-root ro crashkernel=auto rd.lvm.lv=unikylin/root rd.lvm.lv=unikylin/swap rhgb quiet LANG=zh_CN.UTF-8
    ...
    Dec 20 12:05:14 localhost rsyslogd: [origin software="rsyslogd" swVersion="8.24.0-34.ky3.kb2" x-pid="40690" x-info="http://www.rsyslog.com"] start
    ...
`系统强制重启前一直连续打印如下信息:`

    Dec 19 10:51:14 localhost systemd: Started Session 820 of user root.
    Dec 19 10:52:44 localhost kernel: INFO: task kworker/35:0:149580 blocked for more than 120 seconds.
    Dec 19 10:52:44 localhost kernel: "echo 0 > /proc/sys/kernel/hung_task_timeout_secs" disables this message.
    Dec 19 10:52:44 localhost kernel: kworker/35:0    D ffff9b5c6db65140     0 149580      2 0x00000080
    Dec 19 10:52:44 localhost kernel: Call Trace:
    Dec 19 10:52:44 localhost kernel: [<ffffffff93d6cd09>] schedule+0x29/0x70
    Dec 19 10:52:44 localhost kernel: [<ffffffff93d6a7e1>] schedule_timeout+0x221/0x2d0
    ...
### 二、排查问题(google后...)
该问题为一个内核bug，一般会在超大内存的服务器上触发此bug
相关链接:
https://blog.csdn.net/weixin_31845243/article/details/116968817

https://blog.csdn.net/electrocrazy/article/details/79377214

https://www.cnblogs.com/wshenjin/p/7093505.html
**问题原因:**
By default Linux uses up to 40% of the available memory for file system caching.
After this mark has been reached the file system flushes all outstanding data to disk causing all following IOs going synchronous.
For flushing out this data to disk this there is a time limit of 120 seconds by default.
In the case here the IO subsystem is not fast enough to flush the data withing 120 seconds.
This especially happens on systems with a lot of memory.
The problem is solved in later kernels。

翻译过来就是：一般情况下，linux会把可用内存的40%的空间作为文件系统的缓存。当缓存快满时，文件系统将缓存中的数据整体同步到磁盘中。但是系统对同步时间有最大120秒的限制。如果文件系统不能在时间限制之内完成数据同步，则会发生上述的错误。这通常发生在内存很大的系统上。系统内存大，则缓冲区大，同步数据所需要的时间就越长，超时的概率就越大。

默认情况下， Linux最多会使用40%的可用内存作为文件系统缓存。当超过这个阈值后，文件系统会把将缓存中的内存全部写入磁盘， 导致后续的IO请求均同步。将缓存写入磁盘的时候，有一个默认120秒的超时时间。 出现上面的问题是因为IO子系统的处理不够快速，不能在120秒将缓存中的数据全部写入到磁盘的。IO系统响应缓慢，导致越来越多的请求堆积，最终系统内存全部被占用，导致系统失去响应。

### 三、解决办法

#### 1、缩小文件系统缓存大小
此种方案是降低缓存占内存的比例，比如由40%降到10%，这样的话需要同步到磁盘上的数据量会变小，IO写时间缩短，会相对比较平稳。

文件系统缓存的大小是由内核参数vm.dirty_ratio 和 vm.dirty_backgroud_ratio控制决定的。

vm.dirty_background_ratio指定当文件系统缓存脏页数量达到系统内存百分之多少时（如5%）就会触发pdflush/flush/kdmflush等后台回写进程运行，将一定缓存的脏页异步地刷入外存。

vm.dirty_ratio则指定了当文件系统缓存脏页数量达到系统内存百分之多少时（如10%），系统不得不开始处理缓存脏页（因为此时脏页数量已经比较多，为了避免数据丢失需要将一定脏页刷入外存），在此过程中很多应用进程可能会因为系统转而处理文件IO而阻塞。

通常情况下，vm.dirty_ratio的值要大于vm.dirty_background_ratio的值。

##### 1）查看系统当前内核参数
sysctl -a | grep dirty                  
> 如果系统没有设置相关参数,可以查看下面2个文件

cat /proc/sys/vm/dirty_ratio 
cat /proc/sys/vm/dirty_background_ratio
##### 2）临时修改内核参数
sysctl -w vm.dirty_ratio=10
sysctl -w vm.dirty_background_ratio=5
##### 3）永久修改
vi /etc/sysctl.conf打开文件添加2行
vm.dirty_background_ratio = 5
vm.dirty_ratio = 10
执行这个命令生效sysctl -p
##### 4）查看系统当前内核参数是否修改成功
sysctl -a | grep dirty
#### 2、取消120秒时间限制
此方案就是不让系统有那个120秒的时间限制。文件系统把数据从缓存转到外存慢点就慢点，应用程序对此延时不敏感。就是慢点就慢点，我等着。实际上操作系统是将这个变量设为长整形的最大值。

下面说一下内核hung task检测机制由来。我们知道进程等待IO时，经常处于D状态，即TASK_UNINTERRUPTIBLE状态，处于这种状态的进程不处理信号，所以kill不掉，如果进程长期处于D状态，那么肯定不正常，原因可能有二：

1. IO路径上的硬件出问题了，比如硬盘坏了(只有少数情况会导致长期D，通常会返回错误)；

2. 内核自己出问题了。

这种问题不好定位，而且一旦出现就通常不可恢复，kill不掉，通常只能重启恢复了。
内核针对此种情况开发了一种hung task的检测机制，基本原理是：定时检测系统中处于D状态的进程，如果其处于D状态的时间超过了指定时间(默认120s，可以配置)，则打印相关堆栈信息，也可以通过proc参数配置使其直接panic。

如何修改或者取消120秒的时间限值呢。120秒的时间限值由内存参数kernel.hung_task_timeout_secs决定的。直接像方案一那样修改此内核参数的值就可。如果kernel.hung_task_timeout_secs的值设置为0，那就是把此种设置为长整型的最大值。

下面说一下修改调度器的流程。

##### 1）查看系统当前内核参数hung_task_timeout_secs值。

在命令行中输入如下指令：

    [root@node1 ~]# sysctl -a | grep hung_task_timeout_secs
    kernel.hung_task_timeout_secs = 120
有内核返回信息，可知当前设置的hung_task超时时间为120秒。
> 如果系统没有设置相关参数,可以查看下面这个文件

    [root@node1 ~]# cat /proc/sys/kernel/hung_task_timeout_secs 
    120

##### 2）临时修改内核参数hung_task_timeout_secs值。
把hung_task_timeout_secs的值修改为0，在命令行中输入如下指令：

sysctl -w kernel.hung_task_timeout_secs=0
##### 3）永久修改
vi /etc/sysctl.conf打开文件添加1行
kernel.hung_task_timeout_secs=0
执行这个命令生效sysctl -p
##### 4）查看系统当前内核参数是否修改成功
sysctl -a | grep hung_task_timeout_secs 