### nsq-docker docker 集群部署
需要注意 在host 里面配虚拟地址
```shell
由于是docker环境,要注意ip地址是否映射到外网,所以这里还是老老实实写127.0.0.1
127.0.0.1 nsqlookupd
127.0.0.1 nsqd1
127.0.0.1 nsqd2
127.0.0.1 nsqd3
```
