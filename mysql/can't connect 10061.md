[root@centos ~]# mysql -h192.168.189.100 -uroot -pqwer1234
mysql: [Warning] Using a password on the command line interface can be insecure.
ERROR 2003 (HY000): Can't connect to MySQL server on '192.168.189.100' (111)

Navicat 连接数据库报错
Can't connect to MySQL server on 'localhost' (10061)
问题原因: 刚改了配置文件添加了一行又注释掉#skip-grant-tables，没有重启mysqld
解决方法: 重启下mysqld再连接就可以解决