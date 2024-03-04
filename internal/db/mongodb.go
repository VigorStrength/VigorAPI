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
	// Users Collection
	usersCollection := Client.Database(dbName).Collection("users")

	// Index for 'email'
	emailIndexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // Unique index on 'email'
		Options: options.Index().SetUnique(true),
	}
	if _, err := usersCollection.Indexes().CreateOne(ctx, emailIndexModel); err != nil {
		log.Println("Error creating unique index for 'email' field in 'users' collection")
		return err
	}
	log.Println("Successfully created unique index for 'email' field in 'users' collection.")

	// Index for 'username'
	usernameIndexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1}, // Unique index on 'username'
		Options: options.Index().SetUnique(true),
	}
	if _, err := usersCollection.Indexes().CreateOne(ctx, usernameIndexModel); err != nil {
		log.Println("Error creating unique index for 'username' field in 'users' collection")
		return err
	}
	log.Println("Successfully created unique index for 'username' field in 'users' collection.")

	// Exercises Collection
	exercisesCollection := Client.Database(dbName).Collection("exercises")

	// Index for Exercise Name
	exerciseNameIndexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1}, // Unique index on 'name'
		Options: options.Index().SetUnique(true),
	}
	if _, err := exercisesCollection.Indexes().CreateOne(ctx, exerciseNameIndexModel); err != nil {
		log.Println("Error creating unique index for 'name' field in 'exercises' collection")
		return err
	}
	log.Println("Successfully created unique index for 'name' in 'exercises' collection")

	// Meals Collection
	mealsCollection := Client.Database(dbName).Collection("meals")

	// Index for Meal Name
	mealNameIndexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1}, // Unique index on 'name'
		Options: options.Index().SetUnique(true),
	}
	if _, err := mealsCollection.Indexes().CreateOne(ctx, mealNameIndexModel); err != nil {
		log.Println("Error creating unique index for 'name' field in 'meals' collection")
		return err
	}
	log.Println("Successfully created unique index for 'name' in 'meals' collection")

	// Additional collections and indexes can follow the same pattern as above.
	userMealStatusCollection := Client.Database(dbName).Collection("userMealStatus")
	userMealStatusIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "mealId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := userMealStatusCollection.Indexes().CreateOne(ctx, userMealStatusIndexModel); err != nil {
		log.Println("Error creating unique index for 'userId' and 'mealId' fields in 'userMealStatus' collection")
		return err
	}
	log.Println("Successfully created index for userMealStatuses collection")

	// Index for userCircuitStatus
	userCircuitStatusCollection := Client.Database(dbName).Collection("userCircuitStatus")
	userCircuitStatusIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "circuitId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := userCircuitStatusCollection.Indexes().CreateOne(ctx, userCircuitStatusIndexModel); err != nil {
		log.Println("Error creating unique index for 'userCircuitStatus' collection")
		return err
	}
	log.Println("Successfully created index for userCircuitStatus collection")

	// Index for userWorkoutDayStatus
	userWorkoutDayStatusCollection := Client.Database(dbName).Collection("userWorkoutDayStatus")
	userWorkoutDayStatusIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "workoutDayId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := userWorkoutDayStatusCollection.Indexes().CreateOne(ctx, userWorkoutDayStatusIndexModel); err != nil {
		log.Println("Error creating unique index for 'userWorkoutDayStatus' collection")
		return err
	}
	log.Println("Successfully created index for userWorkoutDayStatus collection")

	// Index for userWorkoutWeekStatus
	userWorkoutWeekStatusCollection := Client.Database(dbName).Collection("userWorkoutWeekStatus")
	userWorkoutWeekStatusIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "workoutWeekId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := userWorkoutWeekStatusCollection.Indexes().CreateOne(ctx, userWorkoutWeekStatusIndexModel); err != nil {
		log.Println("Error creating unique index for 'userWorkoutWeekStatus' collection", err)
		return err
	}
	log.Println("Successfully created index for userWorkoutWeekStatus collection")

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
