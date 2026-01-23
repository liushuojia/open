package subscribe

import "context"

var _ CallBack = (*callBack)(nil)

type CallBack interface {
	FN(context.Context, string, []byte) error // 回调函数
	Channel() string                          // 订阅 channel / queue
}

type callBack struct {
	channel, name string
	fn            func(context.Context, string, []byte) error
}

func NewCallBack(channel string, fn func(context.Context, string, []byte) error) CallBack {
	return &callBack{
		channel: channel,
		fn:      fn,
	}
}

func (cb *callBack) FN(ctx context.Context, channel string, body []byte) error {
	return cb.fn(ctx, channel, body)
}
func (cb *callBack) Channel() string {
	return cb.channel
}
