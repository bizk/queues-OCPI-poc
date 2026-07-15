package main

import (
	"fmt"
	queue "queues-ocpi-poc/internal"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	fmt.Println("### Consumer start ###!")

	fmt.Println("Starting queue...")
	queue, err := queue.NewNATSQueue()
	if err != nil {
		fmt.Printf("error starting queue: %v\n", err)
		return
	}

	for true {
		fmt.Println("Waiting for events...")
		err := queue.Subscribe("event", func(m *nats.Msg) {
			m.Respond([]byte("event received"))
		})
		if err != nil {
			fmt.Printf("error subscribing to queue: %v\n", err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}
