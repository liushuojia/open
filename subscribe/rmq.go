package subscribe

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	reconnectDelay   = 8 * time.Second // 连接断开后多久重连
	reconnectMaxTime = 0               // 发送消息或订阅时，等待重连次数 0 一直重复连接
)

var _ Conn = (*rmq)(nil)
var _ RmqExchange = (*rmq)(nil)

type rmq struct {
	ctx    context.Context
	cancel context.CancelFunc

	url        string
	connection *amqp.Connection

	notifyConnect chan struct{}    // 连接成功通知
	notifyClose   chan *amqp.Error // 异常关闭通知
	channelCancel sync.Map         // 订阅后单独取消

	isRunning bool
	lock      sync.Mutex

	subscribeMap sync.Map // sync.map[channel] => func(context.Context, string, []byte) error) error
}

func NewRmq(user, password, host string, port int, vhost string) Conn {
	conn := &rmq{
		url:           fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, password, host, port, vhost),
		notifyConnect: make(chan struct{}),
		notifyClose:   make(chan *amqp.Error),
		isRunning:     false,
	}
	return conn
}
func NewRmqExchange(user, password, host string, port int, vhost string) RmqExchange {
	conn := &rmq{
		url:           fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, password, host, port, vhost),
		notifyConnect: make(chan struct{}),
		notifyClose:   make(chan *amqp.Error),
		isRunning:     false,
	}
	return conn
}

func (s *rmq) connect() error {
	var err error
	s.connection, err = amqp.Dial(s.url)
	if err != nil {
		return err
	}

	s.lock.Lock()
	s.isRunning = true
	s.lock.Unlock()

	// 这个必须在这里重新初始化
	s.notifyClose = make(chan *amqp.Error)
	s.connection.NotifyClose(s.notifyClose)
	close(s.notifyConnect)
	return nil
}
func (s *rmq) connectLoop() {
	i := 1
	for {
		if !s.IsRunning() {
			log.Println("[rabbitMQ]", "test to connect")
			if err := s.connect(); err != nil {
				log.Println("[rabbitMQ]", i, err.Error(), "Failed to connect rabbitMQ. Retrying...")
				i++
				time.Sleep(reconnectDelay)
			} else {
				i = 1
			}
		}
		select {
		case <-s.ctx.Done():
			goto END
		case <-s.notifyClose:
			s.init()
		}
	}
END:
	fmt.Println("[rabbitMQ]", "close")
}
func (s *rmq) init() {
	s.lock.Lock()
	s.isRunning = false
	s.lock.Unlock()

	s.notifyConnect = make(chan struct{})
	if s.connection != nil && !s.connection.IsClosed() {
		_ = s.connection.Close()
	}
}

