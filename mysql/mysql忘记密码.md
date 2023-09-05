登录mysql

    [root@centos ~]# mysql -uroot -p
    Enter password: 
    ERROR 1045 (28000): Access denied for user 'root'@'localhost' (using password: YES)
权限被拒绝，忘记密码
修复配置文件，添加`skip-grant-tables`行

    [root@centos ~]# vim /etc/my.cnf
    [mysqld]
    skip-grant-tables
跳过认证直接登录

    [root@centos ~]# mysql -uroot
    mysql> use mysql;
    mysql> alter user 'root'@'localhost' IDENTIFIED BY 'qwer1234';
    ERROR 1290 (HY000): The MySQL server is running with the --skip-grant-tables option so it cannot execute this statement
使用--skip-grant-tables跳过认证，无法执行此语句修改密码。刷新权限再执行即可

    mysql> flush privileges;
    Query OK, 0 rows affected (0.00 sec)
    mysql> alter user 'root'@'localhost' identified by 'qwer1234';
    ERROR 1819 (HY000): Your password does not satisfy the current policy requirements
密码不符合当前策略要求。mysql8.0以上密码强度默认要求较高，更改密码强度
查看密码强度

    mysql> show variables like 'validate_password%';
    +--------------------------------------+--------+
    | Variable_name                        | Value  |
    +--------------------------------------+--------+
    | validate_password.check_user_name    | ON     |
    | validate_password.dictionary_file    |        |
    | validate_password.length             | 8      |
    | validate_password.mixed_case_count   | 1      |
    | validate_password.number_count       | 1      |
    | validate_password.policy             | MEDIUM |
    | validate_password.special_char_count | 1      |
    +--------------------------------------+--------+
    7 rows in set (0.00 sec)
更改密码策略为low

    mysql> set global validate_password.policy=LOW;
    Query OK, 0 rows affected (0.00 sec)
更改密码长度为6

    mysql> set global validate_password.length=6;
    Query OK, 0 rows affected (0.00 sec)
再修改密码

    mysql> alter user 'root'@'localhost' identified by 'qwer1234';
    ERROR 1396 (HY000): Operation ALTER USER failed for 'root'@'localhost'
alter user失败

    mysql> select host,user from user;
    +-----------+------------------+
    | host      | user             |
    +-----------+------------------+
    | %         | root             |
    | localhost | mysql.infoschema |
    | localhost | mysql.session    |
    | localhost | mysql.sys        |
    +-----------+------------------+
    4 rows in set (0.00 sec)
host列的值为%，更改localhost为%

    mysql> alter user 'root'@'%' identified by 'qwer1234';
    Query OK, 0 rows affected (0.01 sec)

    mysql> flush privileges;
    Query OK, 0 rows affected (0.00 sec)

    mysql> quit
    Bye
恢复配置文件

    [root@centos ~]# vim /etc/my.cnf
    #skip-grant-tables
重新登录

    [root@centos ~]# mysql -uroot -pqwer1234
    mysql: [Warning] Using a password on the command line interface can be insecure.
    Welcome to the MySQL monitor.  Commands end with ; or \g.

