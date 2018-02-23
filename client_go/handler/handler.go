package handler

import (
    "net/http"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "fmt"
    "github.com/streadway/amqp"
)

func Hello() echo.HandlerFunc {
    return func(c echo.Context) error {     //c をいじって Request, Responseを色々する 
        return c.String(http.StatusOK, "Hello World")
    }
}
func Clac() echo.HandlerFunc {
    return func(c echo.Context) error {
        calcValue := c.Param("calcValue")
        return c.String(http.StatusOK, "result:" + calcValue)
    }
}
func Send() echo.HandlerFunc {
    return func(c echo.Context) error {
        sendValue := c.Param("sendValue")

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

        // publishing a message
        err = ch.Publish(
            "", // exchange
            q.Name, // routing key
            false, // mandatory
            false, // immediate
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(sendValue),
            },
        )
        if err != nil {
            e.Logger.Fatalf("Failed to publish a message: %v", err)
        }

        fmt.Printf("send message-queue: %s¥n", sendValue)
        return c.String(http.StatusOK, sendValue)
    }
}
