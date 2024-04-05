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
