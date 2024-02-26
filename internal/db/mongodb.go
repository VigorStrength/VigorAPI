package db

import (
	"context"
	"log"
	"time"

	"github.com/GhostDrew11/vigor-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Client = client
	log.Println("Connected to MongoDB successfully.")

	// Ensure indexes after successful connection. If index creation fails, return the error.
	if err := ensureIndexes(ctx); err != nil {
		return err
	}

	return nil
}

func ensureIndexes(ctx context.Context) error {
	usersCollection := Client.Database("yourDatabaseName").Collection("users")

	// Index for 'email'
	emailIndexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // Unique index on 'email'
		Options: options.Index().SetUnique(true),
	}
	if _, err := usersCollection.Indexes().CreateOne(ctx, emailIndexModel); err != nil {
		return err
	}
	log.Println("Successfully created unique index for 'email' field in 'users' collection.")

	// Index for 'username'
	usernameIndexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1}, // Unique index on 'username'
		Options: options.Index().SetUnique(true),
	}
	if _, err := usersCollection.Indexes().CreateOne(ctx, usernameIndexModel); err != nil {
		return err
	}
	log.Println("Successfully created unique index for 'username' field in 'users' collection.")

	return nil
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
}
