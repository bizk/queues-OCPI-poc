package queue

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

type NATSQueue struct {
	nc *nats.Conn
}

func NewNATSQueue() (*NATSQueue, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		fmt.Printf("error connecting to NATS: %v\n", err)
		return nil, err
	}

	return &NATSQueue{nc: nc}, nil
}

func (q *NATSQueue) Publish(subject string, data []byte) error {
	fmt.Println("Publishing subject: ", subject)
	err := q.nc.Publish(subject, data)
	if err != nil {
		fmt.Printf("error publishing to NATS: %v\n", err)
		return err
	}

	return nil
}

func (q *NATSQueue) Subscribe(subject string, callback nats.MsgHandler) error {
	fmt.Println("Subscribing to subject: ", subject)
	_, err := q.nc.Subscribe(subject, callback)
	if err != nil {
		fmt.Printf("error subscribing to NATS: %v\n", err)
		return err
	}

	return nil
}
