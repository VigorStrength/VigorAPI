package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrDailySupersetNotFound = fmt.Errorf("superset not found")
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

	exercisesMap := make(map[primitive.ObjectID]models.Exercise)
	for _, exercise := range exercises {
		exercisesMap[exercise.ID] = exercise
	}

	orderedExercises := make([]models.Exercise, len(exercisesIDs))
	for i, id := range exercisesIDs {
		if exercise, exists := exercisesMap[id]; exists {
			orderedExercises[i] = exercise
		} else {
			return nil, fmt.Errorf("exercise with ID %s not found", id.Hex())
		}
	}

	return orderedExercises, nil
}

func (us * UserService) GetDailySupersetByID(ctx context.Context, supersetID primitive.ObjectID) (*models.Superset, error) {
	var superset models.Superset
	supersetCollection := us.database.Collection("supersets")
	filter := bson.M{"_id": supersetID}
	err := supersetCollection.FindOne(ctx, filter).Decode(&superset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSupersetNotFound
		}
		return nil, fmt.Errorf("error finding superset: %w", err)
	}

	return &superset, nil
}

func (us *UserService) GetDailySupersetsByIDs(ctx context.Context, supersetsIDs []primitive.ObjectID) ([]models.Superset, error) {
	supersetCollection := us.database.Collection("supersets")

	filter := bson.M{"_id": bson.M{"$in": supersetsIDs}}
	cursor, err := supersetCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding supersets: %w", err)
	}
	defer cursor.Close(ctx)

	var supersets []models.Superset
	if err := cursor.All(ctx, &supersets); err != nil {
		return nil, fmt.Errorf("error decoding supersets: %w", err)
	}

	supersetsMap := make(map[primitive.ObjectID]models.Superset)
	for _, superset := range supersets {
		supersetsMap[superset.ID] = superset
	}

	orderedSupersets := make([]models.Superset, len(supersetsIDs))
	for i, id := range supersetsIDs {
		if superset, exists := supersetsMap[id]; exists {
			orderedSupersets[i] = superset
		} else {
			return nil, fmt.Errorf("superset with ID %s not found", id.Hex())
		}
	}

	return orderedSupersets, nil
}


	
