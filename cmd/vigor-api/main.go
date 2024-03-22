package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GhostDrew11/vigor-api/internal/api"
	"github.com/GhostDrew11/vigor-api/internal/config"
	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	//Load configuration
	cfg, err := config.LoadConfig(false)
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// Instantiate the database service
	dbService := &db.MongoDBService{}

	// Create a context for the database connction
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect the database
	if err := dbService.ConnectDB(ctx, cfg); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
	}
	defer dbService.DisconnectDB(ctx)

	// Create a token service
	handler := &utils.DefaultJWTHandler{}
	jwtService := utils.NewJWTService(cfg.JWTSecretKey, handler)

	// Set up your Gin router
	router := gin.Default()

	// Set up your routes
	api.SetupRoutes(router, jwtService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start API server in a go routine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	// catch SIGINT and SIGTERM signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Doesn't block if no connections, but will wait for the duration of the context deadline
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exiting")
}
