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

func main() {
	connection, channel, err := gorabbitmq.Connector{
		Username: "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
	}.Dial()
	handlerError(err, true)
	defer connection.Close()
	defer channel.Close()

	var iteration = 0
	forever := make(chan bool)
	go func() {
		for true {
			iteration++
			message := fmt.Sprintf("Message number %d", iteration)

			err = gorabbitmq.Message{
				ExchangeName: "email_exchanges.fanout",
				ExchangeKey:  "email.*",
				Mandatory:    false,
				Immediate:    false,
				Msg: amqp.Publishing{
					ContentType: gorabbitmq.TEXT,
					Body:        []byte(message),
				},
			}.Publish(channel)
			handlerError(err, true)

			fmt.Println(fmt.Sprintf("Message Sent : %s", message))
			time.Sleep(500 * time.Millisecond)
		}
	}()
	<-forever
}
