package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrExerciseAlreadyExists = fmt.Errorf("exercise already exists")
	ErrExerciseNotFound = fmt.Errorf("exercise not found")
)

func (as *AdminService) GetExerciseByID(ctx context.Context, exerciseID primitive.ObjectID) (models.Exercise, error) {
	//Get the exercise collection
	exerciseCollection := as.database.Collection("exercises")

	// Find the exercise by ID
	filter := bson.M{"_id": exerciseID}
	var exercise models.Exercise
	err := exerciseCollection.FindOne(ctx, filter).Decode(&exercise)
	if err != nil {
		return models.Exercise{}, fmt.Errorf("error finding exercise: %w", err)
	}

	return exercise, nil
}

func (as *AdminService) GetExercises(ctx context.Context) ([]models.Exercise, error) {
	//Get the exercise collection
	exerciseCollection := as.database.Collection("exercises")

	// Find all exercises
	cursor, err := exerciseCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding exercises: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the exercises
	var exercises []models.Exercise
	if err := cursor.All(ctx, &exercises); err != nil {
		return nil, fmt.Errorf("error decoding exercises: %w", err)
	}

	return exercises, nil
}

func (as *AdminService) CreateExercise(ctx context.Context, exerciseInput models.Exercise) error {
	//Get the exercise collection
	exerciseCollection := as.database.Collection("exercises")

	// check if the new exercise already exists
	filter := bson.M{"name": exerciseInput.Name}
	count, err := exerciseCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking if exercise already exists: %w", err)
	}

	if count > 0 {
		return ErrExerciseAlreadyExists
	}

	exerciseInput.ID = primitive.NewObjectID()

	// Insert the exercise into the database
	_, err = exerciseCollection.InsertOne(ctx, exerciseInput)
	if err != nil {
		return err
	}

	return nil
}

func (as *AdminService) UpdateExercise(ctx context.Context, exerciseID primitive.ObjectID, updateInput models.ExerciseUpdateInput) error {
	//Get the exercise collection
	exerciseCollection := as.database.Collection("exercises")
	filter := bson.M{"_id": exerciseID}
	
	// Convert the update input to a bson document
	updateDoc := as.parser.StructToBson(updateInput)

	// Use $set to only update the provided fields
	update := bson.M{"$set": updateDoc}
	result, err := exerciseCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating exercise: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrExerciseNotFound
	}

	return nil
}

func (as *AdminService) DeleteExercise(ctx context.Context, exerciseID primitive.ObjectID) error {
	//Get the exercise collection
	exerciseCollection := as.database.Collection("exercises")
	filter := bson.M{"_id": exerciseID}

	result, err := exerciseCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting exercise: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrExerciseNotFound
	}

	return nil
}

func (as *AdminService) SearchExercisesByName(ctx context.Context, name string) ([]models.Exercise, error) {
	//Get the exercise collection
	exerciseCollection := as.database.Collection("exercises")

	// Find all exercises that contain the name
	filter := bson.M{"name": primitive.Regex{Pattern: name, Options: "i"}}
	cursor, err := exerciseCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding exercises: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the exercises
	var exercises []models.Exercise
	if err := cursor.All(ctx, &exercises); err != nil {
		return nil, fmt.Errorf("error decoding exercises: %w", err)
	}

	return exercises, nil
}
