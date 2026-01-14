package etcd

import "time"

// Option represents the optional function.
type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

type Options struct {
	hosts       []string      // etcd 服务地址	127.0.0.1:2379
	dialTimeout time.Duration // 连接超时时间
	username    string        // 用户名
	password    string        // 密码
}

func WithHost(host string) Option {
	return func(opts *Options) {
		opts.hosts = append(opts.hosts, host)
	}
}
func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(opts *Options) {
		opts.dialTimeout = dialTimeout
	}
}
func WithUsername(username string) Option {
	return func(opts *Options) {
		opts.username = username
	}
}
func WithPassword(password string) Option {
	return func(opts *Options) {
		opts.password = password
	}
}
