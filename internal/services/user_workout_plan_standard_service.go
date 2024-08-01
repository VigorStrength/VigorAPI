package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (us *UserService) GetWorkoutPlanByID(ctx context.Context, workoutPlanID primitive.ObjectID) (models.WorkoutPlan, error) {
	workoutPlanCollection := us.database.Collection("workoutPlans")

	filter := bson.M{"_id": workoutPlanID}
	var workoutPlan models.WorkoutPlan
	err := workoutPlanCollection.FindOne(ctx, filter).Decode(&workoutPlan)
	if err != nil {
		return models.WorkoutPlan{}, fmt.Errorf("error finding workout plan: %w", err)
	}

	return workoutPlan, nil
}

func  (uc *UserService) GetDailyExercisesByIDs(ctx context.Context, exercisesIDs []primitive.ObjectID) ([]models.Exercise, error) {
	exerciseCollection := uc.database.Collection("exercises")

	filter := bson.M{"_id": bson.M{"$in": exercisesIDs}}
	cursor, err := exerciseCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding exercises: %w", err)
	}
	defer cursor.Close(ctx)

	var exercises []models.Exercise
	if err := cursor.All(ctx, &exercises); err != nil {
		return nil, fmt.Errorf("error decoding exercises: %w", err)
	}

	return exercises, nil
}