package db

import (
	"context"
	"embed"
	"encoding/json"
	"log"

	"github.com/GhostDrew11/vigor-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:embed schemas/*
var files embed.FS

type MongoDBService struct {
	Client      MongoClient
	SchemaFiles embed.FS
}

func NewMongoDBService(client MongoClient) *MongoDBService {
	return &MongoDBService{Client: client, SchemaFiles: files}
}

func (ms *MongoDBService) ConnectDB(ctx context.Context, cfg *config.Config) error {
	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	_, err := ms.Client.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = ms.Client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB successfully.")

	db := ms.Client.Database(cfg.DatabaseName)
	if err := ms.InitializeCollections(ctx, db); err != nil {
		return err
	}

	if err := ms.EnsureIndexes(ctx, db); err != nil {
		return err
	}

	return nil
}

func (ms *MongoDBService) EnsureIndexes(ctx context.Context, db MongoDatabase) error {
	collections := map[string][]mongo.IndexModel{
		"users": {
			{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
			{Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true)},
		},
		"exercises": {
			{Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true)},
		},
		"meals": {
			{Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true)},
		},
		"userMealStatus": {
			{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "mealId", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"userCircuitStatus": {
			{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "circuitId", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"userWorkoutDayStatus": {
			{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "workoutDayId", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"userWorkoutWeekStatus": {
			{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "workoutWeekId", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"UserDailyNutritionalLogs": {
			{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "date", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
	}

	for collection, indexes := range collections {
		if err := ms.CreateIndexes(ctx, db, collection, indexes); err != nil {
			return err
		}
	}
	log.Println("Successfully created indexes for all collections.")
	return nil
}

func (ms *MongoDBService) CreateIndexes(ctx context.Context, db MongoDatabase, collectionName string, indexes []mongo.IndexModel) error {
	collection := db.Collection(collectionName)
	for _, indexModel := range indexes {
		if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
			log.Printf("Error creating index for collection '%s'!\n", collectionName)
			return err
		}
		log.Printf("Successfully created index for collection '%s'.\n", collectionName)
	}

	return nil
}

func (ms *MongoDBService) InitializeCollections(ctx context.Context, db MongoDatabase) error {
	schemas := []struct {
		collectionName string
		schemaFile     string
	}{
		{"users", "schemas/user/userSchema.json"},
		{"exercises", "schemas/workoutPlan/exerciseSchema.json"},
		{"userCircuitStatus", "schemas/workoutPlan/userCircuitStatusSchema.json"},
		{"userWorkoutDayStatus", "schemas/workoutPlan/userWorkoutDayStatusSchema.json"},
		{"userWorkoutWeekStatus", "schemas/workoutPlan/userWorkoutWeekStatusSchema.json"},
		{"workoutPlans", "schemas/workoutPlan/workoutPlanSchema.json"},
		{"meals", "schemas/mealPlan/mealSchema.json"},
		{"userMealStatus", "schemas/mealPlan/userMealStatusSchema.json"},
		{"mealPlans", "schemas/mealPlan/mealPlanSchema.json"},
		{"userDailyNutritionalLogs", "schemas/mealPlan/userDailyNutritionalLogSchema.json"},
		{"messages", "schemas/messaging/groupSchema.json"},
		{"messagesGroups", "schemas/messaging/messageSchema.json"},
	}

	for _, s := range schemas {

		schema, err := ms.SchemaFiles.ReadFile(s.schemaFile)
		if err != nil {
			log.Printf("Error Reading %s schema file!\n", s.collectionName)
			return err
		}

		var schemaBson bson.M
		if err := json.Unmarshal(schema, &schemaBson); err != nil {
			log.Printf("Error unmarshalling %s schema!\n", s.collectionName)
			return err
		}

		if err := ms.ApplyCollectionValidation(ctx, db, s.collectionName, schemaBson); err != nil {
			return err
		}
	}
	log.Println("Successfully initialized new collections with validation.")

	return nil
}

func (ms *MongoDBService) ApplyCollectionValidation(ctx context.Context, db MongoDatabase, collectionName string, schemaBson bson.M) error {
	opts := options.CreateCollection().SetValidator(schemaBson)
	//Attempt to create the v=collection with validation rules
	err := db.CreateCollection(ctx, collectionName, opts)
	if mongo.IsDuplicateKeyError(err) {
		//If the collection already exists attempt to strenghten it with the new validation rules
		collModOpts := bson.D{
			{Key: "collMod", Value: collectionName},
			{Key: "validator", Value: schemaBson},
		}
		return db.RunCommand(ctx, collModOpts).Err()
	} else if err != nil {
		log.Printf("Error applying validation rules to %s collection: %v\n", collectionName, err)
		return err
	}
	log.Printf("Successfully applied validation rules to %s collection.\n", collectionName)

	return nil
}

func (ms *MongoDBService) DisconnectDB(ctx context.Context) {
	if err := ms.Client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
}
