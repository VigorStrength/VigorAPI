package main

import (
	"log"
	"os"

	"github.com/GhostDrew11/vigor-api/internal/api"
	"github.com/GhostDrew11/vigor-api/internal/config"
	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	//Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	//Connect to the database
	if err := db.ConnectDB(cfg); err != nil {
		log.Printf("Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}
	defer db.DisconnectDB()

	// Set up your Gin router
	router := gin.Default()

	// Set up you routes
	api.SetupRoutes(router)

	// Start you API server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
