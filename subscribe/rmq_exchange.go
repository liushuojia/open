package subscribe

import (
	"github.com/streadway/amqp"
)

type RmqExchange interface {

	/*
		CreateExchange
			//		name 		队列名称 为空时，名称随机
			//		kind		类型 amqp.ExchangeFanout amqp.ExchangeDirect amqp.ExchangeTopic amqp.ExchangeHeaders
			//		durable 	是否持久化
			//		autoDelete 	delete when unused 是否自动删除
			//		exclusive	是否设置排他
			//		noWait		是否非阻塞
			//		args		amqp.Table map[string]interface{} 参数
	*/
	CreateExchange(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error

	/*
		CreateQueue
			//		name 		队列名称 为空时，名称随机
			//		durable 	是否持久化
			//		autoDelete 	delete when unused 是否自动删除
			//		exclusive	是否设置排他
			//		noWait		是否非阻塞
			//		args		amqp.Table map[string]interface{} 参数
	*/
	CreateQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (queue amqp.Queue, err error)

	/*
		CreateExchangeBind
			//		name 		队列名称 为空时，名称随机
			//		exchange	kind 名称
			//		noWait		是否非阻塞
			//		args		amqp.Table map[string]interface{} 参数
			// 		routingKeys	绑定key
	*/
	CreateExchangeBind(name string, exchange string, noWait bool, args amqp.Table, routingKeys ...string) error

	/*
		CreateQueueBind
			//		name 		Queue名称
			//		exchange	kind 名称
			//		noWait		是否非阻塞
			//		args		amqp.Table map[string]interface{} 参数
			// 		routingKeys	绑定key
	*/
	CreateQueueBind(name string, exchange string, noWait bool, args amqp.Table, routingKeys ...string) error
}
