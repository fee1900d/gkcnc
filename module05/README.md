## 5.1 本地创建单节点基于 HTTPS 集群，增查删
- https://github.com/cncamp/101/blob/master/module5/etcd-binary-setup.MD

### 创建工作目录

```sh
mkdir -p /data/k8s-work
cd /data/k8s-work
```

### 下载 cfssl 工具

```sh
wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
wget https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64

chmod +x cfssl*
mv cfssl_linux-amd64 /usr/local/bin/cfssl
mv cfssljson_linux-amd64 /usr/local/bin/cfssljson
mv cfssl-certinfo_linux-amd64 /usr/local/bin/cfssl-certinfo
```

### 生成 ca 配置文件

```sh
cat > ca-csr.json <<"EOF"
{
"CN": "kubernetes",
"key": {
"algo": "rsa",
"size": 2048
},
"names": [
{
"C": "CN",
"ST": "Shanghai",
"L": "Shanghai",
"O": "cncamp",
"OU": "cncamp"
}
],
"ca": {
"expiry": "87600h"
}
}
EOF
```

### 生成 ca 证书文件

```sh
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

### 配置 ca 证书策略

```sh
cat > ca-config.json <<"EOF"
{
  "signing": {
      "default": {
          "expiry": "87600h"
        },
      "profiles": {
          "kubernetes": {
              "usages": [
                  "signing",
                  "key encipherment",
                  "server auth",
                  "client auth"
              ],
              "expiry": "87600h"
          }
      }
  }
}
EOF
```

### 配置 etcd 请求 csr 文件
> 注意修改IP
```sh
cat > etcd-csr.json <<"EOF"
{
  "CN": "etcd",
  "hosts": [
    "127.0.0.1",
    "192.168.34.2"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [{
    "C": "CN",
    "ST": "Shanghai",
    "L": "Shanghai",
    "O": "cncamp",
    "OU": "cncamp"
  }]
}
EOF
```

### 生成证书

```sh
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kubernetes etcd-csr.json | cfssljson  -bare etcd
```

```sh
wget https://github.com/etcd-io/etcd/releases/download/v3.5.0/etcd-v3.5.0-linux-amd64.tar.gz
tar -xvf etcd-v3.5.0-linux-amd64.tar.gz
cp -p etcd-v3.5.0-linux-amd64/etcd* /usr/local/bin/
```

### 生成 etcd 配置文件
> 注意端口，不要与k8s服务冲突
```sh
cat >  etcd.conf <<"EOF"
#[Member]
ETCD_NAME="etcd1"
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="https://192.168.34.2:12380"
ETCD_LISTEN_CLIENT_URLS="https://192.168.34.2:12379,http://127.0.0.1:12379"

#[Clustering]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://192.168.34.2:12380"
ETCD_ADVERTISE_CLIENT_URLS="https://192.168.34.2:12379"
ETCD_INITIAL_CLUSTER="etcd1=https://192.168.34.2:12380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
EOF
```

### 创建启动service
```shell
cat > etcd.service <<"EOF"
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
EnvironmentFile=/etc/etcd/etcd.conf
WorkingDirectory=/var/lib/etcd/
ExecStart=/usr/local/bin/etcd \
  --cert-file=/etc/etcd/ssl/etcd.pem \
  --key-file=/etc/etcd/ssl/etcd-key.pem \
  --trusted-ca-file=/etc/etcd/ssl/ca.pem \
  --peer-cert-file=/etc/etcd/ssl/etcd.pem \
  --peer-key-file=/etc/etcd/ssl/etcd-key.pem \
  --peer-trusted-ca-file=/etc/etcd/ssl/ca.pem \
  --peer-client-cert-auth \
  --client-cert-auth
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
```

### etcd 目录
```shell
mkdir -p /etc/etcd
mkdir -p /etc/etcd/ssl
mkdir -p /var/lib/etcd/default.etcd
```

### 分发证书和配置
```shell
cp ca*.pem /etc/etcd/ssl/
cp etcd*.pem /etc/etcd/ssl/
cp etcd.conf /etc/etcd/
cp etcd.service /usr/lib/systemd/system/
```

### 启动集群
```shell
systemctl daemon-reload
systemctl enable --now etcd.service
systemctl status etcd
```
> 如果启动失败，提示 Failed to start etcd.service: Unit is masked.
> 则手动执行 `systemctl unmask etcd.service` 然后再尝试启动

### 查看集群状态
> 注意端口
```shell
ETCDCTL_API=3 etcdctl --write-out=table --cacert=/etc/etcd/ssl/ca.pem --cert=/etc/etcd/ssl/etcd.pem --key=/etc/etcd/ssl/etcd-key.pem --endpoints=https://192.168.128.200:12379 endpoint health
```

## 5.2 k8s 中创建高可用 ectd 集群
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/setup-ha-etcd-with-kubeadm/