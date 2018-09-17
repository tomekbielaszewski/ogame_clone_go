package main

import (
	"log"
	"time"
	"github.com/tomekbielaszewski/ogame_clone_go/utils"
)

func main() {
	queue := utils.NewQueue("amqp://guest:guest@localhost:5672/", "hello")
	defer queue.Close()

	queue.Consume(func(i string) {
		log.Printf("Received message: %s", i)
	})

	for i := 0; i < 10; i++ {
		log.Println("Sending message...")
		queue.Send("dupa")
		time.Sleep(500 * time.Millisecond)
	}
}
