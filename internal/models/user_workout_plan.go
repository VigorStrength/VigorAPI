package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserExerciseStatus struct {
	UserID        primitive.ObjectID `bson:"userId" json:"userId" binding:"required" validate:"required"`
	ExerciseID    primitive.ObjectID `bson:"exerciseId" json:"exerciseId" binding:"required" validate:"required"`
	CircuitID     primitive.ObjectID `bson:"circuitId" json:"circuitId" binding:"required" validate:"required"`
	WorkoutPlanID primitive.ObjectID `bson:"workoutPlanId" json:"workoutPlanId" binding:"required" validate:"required"`
	Completed     bool               `bson:"completed" json:"completed" binding:"required" validate:"omitempty"`
	CompletedLogs []ExerciseLog      `bson:"completedLogs" json:"completedLogs" binding:"required" validate:"required,dive,required"`
}

// UserCircuitStatus tracks the completion status of a circuit for a specific user.
type UserCircuitStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	CircuitID primitive.ObjectID `bson:"circuitId" json:"circuitId" binding:"required"`
	WorkoutDayID primitive.ObjectID `bson:"workoutDayId" json:"workoutDayId" binding:"required"` // Reference the workout day
	WorkoutPlanID  primitive.ObjectID `bson:"workoutPlanId" json:"workoutPlanId" binding:"required"` // Reference to the WorkoutPlan
	Completed bool               `bson:"completed" json:"completed"`
}

// UserWorkoutDayStatus tracks the completion status of a workout day for a specific user.
type UserWorkoutDayStatus struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	WorkoutDayID primitive.ObjectID `bson:"workoutDayId" json:"workoutDayId" binding:"required"`
	WorkoutWeekID primitive.ObjectID `bson:"workoutWeekId" json:"workoutWeekId" binding:"required"` // Reference the workout week
	WorkoutPlanID  primitive.ObjectID `bson:"workoutPlanId" json:"workoutPlanId" binding:"required"` // Reference to the WorkoutPlan
	Completed    bool               `bson:"completed" json:"completed"`
}

// UserWorkoutWeekStatus tracks the completion status of a workout week for a specific user.
type UserWorkoutWeekStatus struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	WorkoutWeekID primitive.ObjectID `bson:"workoutWeekId" json:"workoutWeekId" binding:"required"`
	WorkoutPlanID  primitive.ObjectID `bson:"workoutPlanId" json:"workoutPlanId" binding:"required"` // Reference to the WorkoutPlan
	CompletedDays int                `bson:"completedDays" json:"completedDays"`
}

type UserWorkoutPlanStatus struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`               // Reference to the User
	WorkoutPlanID  primitive.ObjectID `bson:"workoutPlanId" json:"workoutPlanId" binding:"required"` // Reference to the WorkoutPlan
	StartDate      time.Time          `bson:"startDate" json:"startDate" binding:"required"`
	CompletionDate *time.Time         `bson:"completionDate" json:"completionDate"` // nil if not completed
	Completed      bool               `bson:"completed" json:"completed"`
	// More fields as necessary to track progress, such as completed workouts or weeks
}