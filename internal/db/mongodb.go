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
	if err := ensureIndexes(ctx); err != nil {
		return err
	}

	if err := InitializeCollections(cfg.DatabaseName); err != nil {
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
		log.Println("Error trying to create unique index for 'email' field in 'users' collection")
		return err
	}
	log.Println("Successfully created unique index for 'email' field in 'users' collection.")

	// Index for 'username'
	usernameIndexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1}, // Unique index on 'username'
		Options: options.Index().SetUnique(true),
	}
	if _, err := usersCollection.Indexes().CreateOne(ctx, usernameIndexModel); err != nil {
		log.Println("Error trying to create unique index for 'username' field in 'users' collection")
		return err
	}
	log.Println("Successfully created unique index for 'username' field in 'users' collection.")

	return nil
}

func InitializeCollections(dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := Client.Database(dbName)

	userSchema, err := schemaFiles.ReadFile("schemas/user/userSchema.json")
	if err != nil {
		return err
	}

	var userSchemaBson bson.M 
	if err := json.Unmarshal(userSchema, &userSchemaBson); err != nil {
		return err
	}

	workoutPlanSchema, err := schemaFiles.ReadFile("schemas/workoutPlan/workoutPlanSchema.json")
	if err != nil {
		return err
	}

	var workoutPlanSchemaBson bson.M
	if err := json.Unmarshal(workoutPlanSchema, &workoutPlanSchemaBson); err != nil {
		return err
	}

	mealPlanSchema, err := schemaFiles.ReadFile("schemas/mealPlan/mealPlanSchema.json")
	if err != nil {
		return err
	}

	var mealPlanSchemaBson bson.M
	if err := json.Unmarshal(mealPlanSchema, &mealPlanSchemaBson); err != nil {
		return err
	}

	messageSchema, err := schemaFiles.ReadFile("schemas/messaging/messageSchema.json")
	if err != nil {
		return err 
	}

	var messageSchemaBson bson.M 
	if err := json.Unmarshal(messageSchema, &messageSchemaBson); err != nil {
		return err
	}

	messagesGroupSchema, err := schemaFiles.ReadFile("schemas/messaging/groupSchema.json")
	if err != nil {
		return err
	}

	var messagesGroupSchemaBson bson.M
	if err := json.Unmarshal(messagesGroupSchema, &messagesGroupSchemaBson); err != nil {
		return err
	}

	//Applying user validation rules during collection creation or modification
	if err := applyCollectionValidation(ctx, db, "users", userSchemaBson); err != nil {
		return err
	}

	if err := applyCollectionValidation(ctx, db, "workoutPlans", workoutPlanSchemaBson); err != nil {
		return err
	}

	if err := applyCollectionValidation(ctx, db, "mealPlans", mealPlanSchemaBson); err != nil {
		return err
	}

	if err := applyCollectionValidation(ctx, db, "messages", messageSchemaBson); err != nil {
		return err
	}

	if err := applyCollectionValidation(ctx,db, "messagesGroups", messagesGroupSchemaBson); err != nil {
		return err
	}
	
	return nil 
}

//Helper function to apply validation rules to a collection
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


