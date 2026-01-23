package subscribe

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var _ Conn = (*rds)(nil)

type rds struct {
	ctx    context.Context
	cancel context.CancelFunc

	rds *redis.Client
	sub *redis.PubSub

	isRunning    bool
	lock         sync.Mutex
	subscribeMap sync.Map // sync.map[channel] => CallBack
}

func NewRds(client *redis.Client) Conn {
	return &rds{
		rds:       client,
		isRunning: false,
	}
}

func (s *rds) Start(ctx context.Context) error {
	if s.isRunning {
		return nil
	}
	if s.rds == nil {
		return errors.New("redis is nil")
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	_ = s.Register("ping", func(ctx context.Context, channel string, body []byte) error {
		fmt.Println(channel, "pong", string(body))
		return nil
	})

	go s.subscribe()
	return nil
}
func (s *rds) Stop() error {
	if s.IsRunning() {
		s.cancel()
		_ = s.sub.Close()
	}
	return nil
}
func (s *rds) subscribe() {
	if s.isRunning {
		return
	}

	s.lock.Lock()
	s.isRunning = true
	s.lock.Unlock()

	channelList := make([]string, 0)
	s.subscribeMap.Range(func(key, value any) bool {
		if k, ok := key.(string); ok {
			channelList = append(channelList, k)
		}
		return true
	})
	if len(channelList) <= 0 {
		channelList = append(channelList, "ping")
	}

	log.Println("[redis] subscribe channel:", channelList)
	s.sub = s.rds.Subscribe(s.ctx, channelList...)

	for {
		select {
		case msg := <-s.sub.Channel():
			if vv, ok := s.subscribeMap.Load(msg.Channel); ok {
				if fn, ok := vv.(func(ctx context.Context, channel string, body []byte) error); ok {
					go fn(s.ctx, msg.Channel, []byte(msg.Payload))
				}
			}
		case <-s.ctx.Done():
			goto END
		}
	}

END:
	l := make([]string, 0)
	s.subscribeMap.Range(func(key, value any) bool {
		if k, ok := key.(string); ok {
			l = append(l, k)
		}
		return true
	})
	log.Println("[redis] close subscribe", l)
}
func (s *rds) IsRunning() bool {
	return s.isRunning
}
func (s *rds) Register(channel string, fn func(context.Context, string, []byte) error) error {
	if _, ok := s.subscribeMap.Load(channel); ok {
		s.subscribeMap.Store(channel, fn)
		return nil
	}

	s.subscribeMap.Store(channel, fn)
	if s.isRunning && s.sub != nil {
		log.Println("[redis] subscribe channel:", channel)
		return s.sub.Subscribe(s.ctx, channel)
	}
	return nil
}
func (s *rds) UnRegister(channel string) {
	if _, ok := s.subscribeMap.LoadAndDelete(channel); !ok {
		return
	}
	if s.isRunning && s.sub != nil {
		log.Println("[redis] unsubscribe channel:", channel)
		_ = s.sub.Unsubscribe(s.ctx, channel)
	}
}

func (s *rds) Publish(ctx context.Context, channel string, body []byte) error {
	return s.rds.Publish(ctx, channel, body).Err()
}
func (s *rds) PublishExchange(ctx context.Context, exchange, key string, body []byte) error {
	return s.Publish(ctx, exchange, body)
}
