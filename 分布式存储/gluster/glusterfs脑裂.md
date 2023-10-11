###glusterfs脑裂


>裂脑：
裂脑是指文件的两个或多个复制副本变得不同的情况。当文件处于裂脑状态时，副本砖块中文件的数据或元数据不一致，并且没有足够的信息来权威地选择原始副本并修复坏副本，尽管所有砖块都已启动并联机。对于目录，还有一个条目拆分大脑，其中的文件可以在副本的砖块上具有不同的 gfid/文件类型。发生裂脑主要是因为两个原因：

* 由于网络断开连接，客户端暂时失去与brick的连接。
2副本，服务器 1 上的brick 1 和服务器 2 上的brick 2。
 由于网络拆分，客户端 1 失去与brick 2 的连接，客户端 2 失去与brick 1 的连接。
从客户端 1 写入到brick 1，从客户端 2 转到brick 2

* gluster进程下降或返回错误：
服务器 1 关闭，服务器 2 启动：写入发生在服务器 2 上。
服务器 1 启动，服务器 2 关闭（未发生修复/服务器 2 上的数据未在服务器 1 上复制）：写入发生在服务器 1 上。
服务器 2 启动：服务器 1 和服务器 2 都有彼此独立的数据。

**脑裂分为3种**
1. 数据脑裂：文件中的数据在副本组的brick中不同
2. 元数据脑裂：brick中元数据不同
3. GFID脑裂：副本brick上的文件GFID不同，或者副本上的文件类型不同，文件类型不同无法修复，GFID可以修复，GFID脑裂对外表现为目录脑裂


#### 1、gluster cli修复data/metadata split-brain
>data/metadata的split-brain，\<FILE\>可以以指定gfid的方式标识,
\<FILE>\都是相对路径，相对于brick的相对路径，brick目录就是\<FILE>\的/
\<FILE>\可以为目录，当指定目录为源时，后面不能跟/，比如/date,不能写为/date/
```
eg：
gluster volume heal test split-brain source-brick test-host:/test/b1 gfid:c3c94de2-232d-4083-b534-5da17fc476ac
gluster volume heal test split-brain bigger-file /dir/file1
```
###### 选择较大的文件做为源文件
gluster volume heal \<VOLNAME\> split-brain bigger-file \<FILE\>
```shell
eg:gluster v heal replica2 split-brain bigger-file /stbn
No bigger file for file /stbn 如果文件size都一致，会报错
```
###### 选择最新的mtime的文件做为源文件
gluster volume heal \<VOLNAME\> split-brain latest-mtime \<FILE\>
```shell
eg:gluster volume heal replica2 split-brain latest-mtime /stbn
```
###### 选择副本中的一个brick的文件做为源文件
gluster volume heal \<VOLNAME\> split-brain source-brick <HOSTNAME:brick-directory-absolute-path> \<FILE\>
```shell
eg:gluster volume heal replica2 split-brain source-brick glusterfs-node2:/glusterfs/replica/brick1/ /stbn
```
###### 选择副本中的一个brick作为所有文件的源文件
gluster volume heal \<VOLNAME\> split-brain source-brick <HOSTNAME:brick-directory-absolute-path>
```shell
eg:gluster volume heal replica2 split-brain source-brick glusterfs-node2:/glusterfs/replica/brick1/
Healing gfid:00000000-0000-0000-0000-000000000001 failed:Operation not permitted.如果有GFID split-brain会报错不允许的操作
```
#### 2、gluster cli修复GFID split-brain
使用getfattr -d -e hex -m. <path-to-file>查看各brick下的文件权限trusted.gfid是否一致，不一致就是GFID脑裂
>data/metadata的split-brain，\<FILE\>可以以不能指定gfid的方式标识，因为此时gfid各副本不同。
\<FILE>\都是相对路径，相对于brick的相对路径，brick目录就是\<FILE>\的/
\<FILE>\可以为目录，当指定目录为源时，后面不能跟/，比如/date,不能写为/date/
###### 选择较大的文件作为源
gluster volume heal \<VOLNAME\> split-brain bigger-file \<FILE\>
###### 选择最新的mtime的文件作为源
gluster volume heal \<VOLNAME\> split-brain latest-mtime \<FILE\>
###### 选择副本中的一个brick作为特定文件的源
gluster volume heal \<VOLNAME\> split-brain source-brick <HOSTNAME:brick-directory-absolute-path> \<FILE\>

