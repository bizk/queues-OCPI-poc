package queue

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
)

type NATSQueue struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewNATSQueue() (*NATSQueue, error) {
	log.Infof("[QUEUE] Starting NATS connection")
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		log.Errorf("[QUEUE] error connecting to NATS: %v\n", err)
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Errorf("[QUEUE] error creating JetStream context: %v\n", err)
		return nil, err
	}

	return &NATSQueue{nc: nc, js: js}, nil
}

func (q *NATSQueue) Subscribe(subject string, callback nats.MsgHandler) error {
	log.Infof("[QUEUE] Subscribing to subject: %s", subject)
	_, err := q.nc.Subscribe(subject, callback)
	if err != nil {
		log.Errorf("[QUEUE] error subscribing to NATS: %v\n", err)
		return err
	}

	return nil
}

func (q *NATSQueue) Publish(subject string, data []byte) error {
	log.Infof("[QUEUE] Publishing event")

	response, err := q.nc.Request(subject, data, time.Second*2)
	if err != nil {
		log.Errorf("[QUEUE] error requesting from NATS: %v\n", err)
		return err
	}

	if string(response.Data) != "Ok" {
		log.Errorf("[QUEUE] error response from NATS: %s", string(response.Data))
		return fmt.Errorf("error response from NATS: %s", string(response.Reply))
	}

	return nil
}

func (q *NATSQueue) PushToDLQ(data []byte) error {
	log.Infof("[QUEUE] Pushing event to DLQ")

	if err := q.nc.Publish("DLQ", data); err != nil {
		log.Errorf("[QUEUE] error publishing to DLQ: %v\n", err)
		return err
	}

	return nil
}

func (q *NATSQueue) ListEDQ() ([]string, error) {
	log.Infof("[QUEUE] Listing DLQ")

	// Todo move this out
	sub, err := q.js.PullSubscribe("DLQ", "DLQ", nats.BindStream("DLQ"))
	if err != nil {
		log.Errorf("[QUEUE] error creating pull subscription to DLQ: %v\n", err)
		return nil, err
	}
	defer sub.Unsubscribe()

	// Attempt to fetch messages without "destructively" draining the queue
	msgs, err := sub.Fetch(15, nats.MaxWait(5*time.Second))
	if err != nil && err != nats.ErrTimeout {
		log.Errorf("[QUEUE] error fetching from DLQ: %v\n", err)
		return nil, err
	}

	if len(msgs) == 0 {
		log.Debugf("[QUEUE] No messages in DLQ")
		return []string{}, nil
	}

	responses := make([]string, len(msgs))
	for _, msg := range msgs {
		log.Debugf("[QUEUE] DLQ message: %s", string(msg.Data))
		responses = append(responses, string(msg.Data))
	}

	return responses, nil
}
