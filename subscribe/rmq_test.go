package subscribe

import (
	"context"
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func Test_rmq_client(t *testing.T) {
	ctx := context.Background()

	s := NewRmq("admin", "liushuojia", "192.168.2.3", 5672, "/")
	if err := s.Start(ctx); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("start")

	_ = s.Register("aaa", func(ctx context.Context, channel string, body []byte) error {
		fmt.Println("aaa", "a001", channel, string(body))
		return nil
	})

	time.Sleep(20 * time.Second)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			log.Println("Publish", "aaa")
			_ = s.Publish(context.Background(), "aaa", []byte(time.Now().Format("20060102 150405")+" pong"))
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			log.Println("UnRegister", "aaa")
			s.UnRegister("aaa")
		}
	}()

	time.Sleep(time.Minute)
	log.Println("Stop")
	_ = s.Stop()

	time.Sleep(time.Second)
	return
}
