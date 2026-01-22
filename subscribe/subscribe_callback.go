package subscribe

import "context"

var _ CallBack = (*callBack)(nil)

type CallBack interface {
	FN(context.Context, string, []byte) error // 回调函数
	Channel() string                          // 订阅 channel / exchange
	Key() string                              // 订阅 key
	Name() string                             // 注册 名称 唯一索引
}

type callBack struct {
	channel, key, name string
	fn                 func(context.Context, string, []byte) error
}

func (cb *callBack) FN(ctx context.Context, channel string, body []byte) error {
	return cb.fn(ctx, channel, body)
}
func (cb *callBack) Channel() string {
	return cb.channel
}
func (cb *callBack) Key() string {
	return cb.key
}
func (cb *callBack) Name() string {
	return cb.name
}
