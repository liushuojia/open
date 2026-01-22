package subscribe

import (
	"context"
	"errors"
	"sync"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var _ Conn = (*rds)(nil)

func NewRdsCB(channel, name string, fn func(context.Context, string, []byte) error) CallBack {
	return &callBack{
		channel: channel,
		name:    name,
		fn:      fn,
	}
}

type rds struct {
	ctx    context.Context
	cancel context.CancelFunc

	rds *redis.Client
	sub *redis.PubSub

	isRunning    bool
	lock         sync.Mutex
	subscribeMap sync.Map // sync.map[channel] => []CallBack
}

func NewRds(client *redis.Client) Conn {
	return &rds{
		rds:       client,
		isRunning: false,
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
	_ = s.Register(NewRdsCB("ping", "key", func(ctx context.Context, channel string, body []byte) error {
		log.Println("PONG", channel, string(body))
		return nil
	}))

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

	log.Println("[redis] subscribe channel:", channelList)
	s.sub = s.rds.Subscribe(s.ctx, channelList...)

	ch := s.sub.Channel()
	defer func() {
		l := make([]string, 0)
		s.subscribeMap.Range(func(key, value any) bool {
			if k, ok := key.(string); ok {
				l = append(l, k)
			}
			return true
		})
		log.Println("[redis] close subscribe", l)
		_ = s.sub.Close()
	}()

	for {
		select {
		case msg := <-ch:
			if vv, ok := s.subscribeMap.Load(msg.Channel); ok {
				if l, ok := vv.([]CallBack); ok {
					for _, cb := range l {
						go cb.FN(s.ctx, msg.Channel, []byte(msg.Payload))
					}
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
func (s *rds) Register(cb CallBack) error {
	if lv, ok := s.subscribeMap.Load(cb.Channel()); ok {
		if l, ok := lv.([]CallBack); ok {
			cbList := make([]CallBack, 0)
			for _, v := range l {
				if v.Name() == cb.Name() {
					continue
				}
				cbList = append(cbList, v)
			}
			cbList = append(cbList, cb)
			s.subscribeMap.Store(cb.Channel(), cbList)
			return nil
		}
	}

	s.subscribeMap.Store(cb.Channel(), []CallBack{
		cb,
	})
	if s.isRunning && s.sub != nil {
		log.Println("[redis] subscribe channel:", cb.Channel())
		return s.sub.Subscribe(s.ctx, cb.Channel())
	}
	return nil
}
func (s *rds) UnRegister(cb CallBack) {
	lv, ok := s.subscribeMap.Load(cb.Channel())
	if !ok {
		return
	}

	l, ok := lv.([]CallBack)
	if !ok {
		return
	}

	cbList := make([]CallBack, 0)
	for _, v := range l {
		if v.Name() == cb.Name() {
			continue
		}
		cbList = append(cbList, v)
	}
	if len(cbList) > 0 {
		s.subscribeMap.Store(cb.Channel(), cbList)
		return
	}

	s.subscribeMap.Delete(cb.Channel())
	if s.isRunning && s.sub != nil {
		log.Println("[redis] unsubscribe channel:", cb.Channel())
		_ = s.sub.Unsubscribe(s.ctx, cb.Channel())
	}
}

func (s *rds) Publish(ctx context.Context, channel string, body []byte) error {
	return s.rds.Publish(ctx, channel, body).Err()
}
func (s *rds) publishExchange(ctx context.Context, exchange, key string, body []byte) error {
	return s.Publish(ctx, exchange, body)
}
