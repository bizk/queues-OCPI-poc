package main

import (
	"context"
	"fmt"
	"queues-ocpi-poc/internal/db"
	"queues-ocpi-poc/internal/queue"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	fmt.Println("Hello, World!")

	router := gin.Default()

	fmt.Println("Starting Queue...")
	queue, err := queue.NewNATSQueue()
	if err != nil {
		fmt.Printf("error starting queue: %v\n", err)
		return
	}

	d, err := db.Connect()
	if err != nil {
		fmt.Printf("error connecting to MongoDB: %v\n", err)
		return
	}
	defer d.Disconnect(context.Background())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "working..."})
	})

	router.POST("/event", func(c *gin.Context) {
		fmt.Println("Publishing events...")
		for i := 0; i < 5; i++ {
			id := bson.NewObjectID()
			event := db.Event{ID: id, Name: fmt.Sprintf("test-%s-%d", id.Hex(), i)}
			eventBytes, err := bson.Marshal(event)
			if err != nil {
				fmt.Printf("error marshaling event: %v\n", err)
				continue
			}
			if i == 4 {
				eventBytes = []byte{1, 2, 16, 0} // Inject broken data
			}
			if err = queue.Publish("event", eventBytes); err != nil {
				fmt.Printf("error publishing event: %v\n", err)
				continue
			}

		}

		c.JSON(http.StatusOK, gin.H{"message": "event received"})
	})

	router.GET("/events", func(c *gin.Context) {
		documents, err := d.ListDocuments(context.Background(), "events")
		if err != nil {
			fmt.Printf("error listing documents: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error listing documents"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": documents})
	})

	router.GET("/dlq", func(c *gin.Context) {
		responses, err := queue.ListEDQ()
		if err != nil {
			fmt.Printf("error listing DLQ: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error listing DLQ"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": responses})
	})

	router.Run(":8080")
}
