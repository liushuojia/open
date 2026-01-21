package subscribe

import (
	"context"
	"errors"
	"sync"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var _ Conn = (*rds)(nil)

type rds struct {
	ctx          context.Context
	cancel       context.CancelFunc
	rds          *redis.Client
	sub          *redis.PubSub
	isRunning    bool
	lock         sync.Mutex
	subscribeMap map[string]map[string]CallBack
}

func NewRds(client *redis.Client) Conn {
	return &rds{
		rds:          client,
		isRunning:    false,
		subscribeMap: make(map[string]map[string]CallBack),
	}
}

func (s *rds) Start(ctx context.Context) error {
	if s.rds == nil {
		return errors.New("redis is nil")
	}

	if s.isRunning {
		return nil
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	_ = s.Register("ping", "key", func(ctx context.Context, channel, msg string) {
		log.Println("PONG", channel, msg)
	})

	go s.subscribe()
	return nil
}
func (s *rds) Stop() error {
	if s.IsRunning() {
		s.cancel()
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
	for channel := range s.subscribeMap {
		channelList = append(channelList, channel)
	}

	log.Println("[redis] subscribe channel:", channelList)
	s.sub = s.rds.Subscribe(s.ctx, channelList...)

	var (
		ch = s.sub.Channel()
	)
	defer func() {
		l := make([]string, 0)
		for channel := range s.subscribeMap {
			l = append(l, channel)
		}
		log.Println("[redis] close subscribe", l)
		_ = s.sub.Close()
	}()
	for {
		select {
		case msg := <-ch:
			if l, ok := s.subscribeMap[msg.Channel]; ok {
				for _, f := range l {
					go f(s.ctx, msg.Channel, msg.Payload)
				}
			}
		case <-s.ctx.Done():
			s.lock.Lock()
			s.isRunning = false
			s.lock.Unlock()
			return
		}
	}
}
func (s *rds) IsRunning() bool {
	return s.isRunning
}
func (s *rds) Register(channel, key string, callback CallBack) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.subscribeMap[channel]; ok {
		s.subscribeMap[channel][key] = callback
		return nil
	}

	s.subscribeMap[channel] = make(map[string]CallBack)
	s.subscribeMap[channel][key] = callback
	if s.isRunning {
		log.Println("[redis] subscribe channel:", channel)
		return s.sub.Subscribe(s.ctx, channel)
	}
	return nil
}
func (s *rds) UnRegister(channel, key string) {
	if _, ok := s.subscribeMap[channel]; !ok {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.subscribeMap[channel], key)
	if len(s.subscribeMap[channel]) <= 0 && s.isRunning {
		delete(s.subscribeMap, channel)
		log.Println("[redis] unsubscribe channel:", channel)
		_ = s.sub.Unsubscribe(s.ctx, channel)
	}
}
