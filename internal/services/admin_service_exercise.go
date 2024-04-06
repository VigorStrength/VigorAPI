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