func (s *rmq) Start(ctx context.Context) error {
	if s.isRunning {
		return nil
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	_ = s.Register("ping", func(ctx context.Context, channel string, body []byte) error {
		log.Println("PONG", channel, string(body))
		return nil
	})

	log.Println("[rabbitMQ]", "connect rabbitMQ")
	s.init()

	if err := s.connect(); err != nil {
		return errors.New("rabbitMQ connect fail")
	}

	// 监控
	go s.connectLoop()

	// 订阅
	go s.subscribe()

	return nil
}
func (s *rmq) Stop() error {
	if s.IsRunning() {
		s.cancel()
		s.init()
	}
	return nil
}

func (s *rmq) subscribe() {
	if !s.IsRunning() {
		return
	}

	channelList := make([]string, 0)
	s.subscribeMap.Range(func(key, value any) bool {
		if channel, ok := key.(string); ok {
			channelList = append(channelList, channel)
		}
		return true
	})

	if len(channelList) > 0 {
		if err := s.CreateQueueAction(channelList...); err != nil {
			log.Println("[rabbitMQ] create Queue fail err:", err.Error())
			return
		}
	}

	for _, channel := range channelList {
		go s.subscribeQueue(channel)
	}
}
func (s *rmq) IsRunning() bool {
	return s.isRunning
}

func (s *rmq) Register(channel string, fn func(context.Context, string, []byte) error) error {
	if _, ok := s.subscribeMap.Load(channel); !ok {
		if s.IsRunning() && s.connection != nil {
			if err := s.CreateQueueAction(channel); err != nil {
				log.Println("[rabbitMQ] create Queue fail err:", err.Error())
				return err
			}
			log.Println("[subscribe] rabbitMQ channel:", channel)
			go s.subscribeQueue(channel)
		}
	}
	s.subscribeMap.Store(channel, fn)
	return nil
}
func (s *rmq) UnRegister(channels ...string) {
	for _, channel := range channels {
		if _, ok := s.subscribeMap.LoadAndDelete(channel); !ok {
			return
		}
		if s.isRunning && s.connection != nil {
			log.Println("[subscribe] unsubscribe channel:", channel)
			if v, ok := s.channelCancel.LoadAndDelete(channel); ok {
				if vv, ok := v.(context.CancelFunc); ok && vv != nil {
					vv()
				}
			}
		}
	}
}

func (s *rmq) subscribeQueue(channel string) {
START:
	log.Println("[subscribe]", "rabbitMQ", "queue", channel, "start")
	c, err := s.connection.Channel()
	if err != nil {
		log.Println("[subscribe]", "get", "Channel", "err:", err)
		time.Sleep(reconnectDelay)
		goto START
	}

	message, err := c.Consume(channel, "", false, false, false, false, nil)
	if err != nil {
		log.Println("[subscribe]", "Consume", "queue", channel, err)
		time.Sleep(reconnectDelay)
		goto START
	}

	cancelCtx, cancel := context.WithCancel(s.ctx)
	s.channelCancel.Store(channel, cancel)

	for {
		select {
		case d, msgIsOpen := <-message:
			//if vv, ok := s.subscribeMap.Load(d.); ok {
			//	if l, ok := vv.([]CallBack); ok {
			//		for _, cb := range l {
			//			go cb.FN(s.ctx, msg.Channel, []byte(msg.Payload))
			//		}
			//	}
			//}
			if !msgIsOpen {
				break
			}

			v, ok := s.subscribeMap.Load(channel)
			if !ok {
				break
			}

			fn, ok := v.(func(context.Context, string, []byte) error)
			if !ok {
				break
			}

			if err := fn(s.ctx, d.RoutingKey, d.Body); err == nil {
				_ = d.Ack(true)
			}

		case <-c.NotifyClose(make(chan *amqp.Error)):
			log.Println("[subscribe]", "rabbitMQ", "channel", channel, "restart")
			goto START
		case <-s.notifyClose:
			log.Println("notifyClose")
			time.Sleep(reconnectDelay)
			goto START
		case <-s.ctx.Done():
			goto END
		case <-cancelCtx.Done():
			goto END
		}
	}
END:
	c.Close()
	log.Println("[subscribe]", "rabbitMQ", "queue", channel, "end")
}

/*
	publish
*/

func (s *rmq) Publish(ctx context.Context, queue string, body []byte) error {
	channel, err := s.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.Publish(
		"",    // exchange
		queue, // queue
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}
func (s *rmq) PublishExchange(ctx context.Context, exchange, key string, body []byte) error {
	channel, err := s.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.Publish(
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

/*
//		amqp.ExchangeFanout
//		amqp.ExchangeDirect
//		amqp.ExchangeTopic
//		amqp.ExchangeHeaders
//
//		topic 举例：
//		item.# ：能够匹配 item.insert.abc 或者 item.insert
//		item.* ：只能匹配 item.insert
*/

func (s *rmq) CreateQueueAction(nameList ...string) error {
	// 队列不存在创建
	for _, name := range nameList {
		_, err := s.CreateQueue(
			name,  // name 队列名称 为空时，名称随机
			true,  // durable 是否持久化
			false, // delete when unused 是否自动删除
			false, // exclusive 是否设置排他
			false, // no-wait 是否非阻塞
			nil,   // arguments 参数
		)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
//		name 		队列名称 为空时，名称随机
//		durable 	是否持久化
//		autoDelete 	delete when unused 是否自动删除
//		exclusive	是否设置排他
//		noWait		是否非阻塞
//		args		amqp.Table map[string]interface{} 参数
*/

func (s *rmq) CreateExchange(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	channel, err := s.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.ExchangeDeclare(
		name,       // name
		kind,       // type
		durable,    // durable
		autoDelete, // auto-deleted
		internal,   // internal true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		noWait,     // no-wait
		args,       // arguments
	)
}
func (s *rmq) CreateQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (queue amqp.Queue, err error) {
	channel, err := s.connection.Channel()
	if err != nil {
		return queue, err
	}
	defer channel.Close()

	return channel.QueueDeclare(
		name,       // name 队列名称 为空时，名称随机
		durable,    // durable 是否持久化
		autoDelete, // delete when unused 是否自动删除
		exclusive,  // exclusive 是否设置排他
		noWait,     // no-wait 是否非阻塞
		args,       // arguments 参数
	)
}
func (s *rmq) CreateExchangeBind(name string, exchange string, noWait bool, args amqp.Table, routingKeys ...string) error {
	channel, err := s.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	for _, k := range routingKeys {
		err := channel.ExchangeBind(
			name,     // name
			k,        // routing key
			exchange, // exchange
			noWait,
			args,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *rmq) CreateQueueBind(name string, exchange string, noWait bool, args amqp.Table, routingKeys ...string) error {
	channel, err := s.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	for _, k := range routingKeys {
		err := channel.QueueBind(
			name,     // queue name
			k,        // routing key
			exchange, // exchange
			noWait,
			args,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