> 注意：
-	不能将文件的 GFID 用作任何 CLI 选项的参数来解决 GFID 裂脑问题。它应该是被视为源的文件的绝对路径。
-	使用brick作为源选项，无法一次性解决所有 GFID 裂脑，因为在解析数据或元数据裂脑时，无需在 CLI 中指定任何文件路径。对于 GFID 裂脑中的每个文件，使用要使用的策略运行 CLI。
-	使用带有“分布式复制”卷中的“brick块”选项的 CLI 解析目录 GFID 裂脑需要在处于此状态的所有子卷上显式完成。由于目录会在所有brick上创建，因此使用一个特定的brick作为目录 GFID 裂脑的源可以修复该特定子卷的目录。源brick的选择方式应使修复后所有子卷的所有brick都具有相同的 GFID。
-	如前所述，无法使用CLI解决文件系统类型不匹配的问题

#### 3、快速修复
###### 1、获取裂脑中文件的路径：
它可以通过
a） 命令获得。gluster volume heal info split-brain
b） 确定从客户端执行的文件操作不断失败并出现输入/输出错误的文件。
###### 2、关闭从装入点打开此文件的应用程序。 对于虚拟机，需要关闭它们的电源。
###### 3、确定正确的副本
这是通过观察文件的 afr 更改日志扩展属性来完成的 使用 getfattr 命令的砖块;然后确定裂脑的类型 （数据裂脑、元数据裂脑、条目裂脑或裂脑由于 GFID-不匹配）;最后确定哪个砖块包含“好副本” 的文件。
也可能一个砖可能包含正确的数据，而 其他可能包含正确的元数据。getfattr -d -m . -e hex \<file-path-on-brick\>

    0x 000003d7 00000001 00000000
            |      |       |
            |      |        \_ changelog of directory entries
            |       \_ changelog of metadata
            \ _ changelog of data
全为0认为是对的，不为0认为有问题，前8位为data，中间8位位metadata，后8位为directory entries
###### 4、重置包含 使用 setfattr 命令的文件数据/元数据的“错误副本”。
setfattr -n \<attribute-name\> -v \<attribute-value\> \<file-path-on-brick\>
###### 5、通过从客户端执行查找来触发对文件的自我修复：
ls -l \<file-path-on-gluster-mount\>

**example:**
###### 1）查看xattr权限
```shell
[root@pranithk-laptop vol]# getfattr -d -m . -e hex /gfs/brick-?/a
getfattr: Removing leading '/' from absolute path names
\#file: gfs/brick-a/a
trusted.afr.vol-client-0=0x000000000000000000000000     -->/gfs/brick-a/a 上的更新日志认为着某些数据和元数据操作成功了本身
trusted.afr.vol-client-1=0x000003d70000000100000000     -->/gfs/brick-a/a 上的更新日志认为着某些数据和元数据操作在/gfs/brick-a/b 上失败
trusted.gfid=0x80acdbd886524f6fbefa21fc356fed57
\#file: gfs/brick-b/a
trusted.afr.vol-client-0=0x000003b00000000100000000     -->/gfs/brick-b/a 上的更新日志认为着某些数据和元数据操作成功了本身
trusted.afr.vol-client-1=0x000000000000000000000000     -->/gfs/brick-b/a 上的更新日志认为着某些数据和元数据操作在/gfs/brick-a/a 上失败
trusted.gfid=0x80acdbd886524f6fbefa21fc356fed57
```
如果互相指责，就需要人为判断指定一个副本为正确的，并修改其他副本认为正确副本的trusted.afr.vol-client-x<即data/metadata/entries的changelog>为0
如果是3副本，根据xattr判断认为changlelog是对的最多的brick为正确的，修改其他节点认为正确副本的trusted.afr.vol-client-x<即data/metadata/entries的changelog>为0

