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

http配置

	http {
		include       mime.types;   #导⼊⽀持的⽂件类型
		default_type  application/octet-stream;   #设置默认的类型，会提示下载不匹配的类型⽂件
	#⽇志配置部分
	#log_format main '$remote_addr - $remote_user [$time_local] "$request" '
	# '$status $body_bytes_sent "$http_referer" '
	# '"$http_user_agent" "$http_x_forwarded_for"';
	#access_log logs/access.log main;

	#⾃定义优化参数
		sendfile        on;   #实现⽂件零拷⻉
		#tcp_nopush     on;   #在开启了sendfile的情况下，合并请求后统⼀发送给客户端。
		#tcp_nodelay    off;  #在开启了keepalived模式下的连接是否启⽤TCP_NODELAY选项，当为off时，延迟0.2s发送，默认On时，不延迟发送，⽴即发送⽤户相应报⽂。
		keepalive_timeout  65 65; #设置会话保持时间,如果写俩个参数表示后面的参数会告诉客户端
		gzip on;  #开启⽂件压缩
	
		server {
			listen       80;        #设置监听地址和端⼝
			server_name  localhost;  #设置server name，可以以空格隔开写多个并⽀持正则表达式，如*.baidu.com www.baidu.* 
			location / {
				root   html;
				index  index.html index.htm;
			}
			error_page   500 502 503 504  /50x.html;  #定义错误页面
			location = /50x.html {
				root   html;
			}
	# proxy the PHP scripts to Apache listening on 127.0.0.1:80
	#
	#location ~ \.php$ { #以http的⽅式转发php请求到指定web服务器
	# proxy_pass http://127.0.0.1;
	#}
	# pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
	#
	#location ~ \.php$ { #以fastcgi的⽅式转发php请求到php处理
	# root html;
	# fastcgi_pass 127.0.0.1:9000;
	# fastcgi_index index.php;
	# fastcgi_param SCRIPT_FILENAME /scripts$fastcgi_script_name;
	# include fastcgi_params;
	#}
	# deny access to .htaccess files, if Apache's document root
	# concurs with nginx's one
	#
	#location ~ /\.ht { #拒绝web形式访问指定⽂件，如很多的⽹站都是通过.htaccess⽂件来改变⾃⼰的重定向等功能。
	# deny all;
	#}
		}
	}
	# another virtual host using mix of IP-, name-, and port-based configuration
	#
	#server { #⾃定义虚拟server
	# listen 8000;
	# listen somename:8080;
	# server_name somename alias another.alias;
	# location / {
	# root html;
	# index index.html index.htm; #指定默认⽹⻚⽂件，此指令由ngx_http_index_module模块提供
	# }
	#}
	# HTTPS server
	#
	#server { #https服务器配置
	# listen 443 ssl;
	# server_name localhost;
	# ssl_certificate cert.pem;
	# ssl_certificate_key cert.key;
	# ssl_session_cache shared:SSL:1m;
	# ssl_session_timeout 5m;
	# ssl_ciphers HIGH:!aNULL:!MD5;
	# ssl_prefer_server_ciphers on;
	# location / {
	# root html;
	# index index.html index.htm;
	# }
	location /linux38/passwd.ht {
	deny all;
	}
	#}



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
		location / {										/访问默认域名不加路径,curl www.ky.com
			root /data/nginx/html/pc;						/不写index默认为index index.html
		}
		#location /about {				/curl www.ky.com/about/		
		#	root /data/nginx/html/pc;	/root:指导web的目录，在定义location时，文件绝对路径为root+location;	此处为/data/nginx/html/pc/about/index.html
		#	index index.html;
		#}
		location /about {				
			alias /data/nginx/html/pc;	/定义路径别名,此处路径为alias的路径,文件绝对路径为alias；此处为/data/nginx/html/pc/index.html;		
			index index.html;			
			allow 192.168.189.201;		/四层访问控制，可以通过匹配客户源ip地址进行限制
			deny all;					/拒绝所有，允许源ip地址为192.168.189.201访问
		}
		location = /1.jpg {				/精确匹配，优先级最高,curl www.ky.com/1.jpg，返回/data/nginx/images/1.jpg
			root /data/nginx/images;
			#index index.html;
		}
		location ~ /A.?\.jpg {			/匹配大小写A[].jpg为精确匹配，区分大小写[]内为任意单个字符，比如Ax.jpg，Ax.jpg这个图片必须存在才可以访问，不能匹配到Ax.JPG的图片
			root /data/nginx/images;
		
		}
		#location ~* /A.?\.jpg {		/匹配不区分大小写，可以访问ax.JPG这类资源
		#	root /data/nginx/images;
		
		#}	
		location ^~ /image {			/匹配以image开始，curl www.ky.com/image/，返回/data/nginx/html/image/index.html的页面
			root /data/nginx/html;
			index index.html;
		}
		location /image1 {				/curl www.ky.com/image1/
			root /data/nginx/html;		
			
		}	
		location ~* \.(gif|jpg|bmp|png|js)$ {	/匹配以(gif|jpg|bmp|png|js)结尾，不区分大小写，比如JPG和jpg结尾的都可以匹配到
			root /data/nginx/images;
		}
		error_page 500 502 503 504 404 /error.html;			/定义错误页面
		location = /error.html {							/自定义错误页面，访问不存在的页面curl www.ky.com/about/1.jpg，返回/data/nginx/html/error.html的页面
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


![image](https://raw.githubusercontent.com/DouDou-sudo/linux/main/images/nginx-download%E6%95%88%E6%9E%9C%E5%9B%BE.png)

反向代理
upstream配置参数

	#⾃定义⼀组服务器，配置在http块内
	upstream name {
		server address [parameters];  #配置⼀个后端web服务器，配置在upstream内，⾄少要有⼀个server服务器配置。  
		ip_hash；  #源地址hash调度⽅法，基于的客户端的remote_addr(源地址)做hash计算，以实现会话保持
		least_conn; #最少连接调度算法，优先将客户端请求调度到当前连接最少的后端服务器
	}


	hash KEY consistent； #基于指定key做hash计算，使⽤consistent参数，将使⽤ketama⼀致性hash算法，适⽤于后端是Cache服务器（如varnish）时使⽤，consistent定义使⽤⼀致性hash运算，⼀致性hash基于取模运算。

	hash $request_uri consistent; #基于⽤户请求的uri做hash
nginx的负载均衡支持的4种调度算法

	1）轮询（默认）。每个请求按时间顺序逐一分配到不同的后端服务器，如果后端某台服务器宕机，故障系统被自动剔除，使用户访问不受影响。Weight 指定轮询权值，Weight值越大，分配到的访问机率越高，主要用于后端每个服务器性能不均的情况下。

	2）ip_hash。每个请求按访问IP的hash结果分配，这样来自同一个IP的访客固定访问一个后端服务器，有效解决了动态网页存在的session共享问题。

	3）fair。这是比上面两个更加智能的负载均衡算法。此种算法可以依据页面大小和加载时间长短智能地进行负载均衡，也就是根据后端服务器的响应时间来分配请求，响应时间短的优先分配。Nginx本身是不支持fair的，如果需要使用这种调度算法，必须下载Nginx的upstream_fair模块。

	4）url_hash。此方法按访问url的hash结果来分配请求，使每个url定向到同一个后端服务器，可以进一步提高后端缓存服务器的效率。Nginx本身是不支持url_hash的，如果需要使用这种调度算法，必须安装Nginx 的hash软件包

server支持的参数

	weight=number          #设置权重，默认为1。
	max_conns=number       #给当前server设置最⼤活动链接数，默认为0表示没有限制。
	fail_timeout=time      #对后端服务器的失败监测超时时间，默认为10秒。
	max_fails=number       #在fail_timeout时间对后端服务器连续监测失败多少次就标记为不可⽤。
	proxy_next_upstream=error timeout;  #指定在哪种检测状态下将请求转发器其他服务器。
	backup #预留的备份机器。当其他所有的非backup机器出现故障或者忙的时候，才会请求backup机器，因此这台机器的压力最轻。
	down #标记为down状态。表示当前的server暂时不参与负载均衡
	resolve #当server定义的是主机名的时候，当A记录发⽣变化会⾃动应⽤新IP⽽不⽤重启Nginx。
	
	###注意：当负载调度算法为ip_hash时，后端服务器在负载均衡调度中的状态不能有backup。
示例：
	[root@k8snode2 conf.d]# cat fd.conf
	upstream webserver {
		server 192.168.189.200:80;
		server 192.168.189.203:80;
		server 127.0.0.1:9090 backup; 	/使用代理服务器本身作为backup服务器
	}
	server {
		listen 8080;
		server_name www.bw.cn;
		location / {
			proxy_pass http://webserver;
			#proxy_pass http://192.168.189.200:80/; #反向代理单台web服务器
			proxy_set_header        X-Real-IP       $remote_addr; 
		}
	}

配置backup服务器，使用代理服务器本身作为backup服务器

	[root@k8snode2 conf.d]# cat default.conf 
	server {
		listen 9090;
		server_name localhost;
		location / {
			root /data/nginx/html/error;
			index index.html;
		}
	}
#重启Nginx 并访问测试

	[root@k8snode2 conf.d]# systemctl reload nginx.service
	[root@k8snode1 html]# curl www.bw.cn:8080
	master
	[root@k8snode1 html]# curl www.bw.cn:8080
	node2
	[root@k8snode1 html]# curl www.bw.cn:8080	/当两台后端服务器都停止时
	sorry

客户端ip透传
正常情况下，access_log日志文件保留的是客户端的ip，但是经过代理，客户都ip一致为代理端的ip，所以需要透传
	
	[root@k8snode2 conf.d]# cat fd.conf 
	upstream webserver {
		server 192.168.189.200:80;
		server 192.168.189.203:80;
		
	}
	server {
		listen 8080;
		server_name www.bw.cn;
		location / {
			proxy_pass http://webserver;
			#proxy_pass http://192.168.189.200:80/; #反向代理单台web服务器
			proxy_set_header        X-Real-IP       $remote_addr; #添加客户端ip到报文头部
		}

	}
后端web服务器配置，此处为httpd

	[root@k8smaster logs]# grep "LogFormat" /etc/httpd/conf/httpd.conf 
    #LogFormat "%h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-Agent}i\"" combined  /注释掉这行，将%h修改为%{X-Real-IP}，如下面这行
	LogFormat "%{X-Real-IP}i %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-Agent}i\"" combined
    LogFormat "%h %l %u %t \"%r\" %>s %b" common
	重启httpd
	tail -f /etc/httpd/logs/access_log 
	192.68.189.202 - - [13/Oct/2022:22:28:18 -0400] "GET / HTTP/1.0" 200 6 "-" "curl/7.29.0"   /代理服务器
	192.168.189.201 - - [13/Oct/2022:22:41:48 -0400] "GET / HTTP/1.0" 200 6 "-" "curl/7.29.0"  /修改后为客户都ip
调度算法为ip_hash，每个请求按访问IP的hash结果分配，这样来自同一个IP的访客固定访问一个后端服务器

	[root@k8snode2 conf.d]# cat fd.conf 
	upstream webserver {
		ip_hash; 
		server 192.168.189.200:80;
		server 192.168.189.203:80;
		#server 127.0.0.1:9090 backup;
	}
	server {
		listen 8080;
		server_name www.bw.cn;
		location / {
			proxy_pass http://webserver;
			#proxy_pass http://192.168.189.200:80/; #反向代理单台web服务器
			#proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
			proxy_set_header        X-Real-IP       $remote_addr; 
		}

	}
	重新加载nginx
	客户端访问，一直调度到一台后端web服务器上
	[root@k8snode1 html]# curl www.bw.cn:8080
	master
	[root@k8snode1 html]# curl www.bw.cn:8080
	master
	[root@k8snode1 html]# curl www.bw.cn:8080
	master
	[root@k8snode1 html]# curl www.bw.cn:8080
	master
	[root@k8snode1 html]# curl www.bw.cn:8080
	master
反向代理

	proxy_pass;
	#⽤来设置将客户端请求转发给的后端服务器的主机，可以是主机名、IP地址：端⼝的⽅式，也可以代理到预先设置的主机群组，需要模块ngx_http_upstream_module⽀持。
	location  /web {
		index index.html;
		proxy_pass http://192.168.189.203;
	}
	#http://192.168.189.203  不带斜线将访问的/web,等于访问后端服务器http://192.168.189.203:8080/web/index.html，

	后端web服务器
	[root@centos web]# pwd
	/var/www/html/web
	[root@centos web]# cat index.html
	web 192.168.189.203
	[root@centos html]# pwd
	/var/www/html
	[root@centos html]# cat index.html 
	node2
	测试不带/
	[root@k8snode2 conf.d]# curl www.bw.cn:8080/web/
	web 192.168.189.203
	
	location /web {
		index index.html;
		proxy_pass http://192.168.10.82:80/;
		}
	#带斜线，等于访问后端服务器的http://192.168.7.103:80/index.html 内容返回给客户端
	#测试
	[root@k8snode2 conf.d]# curl www.bw.cn:8080/web/
	node2

