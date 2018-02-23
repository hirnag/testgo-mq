package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"fmt"
	"github.com/streadway/amqp"
)

var port = ":5000"

type server struct{}

func main() {
	fmt.Println("start - main")

	e := echo.New()
	e.Use(middleware.Logger())

	// making a connection
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		e.Logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// making a chanel
	ch, err := conn.Channel()
	if err != nil {
		e.Logger.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// declaring a queue
	q, err := ch.QueueDeclare(
		"mq-name", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait (wait time for processing)
		nil, // arguments
	)
	if err != nil {
		e.Logger.Fatalf("Failed to declare a queue: %v", err)
	}

	// consuming a message
	msgs, err := ch.Consume(
		q.Name, // queue
		"", // consumer
		true, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
	if err != nil {
		e.Logger.Fatalf("Failed to register a consumer: %v", err)
	}
	for d := range msgs {
		fmt.Printf("Received a message: %s", d.Body) //any kind of further processing code
		fmt.Println("")
	}

	fmt.Println("end   - main")
}