###### 2）确认正确副本
此处认为gfs/brick-b/a的metadta为正确的，gfs/brick-a/a的data为正确的
###### 3）修改changelog
gfs/brick-a/a:

    修改trusted.afr.vol-client-1的0x000003d70000000100000000为0x000003d70000000000000000 
    /gfs/brick-a/a认为/gfs/brick-b/a的metadata是正确的，将中间8位修改为0
设置gfs/brick-a/a的xattr权限如下:
```shell
setfattr -n trusted.afr.vol-client-1 -v 0x000003d70000000000000000 /gfs/brick-a/a
```
gfs/brick-b/a:

    修改trusted.afr.vol-client-0的0x000003b00000000100000000为0x000000000000000100000000 
    /gfs/brick-b/a认为/gfs/brick-a/a的data是正确的，将前8位修改为0
设置gfs/brick-b/a的xattr权限如下:
```shell
setfattr -n trusted.afr.vol-client-0 -v 0x000000000000000100000000 /gfs/brick-b/a
```
完成上述操作后，查看xattr权限如下所示：
```shell
[root@pranithk-laptop vol]# getfattr -d -m . -e hex /gfs/brick-?/a
getfattr: Removing leading '/' from absolute path names
\#file: gfs/brick-a/a
trusted.afr.vol-client-0=0x000000000000000000000000
trusted.afr.vol-client-1=0x000003d70000000000000000
trusted.gfid=0x80acdbd886524f6fbefa21fc356fed57

\#file: gfs/brick-b/a
trusted.afr.vol-client-0=0x000000000000000100000000
trusted.afr.vol-client-1=0x000000000000000000000000
trusted.gfid=0x80acdbd886524f6fbefa21fc356fed57
```
###### 4）触发自我修复：
client执行以触发愈合。
`ls -l <file-path-on-gluster-mount>`
#### 4、设置自动修复脑裂[ctime|mtime|size|majority]
基于CLI和client的方法需要人工干预，有一个卷设置，当设置为各种可用策略之一时，无需用户干预即可自动恢复脑裂，默认被禁用。
设置以后直接生效，并开始
`cluster.favorite-child-policy`
查看help
```
[root@glusterfs-node2 brick]# gluster v set help  | grep -A3 cluster.favorite-child-policy
Option: cluster.favorite-child-policy
Default Value: none
Description: This option can be used to automatically resolve split-brains using various policies without user intervention.
 "size" picks the file with the biggest size as the source. "ctime" and "mtime" pick the file with the latest ctime and mtime respectively as the source. "majority" picks a file with identical mtime and size in more than half the number of bricks in the replica.
设置glusterfs以最新mtime自动修复脑裂
[root@glusterfs-node2 ~]# gluster v set replica2 cluster.favorite-child-policy mtime
volume set: success
[root@dockernode1 /]# gluster v info
Options Reconfigured:		
cluster.favorite-child-policy: mtime
...
```

#### 5、脑裂处理示例
**环境**
|主机名	                  |Ip地址	        |盘符|
---|:--:|---:
glusterfs-node1	  |  192.168.189.131	|    /dev/sdb
glusterfs-node2 	|192.168.189.132	 |   /dev/sdb
glusterfs-client	|192.168.189.150	

