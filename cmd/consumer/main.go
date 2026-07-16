package main

import (
	"context"
	"fmt"
	"queues-ocpi-poc/internal/db"
	"queues-ocpi-poc/internal/queue"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("[CONSUMER] ### Consumer start ###!")
	ctx := context.Background()
	defer ctx.Done()

	q, err := queue.NewNATSQueue()
	if err != nil {
		log.Error("[CONSUMER] error starting queue: %v\n", err)
		return
	}

	d, err := db.Connect()
	if err != nil {
		fmt.Printf("error connecting to MongoDB: %v\n", err)
		return
	}
	defer d.Disconnect(ctx)
	d.DropCollection(ctx, "events")

	err = q.Subscribe("event", func(m *nats.Msg) {
		log.Infof("[CONSUMER] Received event: %s", time.Now().Format(time.RFC3339))
		var event db.Event
		err := bson.Unmarshal(m.Data, &event)
		if err != nil {
			log.Errorf("[CONSUMER] error unmarshalling event: %v\n", err)
			m.Respond([]byte("Bad request"))
			return
		}

		log.Debugf("[CONSUMER] Event: %s", event.Name)

		err = d.InsertDocument(ctx, "events", event)
		if err != nil {
			log.Error("[CONSUMER] error inserting document: %v\n", err)
			m.Respond([]byte("Internal server error"))
			return
		}

		log.Debugf("[CONSUMER] > Event inserted")
		m.Respond([]byte("Ok"))
	})

	for true {
		if err != nil {
			log.Error("[CONSUMER] error subscribing to queue: %v\n", err)
			return
		}
		time.Sleep(5 * time.Second)
	}
}
