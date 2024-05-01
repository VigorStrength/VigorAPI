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

func (us *UserService) getAndValidateCircuit(ctx context.Context, circuitID primitive.ObjectID) (primitive.ObjectID, error) {
	var userCircuitStatus models.UserCircuitStatus
	err := us.database.Collection("userCircuitStatus").FindOne(ctx, bson.M{"circuitId": circuitID}).Decode(&userCircuitStatus)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return userCircuitStatus.WorkoutPlanID, nil
}

func (us *UserService) getWorkoutDayID(ctx context.Context, circuitID primitive.ObjectID) (primitive.ObjectID, error) {
	var userCircuit models.UserCircuitStatus
	err := us.database.Collection("userCircuitStatus").FindOne(ctx, bson.M{"circuitId": circuitID}).Decode(&userCircuit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return primitive.NilObjectID, ErrCircuitNotFound
		}
		return primitive.NilObjectID, fmt.Errorf("error finding user's circuit: %w", err)
	}
	return userCircuit.WorkoutDayID, nil
}

func (us *UserService) getWeekID(ctx context.Context, workoutDayID primitive.ObjectID) (primitive.ObjectID, error) {
	var userWorkoutDay models.UserWorkoutDayStatus
	err := us.database.Collection("userWorkoutDayStatus").FindOne(ctx, bson.M{"workoutDayId": workoutDayID}).Decode(&userWorkoutDay)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error finding user's workout day: %w", err)
	}
	return userWorkoutDay.WorkoutWeekID, nil
}

func (us *UserService) incrementWorkoutWeekCompletedDays(ctx context.Context, userID, workoutWeekID, workoutPlanID primitive.ObjectID) error {
	filter := bson.M{"userId": userID, "workoutWeekId": workoutWeekID, "workoutPlanId": workoutPlanID}
	update := bson.M{"$inc": bson.M{"completedDays": 1}}
	if _, err := us.database.Collection("userWorkoutWeekStatus").UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (us *UserService) checkAndUpdateCircuitStatus(ctx context.Context, userID, circuitID, workoutPlanID primitive.ObjectID) error {
	filter := bson.M{"userId": userID, "circuitId": circuitID, "workoutPlanId": workoutPlanID, "completed": false}
	count, err := us.database.Collection("userExerciseStatus").CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking circuit completion: %w", err)
	}

	if count == 0 {
		filter := bson.M{"userId": userID, "circuitId": circuitID, "workoutPlanId": workoutPlanID}
		update := bson.M{"$set": bson.M{"completed": true}}
		if _, err := us.database.Collection("userCircuitStatus").UpdateOne(ctx, filter, update); err != nil {
			return fmt.Errorf("error updating circuit status: %w", err)
		}

		workoutDayID, err := us.getWorkoutDayID(ctx, circuitID)
		if err != nil {
			return fmt.Errorf("error getting workoutDayID form circuit: %w", err)
		}
		return us.checkAndUpdateDayStatus(ctx, userID, workoutDayID, workoutPlanID)
	}

	return nil
}

func (us *UserService) checkAndUpdateDayStatus(ctx context.Context, userID, workoutDayID, workoutPlanID primitive.ObjectID) error {
	filter := bson.M{"userId": userID, "workoutDayId": workoutDayID, "workoutPlanId": workoutPlanID, "completed": false}
	count, err := us.database.Collection("userCircuitStatus").CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking day completion: %w", err)
	}

	if count == 0 {
		filter := bson.M{"userId": userID, "workoutDayId": workoutDayID, "workoutPlanId": workoutPlanID}
		update := bson.M{"$set": bson.M{"completed": true}}
		if _, err := us.database.Collection("userWorkoutDayStatus").UpdateOne(ctx, filter, update); err != nil {
			return fmt.Errorf("error updating day status: %w", err)
		}

		weekID, err := us.getWeekID(ctx, workoutDayID)
		if err != nil {
			return fmt.Errorf("error getting weekID from day: %w", err)
		}

		if err := us.incrementWorkoutWeekCompletedDays(ctx, userID, weekID, workoutPlanID); err != nil {
			if err == mongo.ErrNoDocuments {
				return ErrWrokoutWeekNotFound
			}

			return fmt.Errorf("error incrementing workout week completed days: %w", err)
		}

		return us.checkAndUpdateWeekStatus(ctx, userID, weekID, workoutPlanID)
	}

	return nil
}

func (us *UserService) checkAndUpdateWeekStatus(ctx context.Context, userID, workoutWeekID, workoutPlanID primitive.ObjectID) error {
	filter := bson.M{"userId": userID, "workoutWeekId": workoutWeekID, "workoutPlanId": workoutPlanID, "completed": false}
	count, err := us.database.Collection("userWorkoutDayStatus").CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking week completion: %w", err)
	}

	if count == 0 {
		filter := bson.M{"userId": userID, "workoutWeekId": workoutWeekID, "workoutPlanId": workoutPlanID}
		update := bson.M{"$set": bson.M{"completed": true}}
		if _, err := us.database.Collection("userWorkoutWeekStatus").UpdateOne(ctx, filter, update); err != nil {
			return fmt.Errorf("error updating week status: %w", err)
		}

		return us.checkAndUpdateWorkoutPlanStatus(ctx, userID, workoutPlanID)
	}

	return nil
}

func (us *UserService) checkAndUpdateWorkoutPlanStatus(ctx context.Context, userID, workoutPlanID primitive.ObjectID) error {
	filter := bson.M{"userId": userID, "workoutPlanId": workoutPlanID, "completed": false}
	count, err := us.database.Collection("userWorkoutWeekStatus").CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking workout plan completion: %w", err)
	}

	if count == 0 {
		filter := bson.M{"userId": userID, "workoutPlanId": workoutPlanID}
		update := bson.M{"$set": bson.M{"completionDate": time.Now() ,"completed": true, }}
		if _, err := us.database.Collection("userWorkoutPlanStatus").UpdateOne(ctx, filter, update); err != nil {
			return fmt.Errorf("error updating workout plan status: %w", err)
		}
	}

	return nil
}