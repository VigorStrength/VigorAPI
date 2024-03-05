package db

import (
	"context"
	"embed"
	"encoding/json"
	"log"
	"time"

	"github.com/GhostDrew11/vigor-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var schemaFiles embed.FS

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
	if err := ensureIndexes(ctx, cfg.DatabaseName); err != nil {
		return err
	}

	if err := InitializeCollections(ctx, cfg.DatabaseName); err != nil {
		return err
	}

	return nil
}

func ensureIndexes(ctx context.Context, dbName string) error {
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
		if err := createIndexes(ctx, dbName, collection, indexes); err != nil {
			return err
		}
	}
	log.Println("Successfully created indexes for all collctions.")
	return nil
}

func createIndexes(ctx context.Context, dbName, collectionName string, indexes []mongo.IndexModel) error {
	collection := Client.Database(dbName).Collection(collectionName)
	for _, indexModel := range indexes {
		if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
			log.Printf("Error creating index for collection '%s'!\n", collectionName)
			return err 
		}
		log.Printf("Successfully created index for collection '%s'.\n", collectionName)
	}
	return nil
}

func InitializeCollections(ctx context.Context, dbName string) error {
	db := Client.Database(dbName)

	schemas := []struct {
		collectionName string
		schemaFile     string
	}{
		{"users","schemas/user/userSchema.json"},
		{"exercises", "schemas/workoutPlan/exerciseSchema.json"},
		{"userCircuitStatus", "schemas/workoutPlan/userCircuitStatusSchema.json"},
		{"userWorkoutDayStatus", "schemas/workoutPlan/userWorkoutDayStatusSchema.json"},
		{"userWorkoutWeekStatus", "schemas/workoutPlan/userWorkoutWeekStatusSchema.json"},
		{"workoutPlans","schemas/workoutPlan/workoutPlanSchema.json"},
		{"meals", "schemas/mealPlan/mealSchema.json"},
		{"userMealStatus", "schemas/mealPlan/userMealStatusSchema.json"},
		{"mealPlans", "schemas/mealPlan/mealPlanSchema.json"},
		{"userDailyNutritionalLogs", "schemas/mealPlan/UserDailyNutritionalLogSchema.json"},
		{"messages", "schemas/messaging/groupSchema.json"},
		{"messagesGroups", "schemas/messaging/messageSchema.json"},
	}

	for _, s := range schemas {
		schema, err := schemaFiles.ReadFile(s.schemaFile)
		if err != nil {
			log.Printf("Error Reading %s schema file!\n", s.collectionName)
			return err
		}

		var schemaBson bson.M
		if err := json.Unmarshal(schema, &schemaBson); err != nil {
			log.Printf("Error unmarshalling %s schema!\n", s.collectionName)
			return err
		}

		if err := applyCollectionValidation(ctx, db, s.collectionName, schemaBson); err != nil {
			return err
		}
	}

	log.Println("Successfully initialized new collections with validation.")
	return nil
}

// Helper function to apply validation rules to a collection
func applyCollectionValidation(ctx context.Context, db *mongo.Database, collectionName string, schemaBson bson.M) error {
	opts := options.CreateCollection().SetValidator(bson.M{
		"$jsonSchema": schemaBson,
	})

	//Attempt to create the v=collection with validation rules
	err := db.CreateCollection(ctx, collectionName, opts)
	if mongo.IsDuplicateKeyError(err) {
		//If the collection already exists attempt to strenghten it with the new validation rules
		collModOpts := bson.D{
			{Key: "collMod", Value: collectionName},
			{Key: "validator", Value: bson.M{"$jsonSchema": schemaBson}},
		}
		return db.RunCommand(ctx, collModOpts).Err()
	} else if err != nil {
		log.Printf("Error applying validation rules to %s collection: %v\n", collectionName, err)
		return err
	}

	log.Printf("Successfully applied validation rules to %s collection.\n", collectionName)
	return nil
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
}
