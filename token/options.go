package token

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
	Subject string
	Issuer  string
	Expire  string
	Secret  []byte
}

func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}
func WithIssuer(issuer string) Option {
	return func(opts *Options) {
		opts.Issuer = issuer
	}
}
func WithExpire(expire string) Option {
	return func(opts *Options) {
		opts.Expire = expire
	}
}
func WithSecret(secret []byte) Option {
	return func(opts *Options) {
		opts.Secret = secret
	}
}
