package main

import (
	"github.com/streadway/amqp"
	"log"
	"reflect"
)

type queue struct {
	url  string
	name string

	errorChannel chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
}

func NewQueue(url string, qName string) *queue {
	q := new(queue)

	q.url = url
	q.name = qName

	q.errorChannel = make(chan *amqp.Error)
	q.connect()
	go q.reconnect()

	return q
}

func (q queue) Send(message string) {

}

func (q queue) connect() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		//q.errorChannel <- err
		log.Println(reflect.TypeOf(err))
	}

	q.connection = conn
}

func (q queue) reconnect() {
	for {
		var err = <-q.errorChannel
		if err != nil {
			log.Printf("Error occured %s\nReconnecting to rabbitmq...", err)
			q.connect()
		}
	}
}
