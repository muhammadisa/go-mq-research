package gorabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

const (
	// TEXT content type
	TEXT = "text/plain"

	// JSON content type
	JSON = "application/json"
)

// IGoRabbitMQ interface
type IGoRabbitMQ interface {
	Dial() (*amqp.Connection, *amqp.Channel, error)
	CreateDirectExchange(channel *amqp.Channel, exchange Exchange) error
	CreateFanoutExchange(channel *amqp.Channel, exchange Exchange) error
	CreateHeaderExchange(channel *amqp.Channel, exchange Exchange) error
	CreateTopicsExchange(channel *amqp.Channel, exchange Exchange) error
	Consume(channel *amqp.Channel) (<-chan amqp.Delivery, error)
	Publish(channel *amqp.Channel) error
}

// Connector struct
type Connector struct {
	Username, Password, Host, Port string
}

// Exchange struct
type Exchange struct {
	ExchangeName                          string
	Durable, AutoDelete, Internal, NoWait bool
	Args                                  amqp.Table
}

// Message struct
type Message struct {
	ExchangeName, ExchangeKey string
	Mandatory, Immediate      bool
	Msg                       amqp.Publishing
}

// Queue struct
type Queue struct {
	QueueName, Consumer                 string
	AutoAck, Exclusive, NoLocal, NoWait bool
	Args                                amqp.Table
}

// Dial connecting to rabbitmq instance
func (rmqc Connector) Dial() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			rmqc.Username,
			rmqc.Password,
			rmqc.Host,
			rmqc.Port,
		),
	)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return conn, channel, nil
}

// Publish publish message with exchange name and exchange key
func (msg Message) Publish(channel *amqp.Channel) error {
	return channel.Publish(
		msg.ExchangeName,
		msg.ExchangeKey,
		msg.Mandatory,
		msg.Immediate,
		msg.Msg,
	)
}

// Consume consume messages inside queue
func (que Queue) Consume(channel *amqp.Channel) (<-chan amqp.Delivery, error) {
	return channel.Consume(
		que.QueueName,
		que.Consumer,
		que.AutoAck,
		que.Exclusive,
		que.NoLocal,
		que.NoWait,
		que.Args,
	)
}

// CreateDirectExchange create direct exchange
func CreateDirectExchange(channel *amqp.Channel, exchange Exchange) error {
	return channel.ExchangeDeclare(
		exchange.ExchangeName,
		amqp.ExchangeDirect,
		exchange.Durable,
		exchange.AutoDelete,
		exchange.Internal,
		exchange.NoWait,
		exchange.Args,
	)
}

// CreateFanoutExchange create fanout exchange
func CreateFanoutExchange(channel *amqp.Channel, exchange Exchange) error {
	return channel.ExchangeDeclare(
		exchange.ExchangeName,
		amqp.ExchangeFanout,
		exchange.Durable,
		exchange.AutoDelete,
		exchange.Internal,
		exchange.NoWait,
		exchange.Args,
	)
}

// CreateHeaderExchange create fanout exchange
func CreateHeaderExchange(channel *amqp.Channel, exchange Exchange) error {
	return channel.ExchangeDeclare(
		exchange.ExchangeName,
		amqp.ExchangeFanout,
		exchange.Durable,
		exchange.AutoDelete,
		exchange.Internal,
		exchange.NoWait,
		exchange.Args,
	)
}

// CreateTopicsExchange create fanout exchange
func CreateTopicsExchange(channel *amqp.Channel, exchange Exchange) error {
	return channel.ExchangeDeclare(
		exchange.ExchangeName,
		amqp.ExchangeTopic,
		exchange.Durable,
		exchange.AutoDelete,
		exchange.Internal,
		exchange.NoWait,
		exchange.Args,
	)
}
