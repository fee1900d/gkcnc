- 构建本地镜像 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化
```shell
sudo docker build -t cnc_httpserver .
```
- 将镜像推送至 docker 官方镜像仓库
```shell
sudo docker login
...
sudo docker push zhanng/cnc_httpserver
```
- 通过 docker 命令本地启动 httpserver
```shell
sudo docker run -p 80:80 -d zhanng/cnc_httpserver
```
- 通过 nsenter 进入容器查看 IP 配置
```shell
sudo lsns -t net
sudo nsenter -t 51031 -n ip a
```