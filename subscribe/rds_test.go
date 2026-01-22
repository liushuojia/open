package subscribe

import (
	"context"
	"fmt"
	"testing"
	"time"

	utils "github.com/liushuojia/open"
	log "github.com/sirupsen/logrus"
)

func Test_rds_client(t *testing.T) {
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

	_ = s.Register(NewRdsCB("aaa", "key_add", func(ctx context.Context, channel string, body []byte) error {
		fmt.Println("aaa", channel, string(body))
		return nil
	}))
	_ = s.Register(NewRdsCB("bbb", "key_add", func(ctx context.Context, channel string, body []byte) error {
		fmt.Println("bbb", channel, string(body))
		return nil
	}))
	go func() {
		for {
			time.Sleep(5 * time.Second)
			_ = s.Publish(context.Background(), "ping", []byte(time.Now().Format("20060102 150405")+" pp"))
			_ = s.Publish(context.Background(), "get", []byte(time.Now().Format("20060102 150405")+" get"))
		}

	}()
	go func() {
		time.Sleep(10 * time.Second)
		//cancel()
		_ = s.Register(NewRdsCB("ping", "key_add", func(ctx context.Context, channel string, body []byte) error {
			fmt.Println("PONG_add", channel, string(body))
			return nil
		}))
		_ = s.Register(NewRdsCB("get", "key", func(ctx context.Context, channel string, body []byte) error {
			fmt.Println("PONG_get", channel, string(body))
			return nil
		}))

		time.Sleep(20 * time.Second)
		s.UnRegister(NewRdsCB("ping", "key_add", func(ctx context.Context, channel string, body []byte) error {
			fmt.Println("PONG_get", channel, string(body))
			return nil
		}))
		time.Sleep(10 * time.Second)
		s.UnRegister(NewRdsCB("get", "key", func(ctx context.Context, channel string, body []byte) error {
			fmt.Println("PONG_get", channel, string(body))
			return nil
		}))
	}()

	go func() {
		time.Sleep(8 * time.Second)
		s.Stop()
	}()

	time.Sleep(time.Minute)
	_ = s.Stop()

	time.Sleep(10 * time.Second)
}
func Test_rds_client1(t *testing.T) {
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

	_ = s.Register(NewRdsCB("aaa", "a001", func(ctx context.Context, channel string, body []byte) error {
		fmt.Println(channel, "a001", string(body))
		return nil
	}))
	_ = s.Register(NewRdsCB("aaa", "a002", func(ctx context.Context, channel string, body []byte) error {
		fmt.Println(channel, "a002", string(body))
		return nil
	}))
	_ = s.Register(NewRdsCB("aaa", "a002", func(ctx context.Context, channel string, body []byte) error {
		fmt.Println(channel, "a001----00001", string(body))
		return nil
	}))

	go func() {
		for {
			time.Sleep(1 * time.Second)
			_ = s.Publish(context.Background(), "aaa", []byte(time.Now().Format("20060102 150405")+" pp"))
		}
	}()
	go func() {
		time.Sleep(5 * time.Second)
		s.UnRegister(NewRdsCB("aaa", "a001", func(ctx context.Context, channel string, body []byte) error {
			return nil
		}))
		time.Sleep(5 * time.Second)
		s.UnRegister(NewRdsCB("aaa", "a002", func(ctx context.Context, channel string, body []byte) error {
			return nil
		}))
	}()

	time.Sleep(time.Minute)
	_ = s.Stop()

	time.Sleep(10 * time.Second)
}
