package subscribe

//import (
//	"context"
//	"fmt"
//	"testing"
//	"time"
//)
//
//func Test_rmq_client(t *testing.T) {
//	ctx := context.Background()
//
//	s := NewRmq("admin", "liushuojia", "192.168.2.3", 5672, "/")
//	if err := s.Start(ctx); err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("start")
//
//	_ = s.Register(NewRmqCB("aaa", "a001", func(ctx context.Context, channel string, body []byte) error {
//		fmt.Println("aaa", "a001", channel, string(body))
//		return nil
//	}))
//	_ = s.Register(NewRmqCB("aaa", "a002", func(ctx context.Context, channel string, body []byte) error {
//		fmt.Println("aaa", "a002", channel, string(body))
//		return nil
//	}))
//
//	time.Sleep(6 * time.Second)
//
//	go func() {
//		for {
//			time.Sleep(1 * time.Second)
//			_ = s.Publish(context.Background(), "aaa", []byte(time.Now().Format("20060102 150405")+" pong"))
//		}
//	}()
//
//	go func() {
//		for {
//			time.Sleep(5 * time.Second)
//			s.UnRegister(NewRmqCB("aaa", "a001", func(ctx context.Context, s string, bytes []byte) error {
//				return nil
//			}))
//			time.Sleep(5 * time.Second)
//			s.UnRegister(NewRmqCB("aaa", "a002", func(ctx context.Context, s string, bytes []byte) error {
//				return nil
//			}))
//		}
//	}()
//
//	time.Sleep(time.Minute)
//	_ = s.Stop()
//
//	time.Sleep(time.Second)
//
//	return
//}
