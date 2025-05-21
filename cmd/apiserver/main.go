package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/akshayw1/antrea-renovate-demo/internal/api"
)

func main() {
	// Create a Gin router with a vulnerable version of Gin
	router := gin.Default()

	// Create API server
	apiServer := api.NewAPIServer()
	apiServer.Setup(router)

	// Add a health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(router.Run(":" + port))
}