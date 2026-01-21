package subscribe

import "context"

type CallBack func(ctx context.Context, channel, msg string)

type Conn interface {
	// Start 启动服务
	Start(ctx context.Context) error

	// Stop 停止服务
	Stop() error

	// IsRunning 检查服务是否运行中
	IsRunning() bool

	// Register 注册回调
	Register(channel, key string, callBack CallBack) error

	// UnRegister 取消注册回调
	UnRegister(channel, key string)
}
