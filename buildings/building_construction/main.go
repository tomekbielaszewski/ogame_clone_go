package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	newQueue := NewQueue("amqp://guest:guest@localhost:5672/", "hello")
	newQueue.Send("dupa")

	for {
		log.Println("sending..")
		newQueue.Send("dupa")
		time.Sleep(5000 * time.Millisecond)
	}
	//
	//conn := getRabbitmqConnection()
	//defer conn.Close()
	//
	//ch := getChannel(conn)
	//defer ch.Close()
	//
	//q, err := ch.QueueDeclare(
	//	"hello", // name
	//	false,   // durable
	//	false,   // delete when unused
	//	false,   // exclusive
	//	false,   // no-wait
	//	nil,     // arguments
	//)
	//failOnError(err, "Failed to declare a queue")
	//
	//body := "Hello World!"
	//err = ch.Publish(
	//	"",     // exchange
	//	q.Name, // routing key
	//	false,  // mandatory
	//	false,  // immediate
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte(body),
	//	})
	//failOnError(err, "Failed to publish a message")
}

func getChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func getRabbitmqConnection() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
