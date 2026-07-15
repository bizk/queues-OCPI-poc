package main

import (
	"fmt"
	queue "queues-ocpi-poc/internal"

	"net/http"

	"github.com/gin-gonic/gin"
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

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "working..."})
	})

	router.POST("/event", func(c *gin.Context) {
		queue.Publish("event", []byte("something..."))
		c.JSON(http.StatusOK, gin.H{"message": "event received"})
	})

	router.GET("/events", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "events..."})
	})

	router.Run(":8080")
}
