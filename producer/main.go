package main

import (
	"fmt"
	"time"

	"github.com/muhammadisa/go-mq-notification/gorabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Message Queue Notification")

	connection, channel, err := gorabbitmq.Connector{
		Username: "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
	}.Dial()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

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
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println(fmt.Sprintf("Message Sent : %s\n", message))
			time.Sleep(500 * time.Millisecond)
		}
	}()
	<-forever

	defer connection.Close()
	defer channel.Close()

}
