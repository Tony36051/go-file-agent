package main

import (
	"log"

	"github.com/Tony36051/go-file-agent/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()
	routers.SetupRoutes(r)

	// Start HTTP server for file download
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("HTTP server failed to start: %v", err)
	}
}
