package mail

type (
	mail struct {
		name    string
		address string
	}
	mailAttach struct {
		name string
		path string
	}
)

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
	subject    string       // 主题
	body       string       // 内容
	from       mail         // 发送人
	mailTo     []mail       // 收件人
	mailCc     []mail       // 抄送
	mailBcc    []mail       // 暗送
	mailAttach []mailAttach // 附件
}

func WithSubject(subject string) Option {
	return func(opts *Options) {
		opts.subject = subject
	}
}
func WithBody(body string) Option {
	return func(opts *Options) {
		opts.body = body
	}
}
func WithFrom(name, account string) Option {
	return func(opts *Options) {
		opts.from = mail{
			name:    name,
			address: account,
		}
	}
}
func WithMailTo(name, address string) Option {
	return func(opts *Options) {
		opts.mailTo = append(opts.mailTo, mail{name: name, address: address})
	}
}
func WithMailCc(name, address string) Option {
	return func(opts *Options) {
		opts.mailCc = append(opts.mailCc, mail{name: name, address: address})
	}
}
func WithMailBcc(name, address string) Option {
	return func(opts *Options) {
		opts.mailBcc = append(opts.mailBcc, mail{name: name, address: address})
	}
}
func WithMailAttach(name, path string) Option {
	return func(opts *Options) {
		opts.mailAttach = append(opts.mailAttach, mailAttach{name: name, path: path})
	}
}
