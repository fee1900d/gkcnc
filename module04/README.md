- 4.1 Kubeadm 安装 Kubernetes 集群

    https://github.com/cncamp/101/blob/master/k8s-install/readme.md

- 4.2 Envoy Deployment 实验

  https://github.com/cncamp/101/blob/master/module4/envoy-deploy.yaml

注意
---
- 安装时不要指定版本号
- 使用GCP等外网云服务，不要指定国内镜像
- 对于公有云，`kubeadm init` 指定的网络选项应该是云服务的内网 (`--pod-network-cidr` , `--apiserver-advertise-address`等 )；
  - 如果指定的网络为外网会初始化失败，通过查看日志出现拒绝连接情况
    ```shell
    $ sudo crictl --runtime-endpoint unix:///var/run/containerd/containerd.sock logs 361c3b6af1aad
    ......
    W0619 03:31:41.230019       1 clientconn.go:1331] [core] grpc: addrConn.createTransport failed to connect to {127.0.0.1:2379 127.0.0.1 <nil> 0 <nil>}. Err: connection error: desc = "transport: Error while dialing dial tcp 127.0.0.1:2379: connect: connection refused". Reconnecting...
    ......
    ```
