package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrWorkoutPlanAlreadyExists = fmt.Errorf("workout plan already exists")
)

func (as *AdminService) CreateWorkoutPlan(ctx context.Context, workoutPlanInput models.WorkoutPlan) error {
	// Get the workout plan collection
	workoutPlanCollection := as.database.Collection("workoutPlans")

	// check if the new workout plan already exists
	filter := bson.M{"name": workoutPlanInput.Name}
	count, err := workoutPlanCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking if workout plan already exists: %w", err)
	}

	if count > 0 {
		return ErrWorkoutPlanAlreadyExists
	}

	// Insert the new workout plan
	_, err = workoutPlanCollection.InsertOne(ctx, workoutPlanInput)
	if err != nil {
		return fmt.Errorf("error inserting workout plan: %w", err)
	}

	return nil
}