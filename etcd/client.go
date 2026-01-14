package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

func New(options ...Option) (*clientv3.Client, error) {
	opt := loadOptions(options...)

	// 1. 配置 etcd 客户端连接参数
	config := clientv3.Config{
		Endpoints:   opt.hosts,       // etcd 服务地址
		DialTimeout: opt.dialTimeout, // 连接超时时间

		// 生产环境需配置认证：Username/Password/TLS
		Username: opt.username,
		Password: opt.password,
	}

	// 2. 建立连接
	return clientv3.New(config)
}
