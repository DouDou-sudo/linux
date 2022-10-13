全局配置

	[16:08:57 root@centos8 ~]#grep -v "#" /apps/nginx/conf/nginx.conf | grep -v "^$"
	user nginx [nginx]; #启动nginx工作进行的用户和组
	worker_processes  [number|auto]; #启动nginx工作进程的数量,可以写个数，也可以写auto自动，系统有几个cpu内核自动开启几个工作进程
	worker_cpu_affinity 0001 0010 0100 1000; #将Nginx⼯作进程绑定到指定的CPU核⼼，默认Nginx是不进⾏进程绑定的，绑定并不是意味着当前nginx进程独占以⼀核⼼CPU，但是可以保证此进程不会运⾏在其他核⼼上，这就极⼤减少了nginx的⼯作进程在不同的cpu核⼼上的来回跳转，减少了CPU对进程的资源分配与回收以及内存管理等，因此可以有效的提升nginx服务器的性能。可以写auto自动绑定

	#错误⽇志记录配置，语法：error_log file [debug | info | notice | warn | error |crit | alert | emerg]可以根据错误级别分别设置
	#error_log logs/error.log;
	#error_log logs/error.log notice;
	error_log /apps/nginx/logs/error.log error;

	#pid⽂件保存路径
	pid /apps/nginx/logs/nginx.pid;
	worker_priority 0;  #⼯作进程nice值即进程运行优先级越小越优先，-20~19
	worker_rlimit_nofile 65536; #这个数字包括Nginx的所有连接（例如与代理服务器的连接等），⽽不仅仅是与客户端的连接,另⼀个考虑因素是实际的并发连接数不能超过系统级别的最⼤打开⽂件数的限制.
	daemon off;  #前台运⾏Nginx服务⽤于测试、docker等环境。
	master_process off|on;  #是否开启Nginx的master-woker⼯作模式，仅⽤于开发调试场景。

	events {   #事件模型配置参数
		worker_connections  1024; #设置单个⼯作进程的最⼤并发连接数
		use epoll;    #使⽤epoll事件驱动，Nginx⽀持众多的事件驱动，⽐如select、poll、epoll，只能设置在events模块中设置。
		accept_mutex on; #优化同⼀时刻只有⼀个请求⽽避免多个睡眠进程被唤醒的设置，on为防⽌被同时唤醒，默认为off，全部唤醒的过程也成为"惊群"，因此nginx刚安装完以后要进⾏适当的优化。
		multi_accept on; #Nginx服务器的每个⼯作进程可以同时接受多个新的⽹络连接，但是需要在配置⽂件中配置，此指令默认为关闭，即默认为⼀个⼯作进程只能⼀次接受⼀个新的⽹络连接，打开后⼏个同时接受多个。
	}




server配置
location的详细使用

	语法规则： location [=|~|~*|^~] /uri/ { … }

	= #⽤于标准uri前，需要请求字串与uri精确匹配，如果匹配成功就停⽌向下匹配并⽴即处理请求。
	~ #⽤于标准uri前，表示包含正则表达式并且区分⼤⼩写，并且匹配
	!~ #⽤于标准uri前，表示包含正则表达式并且区分⼤⼩写，并且不匹配


	~* #⽤于标准uri前，表示包含正则表达式并且不区分⼤写，并且匹配
	!~* #⽤于标准uri前，表示包含正则表达式并且不区分⼤⼩写,并且不匹配

	^~ #⽤于标准uri前，表示包含正则表达式并且匹配以什么开头
	$ #⽤于标准uri前，表示包含正则表达式并且匹配以什么结尾
	\ #⽤于标准uri前，表示包含正则表达式并且转义字符。可以转. * ?等
	* #⽤于标准uri前，表示包含正则表达式并且代表任意⻓度的任意字符

	匹配优先级 =,^~,~/~*,/
	
	[root@k8snode2 conf.d]# cat pc.conf 
	server {												/自定义虚拟server			
		listen 80;                                          /设置监听80端口
		server_name www.ky.com;								/设置server name，
		access_log /data/nginx/logs/www.ky.com_access.log;	/自定义访问日志地址
		error_log /data/nginx/logs/www,ky,com_error.log;
		location / {										/访问默认域名不加路径
			root /data/nginx/html/pc;						/不写index默认为index index.html
		}
		#location /about {									
		#	root /data/nginx/html/pc;	/root:指导web的目录，在定义location时，文件绝对路径为root+location;	此处为/data/nginx/html/pc/about/index.html
		#	index index.html;
		#}
		location /about {				
			alias /data/nginx/html/pc;	/定义路径别名,此处路径为alias的路径,文件绝对路径为alias；此处为/data/nginx/html/pc/index.html;		
			index index.html;			
			allow 192.168.189.201;		/四层访问控制，可以通过匹配客户源ip地址进行限制
			deny all;					/拒绝所有，允许源ip地址为192.168.189.201访问
		}
		location = /1.jpg {				/精确匹配，优先级最高
			root /data/nginx/images;
			#index index.html;
		}
		location ~ /A.?\.jpg {			/匹配大小写A[].jpg为精确匹配，区分大小写[]内为任意单个字符，比如Ax.jpg，Ax.jpg这个图片必须存在才可以访问，不能匹配到Ax.JPG的图片
			root /data/nginx/images;
		
		}
		#location ~* /A.?\.jpg {		/匹配不区分大小写，可以访问ax.JPG这类资源
		#	root /data/nginx/images;
		
		#}	
		location ^~ /image {			/匹配以image开始，www.ky.com/image/
			root /data/nginx/html;
			index index.html;
		}
		location /image1 {
			root /data/nginx/html;
			
		}	
		location ~* \.(gif|jpg|bmp|png|js)$ {	/匹配以(gif|jpg|bmp|png|js)结尾，不区分大小写，比如JPG和jpg结尾的都可以匹配到
			root /data/nginx/images;
		}
		error_page 500 502 503 504 404 /error.html;			/定义错误页面
		location = /error.html {							/自定义错误页面，
			root /data/nginx/html;
		}
		location /download {					/作为下载服务器配置
			autoindex on;						/自动索引功能
			autoindex_exact_size on;			/计算⽂件确切⼤⼩（单位bytes），off只显示⼤概⼤⼩（单位kb、mb、gb）
			autoindex_localtime on;				/显示本机时间⽽⾮GMT(格林威治)时间
			root /data/nginx/html/pc;			/下载文件本机地址/data/nginx/html/pc/download
			limit_rate 10k;						/限制响应给客户端的传输速率，单位是bytes/second，默认值为0表示不限制
		}
	}
查看路径下的文件

	[root@k8snode2 ~]# tree /data/nginx/html/pc/download/
	/data/nginx/html/pc/download/
	├── 123
	└── abc
本地访问
![图片](https://user-images.githubusercontent.com/62938705/195536340-c3b61e81-ada6-4ec5-b7c2-417547b899cb.png)
