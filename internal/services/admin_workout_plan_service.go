package services

import (
	"context"
	"fmt"
	"log"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrWorkoutPlanAlreadyExists = fmt.Errorf("workout plan already exists")
	ErrWorkoutPlanNotFound = fmt.Errorf("workout plan not found")
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

func (as *AdminService) UpdateWorkoutPlan(ctx context.Context, workoutPlanID primitive.ObjectID, updateInput models.WorkoutPlanInput) error {
	// Get the workout plan collection 
	workoutPlanCollection := as.database.Collection("workoutPlans")
	filter := bson.M{"_id": workoutPlanID}

	// Convert the update input to a BSON document
	updateDoc := as.parser.StructToBson(updateInput)
	log.Printf("Document: %v\n", updateDoc)

	// Use $set to only update the provided fields
	update := bson.M{"$set": updateDoc}
	result, err := workoutPlanCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating workout plan: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrWorkoutPlanNotFound
	}

	return nil
}

func (as *AdminService) DeleteWorkoutPlan(ctx context.Context, workoutPlanID primitive.ObjectID) error {
	// Get the workout plan collection
	workoutPlanCollection := as.database.Collection("workoutPlans")
	filter := bson.M{"_id": workoutPlanID}

	// Delete the workout plan
	result, err := workoutPlanCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting workout plan: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrWorkoutPlanNotFound
	}

	return nil
}