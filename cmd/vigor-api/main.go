package main

import (
	"github.com/GhostDrew11/vigor-api/internal/api"
	"github.com/gin-gonic/gin"
)


func main() {
	// Set up your Gin router
	router := gin.Default()

	// Connect to MOngoDB and pass it to your controllers if needed

	// Set up you routes 
	api.SetupRoutes(router)

	// Start you API server
	router.Run(":8080")
}