**example1：**
##### 1、手动制造脑裂
###### 1）修改quorum-type和quorum-count
将cluster.quorum-type改为fixed，cluster.quorum-count改为1
[root@glusterfs-node2 ~]# gluster v get replica2 cluster.quorum-type
Option                              Value                                                                     
cluster.quorum-type                     fixed
[root@glusterfs-node2 ~]# gluster v get replica2 cluster.quorum-count
Option          Value                                                                     
cluster.quorum-count                    1
###### 2）停止node1上的所有gluster相关进程
[root@glusterfs-node1 ~]# ps -ef | awk '/gluster/{print $2}' | xargs kill
###### 3）在client端写入文件
[root@gluster-client replica]# mkdir -p stbn/op
[root@gluster-client op]# echo "1900" >>year
[root@gluster-client op]# date >> DateFile
[root@gluster-client replica]# chmod 777 stbn/
###### 4）停止node2上的所有gluster相关进程
[root@glusterfs-node2 ~]# ps -ef | awk '/gluster/{print $2}' | xargs kill
###### 5）启动node1的glusterd服务，并在client端写入文件
[root@glusterfs-node1 ~]# systemctl start glusterd
[root@gluster-client replica]# mkdir -p stbn/op
[root@gluster-client op]# echo "2000" >>year
[root@gluster-client op]# date >> DateFile
[root@glusterfs-node2 ~]# systemctl start glusterd
###### 6）再启动node2的glusterd服务，此时就会脑裂
```shell
[root@glusterfs-node1 ~]# gluster v heal replica2 info
Brick glusterfs-node1:/glusterfs/replica/brick1
/stbn 
/ - Is in split-brain
/stbn/op 
/stbn/op/year 
/stbn/op/date 
Status: Connected
Number of entries: 5

Brick glusterfs-node2:/glusterfs/replica/brick1
/stbn 
/ - Is in split-brain
/stbn/op 
/stbn/op/year 
/stbn/op/date 
Status: Connected
Number of entries: 5

[root@glusterfs-node1 ~]# gluster v heal replica2 info split-brain
Brick glusterfs-node1:/glusterfs/replica/brick1
/
Status: Connected
Number of entries in split-brain: 1

Brick glusterfs-node2:/glusterfs/replica/brick1
/
Status: Connected
Number of entries in split-brain: 1
```
##### 2、开始修复脑裂
/ - Is in split-brain处于脑裂状态，一般目录脑裂为GFID脑裂的对外表现，查看各个brick的GFID是否一致
###### 1）node1上brick的xattr权限
```shell
[root@glusterfs-node1 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1
getfattr: Removing leading '/' from absolute path names
...
trusted.afr.replica2-client-1=0x000000000000000000000001
trusted.gfid=0x00000000000000000000000000000001             --->brick的/目录GFID相同
trusted.glusterfs.mdata=0x0100000000000000000000000065241b230000000003d07ba40000000065241b230000000003d07ba400000000000000000000000000000000
...

[root@glusterfs-node1 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1/stbn
...
trusted.afr.replica2-client-1=0x000000000000000100000002
trusted.gfid=0xfc9cb1976b11485f9b7104bfdd4c3ed7             --->brick的/stbn目录GFID不相同
trusted.glusterfs.mdata=0x0100000000000000000000000065241b230000000004085d520000000065241b230000000004085d520000000065241b230000000003d07ba4
...

[root@glusterfs-node1 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1/stbn/op
...
trusted.afr.replica2-client-1=0x000000000000000100000003
trusted.gfid=0x721ff242da0c49d787cc3c41b1f9c67e             --->brick的/stbn/op目录GFID不相同
trusted.glusterfs.mdata=0x0100000000000000000000000065241b2d000000003983f19f0000000065241b2d000000003983f19f0000000065241b230000000004085d52
...

[root@glusterfs-node1 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1/stbn/op/year
...
trusted.afr.replica2-client-1=0x000000020000000100000000
trusted.gfid=0x59581e2aeb0c42c98776f0dd906832da             --->brick的/stbn/op/year目录GFID不相同
trusted.gfid2path.ddb678ce11c8979b=0x37323166663234322d646130632d343964372d383763632d33633431623166396
33637652f79656172trusted.glusterfs.mdata=0x0100000000000000000000000065241b280000000013aeb43e0000000065241b280000000013
...

[root@glusterfs-node1 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1/stbn/op/date
...
trusted.afr.replica2-client-1=0x000000020000000100000000
trusted.gfid=0x3698bee4f8304ab09f6b5d6754759e6b             --->brick的/stbn/op/date目录GFID不相同
trusted.gfid2path.c17327bebd856a92=0x37323166663234322d646130632d343964372d383763632d33633431623166396
33637652f64617465trusted.glusterfs.mdata=0x0100000000000000000000000065241b2d0000000039afa1970000000065241b2d0000000039
...
```
###### 2）node2上brick的xattr权限
```shell
[root@glusterfs-node2 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1
...
trusted.afr.replica2-client-0=0x000000000000000000000001
trusted.gfid=0x00000000000000000000000000000001
trusted.glusterfs.mdata=0x0100000000000000000000000065241af60000000036efb7b70000000065241af60000000036efb7b700000000000000000000000000000000
...

[root@glusterfs-node2 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1/stbn
...
trusted.afr.replica2-client-0=0x000000000000000200000002
trusted.gfid=0xe2624aa3831c44c8a4e737eb53189795
trusted.glusterfs.mdata=0x0100000000000000000000000065241b1000000000233260110000000065241af6000000003729a2af0000000065241af60000000036efb7b7
...

[root@glusterfs-node2 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1/stbn/op
...
trusted.afr.replica2-client-0=0x000000000000000100000003
trusted.gfid=0x74b1aded1d514ff7acbef5556de1835c
trusted.glusterfs.mdata=0x0100000000000000000000000065241b0a0000000016c5ab720000000065241b0a0000000016c5ab720000000065241af6000000003729a2
...

[root@glusterfs-node2 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1/stbn/op/year
...
trusted.afr.replica2-client-0=0x000000020000000100000000
trusted.gfid=0x9419682c6ccb49c1a7d4b026f2d3633c
trusted.gfid2path.155e95a8e816ed12=0x37346231616465642d316435312d346666372d616362652d66353535366465313
83335632f79656172
trusted.glusterfs.mdata=0x0100000000000000000000000065241b05000000001bb5e47f0000000065241b05000000001bb5e47f0000000065241b05000000001b933aaa
...

[root@glusterfs-node2 ~]# getfattr -d -m. -e hex /glusterfs/replica/brick1/stbn/op/date
...
trusted.afr.replica2-client-0=0x000000020000000100000000
trusted.gfid=0xb75a8e99161849c79f166f5de3b29c37
trusted.gfid2path.93dee3f3fcd904a5=0x37346231616465642d316435312d346666372d616362652d6635353536646531383335632f64617465
trusted.glusterfs.mdata=0x0100000000000000000000000065241b0a000000001723633c0000000065241b0a000000001723633c0000000065241b0a0000000016c5ab72
...
```
###### 3）修复split-brain
当多级目录都发生split-brain时，需要首先修复父目录，如果父目录脑裂时子目录无法修复，修复时的/不是真正的/目录，brick的绝对路径为相对/目录
```
[root@glusterfs-node1 ~]# gluster v heal replica2 split-brain latest-mtime /stbn/op/
Lookup failed on /stbn:Input/output error
Volume heal failed.
[root@glusterfs-node1 ~]# gluster v heal replica2 split-brain latest-mtime /stbn/op
Lookup failed on /stbn:Input/output error
Volume heal failed.
[root@glusterfs-node1 ~]# gluster v heal replica2 split-brain latest-mtime /stbn/op/year
Lookup failed on /stbn/op:Input/output error
Volume heal failed.
[root@glusterfs-node1 ~]# gluster v heal replica2 split-brain latest-mtime /stbn/op/date
Lookup failed on /stbn/op:Input/output error
Volume heal failed.
[root@glusterfs-node1 ~]# gluster v heal replica2 split-brain latest-mtime /stbn/
Lookup failed on /stbn/:Invalid argument.
Volume heal failed.
[root@glusterfs-node1 ~]# gluster v heal replica2 split-brain latest-mtime /
Lookup failed on /:Invalid argument.
Volume heal failed.
------------------------------------------------------------------
[root@glusterfs-node1 ~]# gluster v heal replica2 split-brain latest-mtime /stbn
GFID split-brain resolved for file /stbn
```
###### 4）修复完成查看文件和权限
修复以后，等待一会就会修复完成
```
[root@glusterfs-node2 ~]# gluster v heal replica2 info 
Brick glusterfs-node1:/glusterfs/replica/brick1
Status: Connected
Number of entries: 0

Brick glusterfs-node2:/glusterfs/replica/brick1
Status: Connected
Number of entries: 0
```
查看文件内容和目录权限都是以最新mtime为准
```shell
[root@glusterfs-node1 ~]# cat /glusterfs/replica/brick1/stbn/op/year 
2000
[root@glusterfs-node1 ~]# cat /glusterfs/replica/brick1/stbn/op/date 
Mon Oct  9 23:24:29 CST 2023
[root@glusterfs-node1 ~]# ll -ad /glusterfs/replica/brick1/stbn
drwxr-xr-x 3 root root 16 Oct  9 21:49 /glusterfs/replica/brick1/stbn
```

