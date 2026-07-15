package main

import (
	"context"
	"fmt"
	"queues-ocpi-poc/internal/db"
	"queues-ocpi-poc/internal/queue"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	fmt.Println("### Consumer start ###!")
	ctx := context.Background()
	defer ctx.Done()

	fmt.Println("Starting queue...")
	q, err := queue.NewNATSQueue()
	if err != nil {
		fmt.Printf("error starting queue: %v\n", err)
		return
	}

	d, err := db.Connect()
	if err != nil {
		fmt.Printf("error connecting to MongoDB: %v\n", err)
		return
	}
	defer d.Disconnect(ctx)
	defer d.DropCollection(ctx, "events")

	err = q.Subscribe("event", func(m *nats.Msg) {
		fmt.Println("Received event: ", time.Now().UTC())
		var event db.Event
		err := bson.Unmarshal(m.Data, &event)
		if err != nil {
			fmt.Printf("error unmarshalling event: %v\n", err)
			return
		}

		fmt.Println("Event: ", event.Name)

		err = d.InsertDocument(ctx, "events", event)
		if err != nil {
			fmt.Printf("error inserting document: %v\n", err)
			return
		}

		m.Respond([]byte("event received"))
	})

	for true {
		if err != nil {
			fmt.Printf("error subscribing to queue: %v\n", err)
			return
		}
		time.Sleep(5 * time.Second)
	}
}
