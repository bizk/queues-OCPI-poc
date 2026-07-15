package main

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello, World!")

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "working..."})
	})

	router.POST("/event", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "event received"})
	})

	router.GET("/events", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "events..."})
	})

	router.Run(":8080")
}
