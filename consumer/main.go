package main

import (
	"fmt"

	"github.com/muhammadisa/go-mq-notification/gorabbitmq"
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

	messagesCustomer, err := gorabbitmq.Queue{
		QueueName: "customer",
		Consumer:  "",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}.Consume(channel)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	messagesDriver, err := gorabbitmq.Queue{
		QueueName: "driver",
		Consumer:  "",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}.Consume(channel)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range messagesCustomer {
			fmt.Println(fmt.Sprintf("Received Message From Customer : %s\n", d.Body))
		}

		for d := range messagesDriver {
			fmt.Println(fmt.Sprintf("Received Message From Driver : %s\n", d.Body))
		}
	}()
	<-forever

	defer connection.Close()
	defer channel.Close()
}
