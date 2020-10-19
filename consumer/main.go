package main

import (
	"fmt"
	"time"

	"github.com/muhammadisa/gorabbitmq"
	"github.com/streadway/amqp"
)

func handlerError(err error, panicking bool) {
	if err != nil {
		fmt.Println(err)
		if panicking {
			panic(err)
		}
	}
}

func delay(milis int64) {
	time.Sleep(time.Duration(milis) * time.Millisecond)
}

func main() {

	/**

	Failure Handling Solution

	Dead Latter Exchange (DLX)
	--------------------------

	Nack and Reject	:
	Ack the Message negatively and if requeue is true will be returned
	to exchange and continue by another queue if you configure DLX, if
	don't, the Message will be requeue as new Message from the exchange
	otherwise the requeue is false and no DLX configured will simply
	discard the Message.

	Ack	:
	Will remove the message from queue, even DLX is configured, because
	the Message Acknowleged positively

	**/

	connection, channel, err := gorabbitmq.Connector{
		Username: "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
	}.Dial()
	handlerError(err, true)
	defer connection.Close()
	defer channel.Close()

	messagesCustomer, err := gorabbitmq.Queue{
		QueueName: "foobar",
		Consumer:  "",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   true,
		NoWait:    true,
		Args: amqp.Table{
			"x-dead-letter-exchange": "email_exchanges.fanout",
		},
	}.Consume(channel)
	handlerError(err, true)

	messagesDriver, err := gorabbitmq.Queue{
		QueueName: "driver",
		Consumer:  "",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   true,
		NoWait:    true,
	}.Consume(channel)
	handlerError(err, true)

	foreverCustomer := make(chan bool)
	go func() {
		for d := range messagesCustomer {
			fmt.Println(fmt.Sprintf("Received Message From Customer : %s", d.Body))

			// Negative Ack and Requeue if failure happen
			err := d.Nack(false, true)
			handlerError(err, false)

			delay(250)
		}
	}()

	foreverDriver := make(chan bool)
	go func() {
		for d := range messagesDriver {
			fmt.Println(fmt.Sprintf("Received Message From Driver : %s", d.Body))

			// Positive Ack if no failure happen
			err := d.Ack(false)
			handlerError(err, false)

			delay(250)
		}
	}()

	<-foreverCustomer
	<-foreverDriver
}
