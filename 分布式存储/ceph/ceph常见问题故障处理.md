### 1、application not enabled on 1 pool(s)
[root@ceph-node2 osd]# ceph -s
  cluster:
    id:     ebc64023-50ef-4418-bd5d-ad16db89a284
    health: HEALTH_WARN
            application not enabled on 1 pool(s)

解决方法：
根据提示把pool设置一下
[root@ceph-node2 osd]# ceph health detail
HEALTH_WARN application not enabled on 1 pool(s)
POOL_APP_NOT_ENABLED application not enabled on 1 pool(s)
    application not enabled on pool 'genesha_data'
    use 'ceph osd pool application enable <pool-name> <app-name>',
    where <app-name> is 'cephfs','rbd','rgw' or freeform for custom applications

[root@ceph-node2 osd]# ceph osd pool application enable genesha_data rgw
enabled application 'rgw' on pool 'genesha_data'
### 2、journal read_header error decoding journal header
ceph 添加osd报错journal read_header error decoding journal header
解决方法：
1、重新格式化数据盘和缓存盘
mkfs.xfs /dev/sdb2 -f
mkfs.xfs /dev/sdb3 -f
2、清理/var/lib/ceph/osd/ceph-0/目录下无关文件，只留下fsid和journal文件
3、dd覆盖一下缓存篇的第一个扇区
dd if=/dev/zero of=/dev/sdb3 bs=512 count=1
4、重新带key执行,osd对应的key在/etc/ceph/目录下
ceph-osd -i 0 --mkfs --keyring /etc/ceph/ceph.osd.0.keyring