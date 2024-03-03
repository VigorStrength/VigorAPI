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
        Keys: bson.M{"email": 1}, // Unique index on 'email'
        Options: options.Index().SetUnique(true),
    }
    if _, err := usersCollection.Indexes().CreateOne(ctx, emailIndexModel); err != nil {
        log.Println("Error trying to create unique index for 'email' field in 'users' collection")
        return err
    }
    log.Println("Successfully created unique index for 'email' field in 'users' collection.")

    // Index for 'username'
    usernameIndexModel := mongo.IndexModel{
        Keys: bson.M{"username": 1}, // Unique index on 'username'
        Options: options.Index().SetUnique(true),
    }
    if _, err := usersCollection.Indexes().CreateOne(ctx, usernameIndexModel); err != nil {
        log.Println("Error trying to create unique index for 'username' field in 'users' collection")
        return err
    }
    log.Println("Successfully created unique index for 'username' field in 'users' collection.")

    // Exercises Collection
    exercisesCollection := Client.Database(dbName).Collection("exercises")
    
    // Index for Exercise Name
    exerciseNameIndexModel := mongo.IndexModel{
        Keys: bson.M{"name": 1}, // Unique index on 'name'
        Options: options.Index().SetUnique(true),
    }
    if _, err := exercisesCollection.Indexes().CreateOne(ctx, exerciseNameIndexModel); err != nil {
        log.Println("Error trying to create unique index for 'name' field in 'exercises' collection")
        return err
    }
    log.Println("Successfully created unique index for 'name' in 'exercises' collection")

    // Meals Collection
    mealsCollection := Client.Database(dbName).Collection("meals")
    
    // Index for Meal Name
    mealNameIndexModel := mongo.IndexModel{
        Keys: bson.M{"name": 1}, // Unique index on 'name'
        Options: options.Index().SetUnique(true),
    }
    if _, err := mealsCollection.Indexes().CreateOne(ctx, mealNameIndexModel); err != nil {
        log.Println("Error trying to create unique index for 'name' field in 'meals' collection")
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
		log.Println("Error trying to create unique index for 'userId' and 'mealId' fields in 'userMealStatus' collection")
		return err
	}
	log.Println("Successfully created index for userMealStatuses collection")

    return nil
}

func InitializeCollections(ctx context.Context, dbName string) error {
	db := Client.Database(dbName)

	userSchema, err := schemaFiles.ReadFile("schemas/user/userSchema.json")
	if err != nil {
		return err
	}

	var userSchemaBson bson.M 
	if err := json.Unmarshal(userSchema, &userSchemaBson); err != nil {
		return err
	}

	exerciseSchema, err := schemaFiles.ReadFile("schemas/workoutPlan/exerciseSchema.json")
	if err != nil {
		return err
	}

	var exerciseSchemaBson bson.M
	if err := json.Unmarshal(exerciseSchema, &exerciseSchemaBson); err != nil {
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

	mealSchema, err := schemaFiles.ReadFile("schemas/mealPlan/mealSchema.json")
	if err != nil {
		return err 
	}

	var mealSchemaBson bson.M
	if err := json.Unmarshal(mealSchema, &mealSchemaBson); err != nil {
		return err
	}

	userMealStatusSchema, err := schemaFiles.ReadFile("schemas/mealPlan/userMealStatusSchema.json")
	if err != nil {
		return err 
	}

	var userMealStatusSchemaBson bson.M
	if err := json.Unmarshal(userMealStatusSchema, &userMealStatusSchemaBson); err != nil {
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

	if err := applyCollectionValidation(ctx, db, "exercises", exerciseSchemaBson); err != nil {
		return err
	}

	if err := applyCollectionValidation(ctx, db, "workoutPlans", workoutPlanSchemaBson); err != nil {
		return err
	}

	if err := applyCollectionValidation(ctx, db, "meals", mealSchemaBson); err != nil {
		return err
	}

	if err := applyCollectionValidation(ctx, db, "meals", userMealStatusSchemaBson); err != nil {
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


