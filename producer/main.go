package main

import (
	"fmt"
	"math/rand"
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
	connection, channel, err := gorabbitmq.Connector{
		Username: "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
	}.Dial()
	handlerError(err, true)
	defer connection.Close()
	defer channel.Close()

	greetings := []string{
		"Bonjour",
		"Hola",
		"Hello",
		"Halo",
		"Kon'nichiwa",
		"Yeoboseyo",
		"Ni Hao",
		"Sawadikap",
	}
	forever := make(chan bool)
	go func() {
		for true {
			greeting := fmt.Sprintf(`{"foobar_content":"%s"}`, greetings[rand.Intn(len(greetings))])
			err = gorabbitmq.Message{
				ExchangeName: "foobar_exchanges",
				ExchangeKey:  "",
				Mandatory:    false,
				Immediate:    false,
				Msg: amqp.Publishing{
					ContentType: gorabbitmq.TEXT,
					Body:        []byte(greeting),
				},
			}.Publish(channel)
			handlerError(err, true)
			fmt.Println(greeting)
			delay(1000)
		}
	}()
	<-forever
}