**example2：**
##### 1、脑裂情况
```shell
[root@glusterfs-node1 ~]# gluster v heal replica2 info
Brick glusterfs-node1:/glusterfs/replica/brick1
/o 
/ - Is in split-brain
/o/p 
/l 
/l/p 
/l/p/date 
/o/p/date 
Status: Connected
Number of entries: 7

Brick glusterfs-node2:/glusterfs/replica/brick1
/l 
/ - Is in split-brain
/l/p 
/o 
/o/p 
/o/p/date 
/l/p/date 
Status: Connected
Number of entries: 7
```
node1
```shell
[root@glusterfs-node1 ~]# cat /glusterfs/replica/brick1/o/p/date 
Sun Oct  8 16:53:58 CST 2023
[root@glusterfs-node1 ~]# cat /glusterfs/replica/brick1/l/p/date 
Sun Oct  8 16:53:55 CST 2023
```
node2
```shell
[root@glusterfs-node2 ~]# cat /glusterfs/replica/brick1/o/p/date 
Sun Oct  8 16:53:17 CST 2023
[root@glusterfs-node2 ~]# cat /glusterfs/replica/brick1/l/p/date 
Sun Oct  8 16:53:21 CST 2023
```
##### 2、开始修复脑裂
###### 1）修复GFID脑裂，以node1为源，修改node2的trusted.afr.\<volume\>-client-x全为0
```
[root@glusterfs-node2 ~]# setfattr -n trusted.afr.replica2-client-0 -v 0x000000000000000000000000 /glusterfs/replica/brick1/
```
###### 2）等待一会就会自动修复
```shell
[root@glusterfs-node2 ~]# gluster v heal replica2 info 
Brick glusterfs-node1:/glusterfs/replica/brick1
Status: Connected
Number of entries: 0

Brick glusterfs-node2:/glusterfs/replica/brick1
Status: Connected
Number of entries: 0
```
node2的文件内容被修改为node1的文件内容
```shell
[root@glusterfs-node2 ~]# cat /glusterfs/replica/brick1/o/p/date 
Sun Oct  8 16:53:58 CST 2023
[root@glusterfs-node2 ~]# cat /glusterfs/replica/brick1/l/p/date 
Sun Oct  8 16:53:55 CST 2023

```