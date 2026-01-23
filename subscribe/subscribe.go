package subscribe

import "context"

type Conn interface {
	// Start 启动服务
	Start(ctx context.Context) error

	// Stop 停止服务
	Stop() error

	// Register 注册订阅
	Register(channel string, fn func(context.Context, string, []byte) error) error

	// UnRegister 注销订阅
	UnRegister(channels ...string)

	// Publish 发送消息
	Publish(ctx context.Context, channel string, body []byte) error
	PublishExchange(ctx context.Context, exchange, key string, body []byte) error
}
