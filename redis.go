package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const (
	poolSize    = 30
	minIdleConn = 6
)

func RedisConnect(address, password string, db int) (*redis.Client, error) {
	log.Println(fmt.Sprintf(
		"connect redis - address:%s password:%s db:%d",
		address, password, db,
	))

	ops := &redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,

		PoolSize:     poolSize,    // 连接池最大socket连接数，默认为5倍CPU数， 5 * runtime.NumCPU
		MinIdleConns: minIdleConn, // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//超时
		DialTimeout:  5 * time.Second, // 连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, // 读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, // 写超时，默认等于读超时，-1表示取消读超时
		PoolTimeout:  4 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//------------------------------------------------------------------------------------------------------
		// ClusterClient管理着一组redis.Client,下面的参数和非集群模式下的redis.Options参数一致，但默认值有差别。
		// 初始化时，ClusterClient会把下列参数传递给每一个redis.Client
		// 钩子函数
		// 仅当客户端执行命令需要从连接池获取连接时，如果连接池需要新建连接则会调用此钩子函数
		OnConnect: func(ctx context.Context, conn *redis.Conn) error {
			//log.Println("redis", "connect", conn)
			return nil
		},
	}

	redisClient := redis.NewClient(ops)
	err := redisClient.Ping(context.Background()).Err()
	return redisClient, err
}
