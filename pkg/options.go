package pkg

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
	pathList  []string // 配置文档
	recover   bool     // panic recover
	noWaiting bool     // 阻塞进程等待停止服务
}

func WithConfig(path ...string) Option {
	return func(opts *Options) {
		opts.pathList = path
	}
}
func WithRecover(recover bool) Option {
	return func(opts *Options) {
		opts.recover = recover
	}
}
func WithNoWaiting(noWaiting bool) Option {
	return func(opts *Options) {
		opts.noWaiting = noWaiting
	}
}
