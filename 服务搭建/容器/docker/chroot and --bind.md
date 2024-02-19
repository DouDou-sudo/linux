一、绑定挂载
通常在是一个目录上挂载一个设备
使用--bind参数可以将一个目录挂载到另外一个目录上
使用--bind挂载后，无论原目录是否有文件，挂载目录下的原有文件都会被隐藏

    创建挂载目录
    [root@k8snode2 ~]# mkdir /home/chroot
    [root@k8snode2 ~]# mkdir /home/chroot/bin
    [root@k8snode2 ~]# mkdir /home/chroot/lib64
    目录绑定挂载
    [root@k8snode2 ~]# mount --bind /bin/ /home/chroot/bin/
    [root@k8snode2 ~]# mount --bind /lib64 /home/chroot/lib64/
使用findmnt 目录名,列出挂载点的挂载信息，包括绑定挂载

    列出绑定挂载
    [root@k8snode2 ~]# findmnt /home/chroot/bin 
    TARGET           SOURCE                            FSTYPE OPTIONS
    /home/chroot/bin /dev/mapper/centos-root[/usr/bin] xfs    rw,relatime,attr2,inode64,noquota
    [root@k8snode2 ~]# findmnt /home/chroot/lib64 
    TARGET             SOURCE                              FSTYPE OPTIONS
    /home/chroot/lib64 /dev/mapper/centos-root[/usr/lib64] xfs    rw,relatime,attr2,inode64,noquota
    设备挂载
    [root@k8snode2 ~]# findmnt /boot 
    TARGET SOURCE    FSTYPE OPTIONS
    /boot  /dev/sda1 xfs    rw,relatime,attr2,inode64,noquota
卸载绑定挂载

    [root@k8snode2 ~]# umount /home/chroot/lib64
二、chroot
chroot是一个为进程提供有限隔离的程序。使用chroot，我们可以在执行程序时设置根目录。例如，我们可以使用chroot在根目录/home/apache 下运行httpd。这将有效地使/home/apache作为该进程的根目录/。因此，例如，目录/home/apache/www将变为/www。这样，httpd进程将无法访问外部的任何文件/home/apache。

上面将/bin和/lib64/目录挂载在了/home/chroot目录下
使用chroot将根目录更改为/home/chroot目录

    [root@k8snode2 ~]# chroot /home/chroot/
    bash-4.2# pwd
    /
    bash-4.2# ls
    bash: ls: command not found
    bash-4.2# ls -l /

使用ldd命令查看ls所需库文件，ls所需库文件全部在/lib64目录下，上面已经将/lib64目录绑定挂载在/home/chroot目录下

    [root@k8snode2 ~]# ldd /usr/bin/ls
        linux-vdso.so.1 =>  (0x00007ffd271ba000)
        libselinux.so.1 => /lib64/libselinux.so.1 (0x00007f469a2ef000)
        libcap.so.2 => /lib64/libcap.so.2 (0x00007f469a0ea000)
        libacl.so.1 => /lib64/libacl.so.1 (0x00007f4699ee1000)
        libc.so.6 => /lib64/libc.so.6 (0x00007f4699b13000)
        libpcre.so.1 => /lib64/libpcre.so.1 (0x00007f46998b1000)
        libdl.so.2 => /lib64/libdl.so.2 (0x00007f46996ad000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f469a516000)
        libattr.so.1 => /lib64/libattr.so.1 (0x00007f46994a8000)
        libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f469928c000)
查看系统命令ls的的位置

    [root@k8snode2 ~]# which ls
    alias ls='ls --color=auto'
        /usr/bin/ls
    mount --bind /usr/bin /home/chroot/usr/bin/
将/usr/bin目录也绑定在/home/chroot目录下

    [root@k8snode2 ~]# mkdir -p /home/chroot/usr/bin
    [root@k8snode2 ~]# mount --bind /usr/bin /home/chroot/usr/bin/
再次使用chroot将/home/chroot目录更改为/目录

    [root@k8snode2 ~]# chroot /home/chroot/
    bash-4.2# ls
    bin  lib64  usr

