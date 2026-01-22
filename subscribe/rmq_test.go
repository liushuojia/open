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
//	time.Sleep(time.Minute)
//	return
//
//	_ = s.Register("aaa", "key_add", func(ctx context.Context, channel string, body []byte) error {
//		fmt.Println("aaa", channel, string(body))
//		return nil
//	})
//	_ = s.Register("bbb", "key_add", func(ctx context.Context, channel string, body []byte) error {
//		fmt.Println("bbb", channel, string(body))
//		return nil
//	})
//
//	time.Sleep(time.Minute)
//	_ = s.Stop()
//
//	time.Sleep(time.Second)
//
//	return
//}
