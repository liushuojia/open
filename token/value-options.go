package token

type ValueOption func(opts *ValueOptions)

func loadValueOptions(options ...ValueOption) *ValueOptions {
	opts := new(ValueOptions)
	for _, option := range options {
		option(opts)
	}
	return opts
}

type ValueOptions struct {
	ID    uint64
	AppID string
	Key   string
	UUID  string
}

func WithValueOptions(options ValueOptions) ValueOption {
	return func(opts *ValueOptions) {
		*opts = options
	}
}
func WithValueID(id uint64) ValueOption {
	return func(opts *ValueOptions) {
		opts.ID = id
	}
}
func WithValueAppID(appID string) ValueOption {
	return func(opts *ValueOptions) {
		opts.AppID = appID
	}
}
func WithValueKey(key string) ValueOption {
	return func(opts *ValueOptions) {
		opts.Key = key
	}
}
func WithValueUUID(uuid string) ValueOption {
	return func(opts *ValueOptions) {
		opts.UUID = uuid
	}
}
