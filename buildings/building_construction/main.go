package main

import (
	"log"
	"time"
)

func main() {
	queue := NewQueue("amqp://guest:guest@localhost:5672/", "hello")
	defer queue.Close()

	queue.Consume(func(i string) {
		log.Printf("Received message: %s", i)
	})

	for {
		log.Println("Sending message...")
		queue.Send("dupa")
		time.Sleep(500 * time.Millisecond)
	}
}
