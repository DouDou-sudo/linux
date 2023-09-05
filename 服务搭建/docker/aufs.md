        "GraphDriver": {
            "Data": {
                "LowerDir": "/var/lib/docker/overlay2/e5313e95e9198732e6e57d2cf8f8f609e8a3959c9e9154b5d54c7169fa70bafc-init/diff:/var/lib/docker/overlay2/12f2b9f25c60f2b4a7999a484de94d68b7b9aec691f70a6310ed2ff349aa4869/diff",                
                "MergedDir": "/var/lib/docker/overlay2/e5313e95e9198732e6e57d2cf8f8f609e8a3959c9e9154b5d54c7169fa70bafc/merged",
                "UpperDir": "/var/lib/docker/overlay2/e5313e95e9198732e6e57d2cf8f8f609e8a3959c9e9154b5d54c7169fa70bafc/diff",
                "WorkDir": "/var/lib/docker/overlay2/e5313e95e9198732e6e57d2cf8f8f609e8a3959c9e9154b5d54c7169fa70bafc/work"
            },
            "Name": "overlay2"
        },

#Lower Dir：image 镜像层(镜像本身，只读)
#Upper Dir：容器的上层(读写)
#Merged Dir：容器的文件系统，使用 Union FS（联合文件系统）将 lowerdir 和 upper Dir：合并给容器使用。
#Work Dir：容器在宿主机的工作目录