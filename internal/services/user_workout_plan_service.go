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
	ErrCircuitNotFound = fmt.Errorf("user circuit not found")
	ErrWrokoutWeekNotFound = fmt.Errorf("user workout week not found")
	ErrAlreadyJoinded = fmt.Errorf("user has already joined this workout plan")
	ErrExerciseAlreadyCompleted = fmt.Errorf("exercise has already been completed")
	ErrActiveWorkoutPlanNotFound = fmt.Errorf("the user hasn't joined any workout plan with that ID")
)

func (us *UserService) GetActiveWorkoutPlan(ctx context.Context, userID primitive.ObjectID) (*models.UserWorkoutPlanStatus, error) {
	var activeWorkoutPlan models.UserWorkoutPlanStatus
	userWorkoutPlanStatusCollection := us.database.Collection("userWorkoutPlanStatus")
	filter := bson.M{"userId": userID, "completed": false}
	err := userWorkoutPlanStatusCollection.FindOne(ctx, filter).Decode(&activeWorkoutPlan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrActiveWorkoutPlanNotFound
		}
		return nil, fmt.Errorf("error finding active workout plan: %w", err)
	}

	//Combining it with progress from the workoutPlan to have be just as one call to the server
	totalDays, completedDays, err := us.progressUtil(ctx, userID, activeWorkoutPlan.WorkoutPlanID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrActiveWorkoutPlanNotFound
		}
		return nil, fmt.Errorf("error getting progress: %w", err)
	}

	if totalDays == 0 {
		activeWorkoutPlan.Progress = 0
	} else {
		activeWorkoutPlan.Progress = completedDays / totalDays * 100
	}

	return &activeWorkoutPlan, nil
}

func (us *UserService) JoinWorkoutPlan(ctx context.Context, userID, workoutPlanID primitive.ObjectID) error {
	var workoutPlan models.WorkoutPlan
	workoutPlanCollection := us.database.Collection("workoutPlans")
	err := workoutPlanCollection.FindOne(ctx, bson.M{"_id": workoutPlanID}).Decode(&workoutPlan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrWorkoutPlanNotFound
		}
		return fmt.Errorf("error finding workout plan: %w", err)
	}

	userWorkoutPlanStatus := models.NewUserWorkoutPlanStatus(userID, workoutPlanID, workoutPlan.Name)
	if _, err := us.database.Collection("userWorkoutPlanStatus").InsertOne(ctx, userWorkoutPlanStatus); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrAlreadyJoinded
		}
		return fmt.Errorf("error inserting user workout plan status: %w", err)
	}

	for _, week := range workoutPlan.Weeks {
		userWeekStatus := models.NewUserWorkoutWeekStatus(userID, week.ID, workoutPlanID)
		if _, err := us.database.Collection("userWorkoutWeekStatus").InsertOne(ctx, userWeekStatus); err != nil {
			return fmt.Errorf("error inserting user workout week status: %w", err)
		}

		for _, day := range week.Days {
			userDayStatus := models.NewUserWorkoutDayStatus(userID, day.ID, week.ID, workoutPlanID)
			if _, err := us.database.Collection("userWorkoutDayStatus").InsertOne(ctx, userDayStatus); err != nil {
				return fmt.Errorf("error inserting user workout day status: %w", err)
			}

			for _, circuit := range append(day.WarmUps, append(day.Workouts, day.CoolDowns...)...) {
				userCircuitStatus := models.NewUserCircuitStatus(userID, circuit.ID, day.ID, workoutPlanID)
				if _, err := us.database.Collection("userCircuitStatus").InsertOne(ctx, userCircuitStatus); err != nil {
					return fmt.Errorf("error inserting user circuit status: %w", err)
				}

				for _, exerciseID := range circuit.ExerciseIDs {
					userExerciseStatus := models.NewUserExerciseStatus(userID, exerciseID, circuit.ID, workoutPlanID)
					if _, err := us.database.Collection("userExerciseStatus").InsertOne(ctx, userExerciseStatus); err != nil {
						return fmt.Errorf("error inserting user exercise status: %w", err)
					}
				}
			}
		}
	}

	return nil
}

func (us *UserService) MarkExerciseAsCompleted(ctx context.Context, userID, exerciseID, circuitID primitive.ObjectID, logs []models.UserExerciseLogInput) error {
    
	workoutPlanID, err := us.getAndValidateCircuit(ctx, circuitID)
    if err != nil {
		if err == mongo.ErrNoDocuments {
			return  ErrCircuitNotFound
		}

        return fmt.Errorf("error getting workout plan ID from circuit: %w", err)
    }

    filter := bson.M{
        "userId":     userID,
        "exerciseId": exerciseID,
        "circuitId":  circuitID,
        "workoutPlanId": workoutPlanID,
    }
    var userExerciseStatus models.UserExerciseStatus
    if err := us.database.Collection("userExerciseStatus").FindOne(ctx, filter).Decode(&userExerciseStatus); err != nil {
        if err == mongo.ErrNoDocuments {
            return ErrExerciseNotFound
        }
        return fmt.Errorf("error retrieving exercise status: %w", err)
    }
    
    if userExerciseStatus.Completed {
        return ErrExerciseAlreadyCompleted
    }
    
    update := bson.M{
        "$set": bson.M{"completed": true},
        "$push": bson.M{"completedLogs": bson.M{"$each": logs}},
    }

    if _, err := us.database.Collection("userExerciseStatus").UpdateOne(ctx, filter, update); err != nil {
        return fmt.Errorf("error updating exercise status: %w", err)
    }

    return us.checkAndUpdateCircuitStatus(ctx, userID, circuitID, workoutPlanID)
}

// func (us *UserService) GetWorkoutPlanProgress(ctx context.Context, userID, workoutPlanID primitive.ObjectID) (float64, error) {
// 	totalDays, err := us.database.Collection("userWorkoutDayStatus").CountDocuments(ctx, bson.M{"userId": userID, "workoutPlanId": workoutPlanID})
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return 0, ErrWorkoutPlanNotFound
// 		}
// 		return 0, fmt.Errorf("error counting workout days: %w", err)
// 	}

// 	completedDays, err := us.database.Collection("userWorkoutDayStatus").CountDocuments(ctx, bson.M{"userId": userID, "workoutPlanId": workoutPlanID, "completed": true})
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return 0, ErrWorkoutPlanNotFound
// 		}
		
// 		return 0, fmt.Errorf("error counting completed workout days: %w", err)
// 	}

// 	if totalDays == 0 {
// 		return 0, nil
// 	}

// 	progress := float64(completedDays) / float64(totalDays) * 100

// 	return progress, nil
// }

