package conf

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
	filePath string
}

func WithFilePath(filePath string) Option {
	return func(opts *Options) {
		opts.filePath = filePath
	}
}
