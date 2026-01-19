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
	pathList []string // 配置文档
}

func WithConfig(path ...string) Option {
	return func(opts *Options) {
		opts.pathList = path
	}
}
