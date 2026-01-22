package subscribe

//import (
//	"context"
//	"fmt"
//	"time"
//
//	"github.com/streadway/amqp"
//)
//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"sync"
//	"time"
//
//	log "github.com/sirupsen/logrus"
//	"github.com/streadway/amqp"
//)
//
//const (
//	reconnectDelay   = 8 * time.Second // 连接断开后多久重连
//	reconnectMaxTime = 0               // 发送消息或订阅时，等待重连次数 0 一直重复连接
//)
//
//var _ Conn = (*rmq)(nil)
//
//type rmq struct {
//	ctx    context.Context
//	cancel context.CancelFunc
//
//	url        string
//	connection *amqp.Connection
//	channel    *amqp.Channel
//
//	errChan       chan error    // 全局错误通道
//	notifyConnect chan struct{} // 连接成功通知
//	isRunning     bool
//	lock          sync.Mutex
//	subscribeMap  map[string]map[string]CallBack
//}
//
//func NewRmq(user, password, host string, port int, vhost string) Conn {
//	return &rmq{
//		url:          fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, password, host, port, vhost),
//		isRunning:    false,
//		errChan:      make(chan error, 1),
//		subscribeMap: make(map[string]map[string]CallBack),
//	}
//}
//
//// monitorConnection 监控连接状态，断开时触发重连
//func (s *rmq) monitorConnection() {
//	var err error
//	for {
//		select {
//		case err := <-s.errChan: // 通道断开
//			if err == nil {
//				goto loop
//			}
//			log.Println("[RabbitMQ]", "channel 断开重连")
//			if !s.connection.IsClosed() {
//				goto reConnect
//			}
//			goto resetChannel
//		case <-s.connection.NotifyClose(make(chan *amqp.Error, 1)): //连接断开
//			goto reConnect
//		case <-s.ctx.Done():
//			goto end
//		}
//	reConnect:
//		s.init()
//		if err := s.connect(); err != nil {
//			log.Printf("重连失败: %v", err)
//			goto loop
//		}
//		goto resetChannel
//	resetChannel:
//		if !s.connection.IsClosed() {
//			goto reConnect
//		}
//		s.channel, err = s.connection.Channel()
//		if err == nil {
//			// 监控信道
//			go s.monitorChannel()
//			// 重新订阅
//			go s.subscribe()
//		}
//		goto loop
//	loop:
//		time.Sleep(reconnectDelay)
//	}
//end:
//}
//
//// monitorChannel 监控信道
//func (s *rmq) monitorChannel() {
//	select {
//	case err := <-s.channel.NotifyClose(make(chan *amqp.Error, 1)):
//		log.Printf("信道关闭: %v", err)
//		s.errChan <- errors.New("channel closed " + err.Error())
//	case <-s.ctx.Done():
//	}
//}
//
//func (s *rmq) connect() error {
//	var err error
//	s.connection, err = amqp.Dial(s.url)
//	if err != nil {
//		return err
//	}
//
//	s.lock.Lock()
//	s.isRunning = true
//	s.lock.Unlock()
//	return nil
//}
//func (s *rmq) init() {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//	s.isRunning = false
//	s.notifyConnect = make(chan struct{})
//
//	// 关闭旧连接/信道（防止资源泄漏）
//	if s.channel != nil {
//		_ = s.channel.Close()
//	}
//	if s.connection != nil && !s.connection.IsClosed() {
//		_ = s.connection.Close()
//	}
//}
//
//func (s *rmq) Start(ctx context.Context) error {
//	if s.isRunning {
//		return nil
//	}
//
//	s.ctx, s.cancel = context.WithCancel(ctx)
//	_ = s.Register("ping", "key", func(ctx context.Context, channel string, body []byte) error {
//		log.Println("PONG", channel, string(body))
//		return nil
//	})
//
//	log.Println("[rabbitMQ]", "connect rabbitMQ")
//	s.init()
//
//	var err error
//	if err = s.connect(); err != nil {
//		return errors.New("rabbitMQ connect fail")
//	}
//	if s.channel, err = s.connection.Channel(); err != nil {
//		return errors.New("rabbitMQ get channel fail")
//	}
//
//	// 订阅
//	go s.subscribe()
//
//	// 监控
//	go s.monitorConnection()
//	return nil
//}
//func (s *rmq) Stop() error {
//	if s.IsRunning() {
//		s.cancel()
//		s.init()
//	}
//	return nil
//}
//func (s *rmq) subscribe() {
//	if s.isRunning {
//		return
//	}
//
//	channelList := make([]string, 0)
//	for channel := range s.subscribeMap {
//		channelList = append(channelList, channel)
//	}
//	if len(channelList) > 0 {
//		if err := s.CreateQueue(channelList...); err != nil {
//			return
//		}
//	}
//
//	/*
//		log.Println("[rabbitMQ] subscribe channel:", channelList)
//		s.sub = s.rds.Subscribe(s.ctx, channelList...)
//
//		var (
//			ch = s.sub.Channel()
//		)
//		defer func() {
//			l := make([]string, 0)
//			for channel := range s.subscribeMap {
//				l = append(l, channel)
//			}
//			log.Println("[redis] close subscribe", l)
//			_ = s.sub.Close()
//		}()
//		for {
//			select {
//			case msg := <-ch:
//				if l, ok := s.subscribeMap[msg.Channel]; ok {
//					for _, f := range l {
//						go f(s.ctx, msg.Channel, msg.Payload)
//					}
//				}
//			case <-s.ctx.Done():
//				s.lock.Lock()
//				s.isRunning = false
//				s.lock.Unlock()
//				return
//			}
//		}
//	*/
//}
//func (s *rmq) IsRunning() bool {
//	return s.isRunning
//}
//func (s *rmq) Register(channel, key string, callback CallBack) error {
//
//	//s.lock.Lock()
//	//defer s.lock.Unlock()
//	//
//	//if v, ok := s.subscribeMap[channel]; ok {
//	//	if _, ok := v[key]; ok {
//	//		s.subscribeMap[channel][key] = callback
//	//	}
//	//	return nil
//	//}
//	//
//	//s.subscribeMap[channel] = make(map[string]CallBack)
//	//s.subscribeMap[channel][key] = callback
//	//if s.isRunning {
//	//	log.Println("[redis] subscribe channel:", channel)
//	//	return s.sub.Subscribe(s.ctx, channel)
//	//}
//
//	return nil
//}
//func (s *rmq) UnRegister(channel, key string) {
//	/*
//		if _, ok := s.subscribeMap[channel]; !ok {
//			return
//		}
//
//		s.lock.Lock()
//		defer s.lock.Unlock()
//
//		delete(s.subscribeMap[channel], key)
//		if len(s.subscribeMap[channel]) <= 0 && s.isRunning {
//			delete(s.subscribeMap, channel)
//			log.Println("[redis] unsubscribe channel:", channel)
//			_ = s.sub.Unsubscribe(s.ctx, channel)
//		}
//	*/
//}
//
//func (s *rmq) SubscribeQueue(callback CallBack, name string) {
//START:
//	channel, err := s.connection.Channel()
//	if err != nil {
//		log.Println("channel", "queue", name, err)
//		time.Sleep(reconnectDelay)
//		goto START
//	}
//
//	notifyClose := make(chan *amqp.Error)
//	channel.NotifyClose(notifyClose)
//
//	log.Println("[subscribe]", "rabbitMQ", "queue", name, "start")
//	if err := s.CreateQueue(name); err != nil {
//		log.Println("CreateQueue", "queue", name, err)
//		time.Sleep(reconnectDelay)
//		goto START
//	}
//
//	message, err := channel.Consume(name, "", false, false, false, false, nil)
//	if err != nil {
//		log.Println("Consume", "queue", name, err)
//		time.Sleep(reconnectDelay)
//		goto START
//	}
//
//	for {
//		select {
//		case d, msgIsOpen := <-message:
//			if !msgIsOpen {
//				break
//			}
//			if err := callback(s.ctx, d.RoutingKey, d.Body); err == nil {
//				_ = d.Ack(true)
//			}
//		case <-notifyClose:
//			log.Println("[subscribe]", "rabbitMQ", "queue", name, "restart")
//			goto START
//		case <-s.ctx.Done():
//			goto END
//		}
//	}
//END:
//	channel.Close()
//	log.Println("[subscribe]", "rabbitMQ", "queue", name, "end")
//}
//func (s *rmq) CreateQueue(nameList ...string) error {
//	// 队列不存在创建
//	fmt.Println(nameList)
//	for _, name := range nameList {
//		_, err := s.channel.QueueDeclare(
//			name,  // name 队列名称 为空时，名称随机
//			true,  // durable 是否持久化
//			false, // delete when unused 是否自动删除
//			false, // exclusive 是否设置排他
//			false, // no-wait 是否非阻塞
//			nil,   // arguments 参数
//		)
//		return err
//	}
//	return nil
//}
//
//func (s *rmq) Publish(ctx context.Context, channel string, body []byte) error {
//	return s.channel.Publish(
//		"",      // exchange
//		channel, // name
//		false,   // mandatory
//		false,   // immediate
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        body,
//		},
//	)
//}
//func (s *rmq) publishExchange(ctx context.Context, exchange, key string, body []byte) error {
//	return s.channel.Publish(
//		exchange, // exchange
//		key,      // routing key
//		false,    // mandatory
//		false,    // immediate
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        body,
//		},
//	)
//}
