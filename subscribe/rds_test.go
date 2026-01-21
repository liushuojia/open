package subscribe

import (
	"context"
	"fmt"
	"testing"
	"time"

	utils "github.com/liushuojia/open"
	log "github.com/sirupsen/logrus"
)

func Test_Client(t *testing.T) {
	rds, err := utils.RedisConnect("192.168.2.3:6379", "liushuojia", 30)
	if err != nil {
		log.Fatalln(err.Error())
	}
	ctx := context.Background()

	s := NewRds(rds)
	if err := s.Start(ctx); err != nil {
		fmt.Println(err)
		return
	}

	_ = s.Register("aaa", "key_add", func(ctx context.Context, channel, msg string) {
		fmt.Println("aaa", channel, msg)
	})
	_ = s.Register("bbb", "key_add", func(ctx context.Context, channel, msg string) {
		fmt.Println("bbb", channel, msg)
	})
	go func() {
		for {
			time.Sleep(5 * time.Second)
			rds.Publish(context.Background(), "ping", time.Now().Format("20060102 150405")+" pp")
			rds.Publish(context.Background(), "get", time.Now().Format("20060102 150405")+" get")
		}

	}()
	go func() {
		time.Sleep(10 * time.Second)
		//cancel()
		_ = s.Register("ping", "key_add", func(ctx context.Context, channel, msg string) {
			fmt.Println("PONG_add", channel, msg)
		})
		_ = s.Register("get", "key", func(ctx context.Context, channel, msg string) {
			fmt.Println("PONG_get", channel, msg)
		})

		time.Sleep(20 * time.Second)
		s.UnRegister("ping", "key_add")
		time.Sleep(10 * time.Second)
		s.UnRegister("get", "key")
	}()

	go func() {
		time.Sleep(8 * time.Second)
		s.Stop()
	}()

	time.Sleep(time.Minute)
	_ = s.Stop()

	time.Sleep(time.Second)

	return
}
