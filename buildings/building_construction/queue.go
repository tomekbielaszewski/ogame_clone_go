package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type queue struct {
	url  string
	name string

	errorChannel chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
}

type consume func(string)

func NewQueue(url string, qName string) *queue {
	q := new(queue)
	q.url = url
	q.name = qName

	q.connect()
	go q.reconnector()

	return q
}

func (q *queue) Send(message string) {
	err := q.channel.Publish(
		"",     // exchange
		q.name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	logError("Sending message to queue failed", err)
}

func (q *queue) Consume(consumer consume) {
	log.Println("Registering consumer...")
	msgs, err := q.channel.Consume(
		q.name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	logError("Consuming messages from queue failed", err)
	log.Println("Consumer registered! Processing messages...")

	go func() {
		for msg := range msgs {
			consumer(string(msg.Body[:]))
		}
	}()
}

func (q *queue) Close() {
	q.channel.Close()
	q.connection.Close()
}

func (q *queue) reconnector() {
	for {
		err := <-q.errorChannel
		logError("Reconnecting after connection closed", err)

		q.connect()
	}
}

func (q *queue) connect() {
	for {
		log.Printf("Connecting to rabbitmq on %s\n", q.url)
		conn, err := amqp.Dial(q.url)
		if err == nil {
			q.connection = conn
			q.errorChannel = make(chan *amqp.Error)
			q.connection.NotifyClose(q.errorChannel)

			log.Println("Connection established!")

			q.openChannel()
			q.declareQueue()

			return
		}

		logError("Connection to rabbitmq failed. Retrying in 1 sec... ", err)
		time.Sleep(1000 * time.Millisecond)
	}
}

func (q *queue) declareQueue() {
	_, err := q.channel.QueueDeclare(
		q.name, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	logError("Queue declaration failed", err)
}

func (q *queue) openChannel() {
	channel, err := q.connection.Channel()
	logError("Opening channel failed", err)
	q.channel = channel
}

func logError(message string, err error) {
	if err != nil {
		log.Printf("%s: %s", message, err)
	}
}
