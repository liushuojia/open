package subscribe

import "context"

type Conn interface {
	// Start 启动服务
	Start(ctx context.Context) error

	// Stop 停止服务
	Stop() error

	// IsRunning 检查服务是否运行中
	IsRunning() bool

	// Register 注册订阅
	Register(callBack CallBack) error

	// UnRegister 注销订阅
	UnRegister(callBack CallBack)

	// Publish 发送消息
	Publish(ctx context.Context, channel string, body []byte) error
	publishExchange(ctx context.Context, exchange, key string, body []byte) error
}
