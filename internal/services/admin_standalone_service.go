package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrStandAloneWorkoutAlreadyExists = fmt.Errorf("stand alone workout already exists")
	ErrStandAloneWorkoutNotFound      = fmt.Errorf("stand alone workout not found")
)

func (as *AdminService) CreateStandAloneWorkoutItem(ctx context.Context, standAloneWorkout models.StandAlone) error {
	standAloneWorkoutsCollection := as.database.Collection("standAloneWorkouts")

	filter := bson.M{"_id": standAloneWorkout.ID}
	count, err := standAloneWorkoutsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking if stand alone workout already exists: %w", err)
	}

	if count > 0 {
		return ErrStandAloneWorkoutAlreadyExists
	}

	_, err = standAloneWorkoutsCollection.InsertOne(ctx, standAloneWorkout)
	if err != nil {
		return fmt.Errorf("error inserting new stand alone workout: %w", err)
	}

	return nil
}

func (as *AdminService) UpdateStandAloneWorkoutItem(ctx context.Context, standAloneWorkoutID primitive.ObjectID, standAloneWorkoutInput models.StandAloneUpdateInput) error {
	// Check if standAloneWorkoutInput is empty
	if (models.StandAloneUpdateInput{}) == standAloneWorkoutInput {
		return fmt.Errorf("standAloneWorkoutInput cannot be empty")
	}

	standAloneWorkoutsCollection := as.database.Collection("standAloneWorkouts")

	filter := bson.M{"_id": standAloneWorkoutID}
	var standAloneWorkout models.StandAlone
	err := standAloneWorkoutsCollection.FindOne(ctx, filter).Decode(&standAloneWorkout)
	if err != nil {
		return fmt.Errorf("error finding stand alone workout: %w", err)
	}

	update := bson.M{"$set": standAloneWorkoutInput}
	_, err = standAloneWorkoutsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating stand alone workout: %w", err)
	}

	return nil
}

func (as *AdminService) GetStandAloneWorkoutItems(ctx context.Context) ([]models.StandAlone, error) {
	standAloneWorkoutsCollection := as.database.Collection("standAloneWorkouts")

	cursor, err := standAloneWorkoutsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding stand alone workouts: %w", err)
	}
	defer cursor.Close(ctx)

	var standAloneWorkouts []models.StandAlone
	if err := cursor.All(ctx, &standAloneWorkouts); err != nil {
		return nil, fmt.Errorf("error decoding stand alone workouts: %w", err)
	}

	return standAloneWorkouts, nil
}