server配置
	
	[root@k8snode2 conf.d]# cat pc.conf 
	server {												/自定义虚拟server			
		listen 80;                                          /设置监听80端口
		server_name www.ky.com;								/设置server name，
		access_log /data/nginx/logs/www.ky.com_access.log;	/自定义日志地址
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
		location = /error.html {				/自定义错误页面，
			root /data/nginx/html;
		}

	}
