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
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//Load configuration
	cfg, err := config.LoadConfig(false)
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// Instantiate the database service
	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v\n", err)
	}

	clientWrapper := db.NewMongoClientWrapper(mongoClient)
	dbService := db.NewMongoDBService(clientWrapper)

	// Create a context for the database connction
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect the database
	if err := dbService.ConnectDB(ctx, cfg); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
	}
	defer dbService.DisconnectDB(ctx)

	// Initilize the collections
	adminCollection := dbService.Client.Database(cfg.DatabaseName).Collection("admins")
	userCollection := dbService.Client.Database(cfg.DatabaseName).Collection("users")

	// Create a token service
	handler := &utils.DefaultJWTHandler{}
	hasher := &utils.DefaultHasher{}
	jwtService := utils.NewJWTService(cfg.JWTSecretKey, handler)
	adminService := services.NewAdminService(adminCollection, hasher)
	userService := services.NewUserService(userCollection, hasher)

	// Set up your Gin router
	router := gin.Default()

	// Set up your routes
	api.SetupRoutes(router, jwtService, *userService, *adminService)

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

// func main() {
// 	cfg, err := config.LoadConfig(false)
// 	if err != nil {
// 		log.Fatalf("Error loading config: %v\n", err)
// 	}

// 	log.Printf("Using database URI: %s\n", cfg.MongoDBURI)
// 	log.Printf("Using database name: %s\n", cfg.DatabaseName)
// }