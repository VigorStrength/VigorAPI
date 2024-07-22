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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//Load configuration
	cfg, err := config.LoadConfig()
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

	// Initialize database 
	database := dbService.Client.Database(cfg.DatabaseName)
	
	// Create required services
	handler := &utils.DefaultJWTHandler{}
	hasher := &utils.DefaultHasher{}
	parser := &utils.DefaultParser{}
	jwtService := utils.NewJWTService(cfg.JWTSecretKey, handler)
	adminService := services.NewAdminService(database, hasher, parser)
	userService := services.NewUserService(database, hasher, parser)

	// Set up your Gin router
	router := gin.Default()

	// Use Logger and Recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Set up CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:    []string{"*"}, //Allow all origins for the moment to be adjusted once the frontend is deployed
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-length", "Content-Type", "Authorization", "RefreshToken", "Accept", "Accept-Encoding", "User-Agent", "Host", "Connection",
		"Postman-Token", // Included Postman-Token to allow testing with Postman, remove in production,
		},
		ExposeHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard: true,
		MaxAge: 12 * time.Hour,
	}))

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