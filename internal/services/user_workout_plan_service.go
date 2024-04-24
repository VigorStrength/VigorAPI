package services

import (
	"context"
	"fmt"
	"time"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrAlreadyJoinded = fmt.Errorf("user has already joined this workout plan")
)

func (us *UserService) JoinWorkoutPlan(ctx context.Context, userID, workoutPlanID primitive.ObjectID) error {
	workoutPlanCollection := us.database.Collection("workoutPlans")
	userWorkoutPlanCollection := us.database.Collection("userWorkoutPlans")
	filter := bson.M{"_id": workoutPlanID}

	if err := workoutPlanCollection.FindOne(ctx, filter).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrWorkoutPlanNotFound
		}

		return fmt.Errorf("error finding workout plan: %w", err)
	}

	userWorkoutPlanFilter := bson.M{"userId": userID, "workoutPlanId": workoutPlanID}
	count, err := userWorkoutPlanCollection.CountDocuments(ctx, userWorkoutPlanFilter)
	if err != nil {
		return fmt.Errorf("error checking if user already joined the workout plan: %w", err)
	}
	if count > 0 {
		return ErrAlreadyJoinded
	}

	userWorkoutPlanStatus := models.UserWorkoutPlanStatus{
		ID: 		  primitive.NewObjectID(),
		UserID:       userID,
		WorkoutPlanID: workoutPlanID,
		StartDate:    time.Now(),
		Completed:    false,
	}

	if _, err := userWorkoutPlanCollection.InsertOne(ctx, userWorkoutPlanStatus); err != nil {
		return fmt.Errorf("error joining workout plan: %w", err)
	}

	return nil
}