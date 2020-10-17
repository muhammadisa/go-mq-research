package main

import (
	"fmt"

	"github.com/muhammadisa/gorabbitmq"
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

	messagesCustomer, err := gorabbitmq.Queue{
		QueueName: "customer",
		Consumer:  "",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}.Consume(channel)
	handlerError(err, true)

	messagesDriver, err := gorabbitmq.Queue{
		QueueName: "driver",
		Consumer:  "",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}.Consume(channel)
	handlerError(err, true)

	foreverCustomer := make(chan bool)
	go func() {
		for d := range messagesCustomer {
			fmt.Println(fmt.Sprintf("Received Message From Customer : %s", d.Body))

			// err := d.Ack(false)
			// handlerError(err, false)
			// fmt.Println(fmt.Sprintf("Message acknowledged message"))
		}
	}()

	foreverDriver := make(chan bool)
	go func() {
		for d := range messagesDriver {
			fmt.Println(fmt.Sprintf("Received Message From Driver : %s", d.Body))

			// err := d.Ack(false)
			// handlerError(err, false)
			// fmt.Println(fmt.Sprintf("Message acknowledged message\n"))
		}
	}()

	<-foreverCustomer
	<-foreverDriver
}
