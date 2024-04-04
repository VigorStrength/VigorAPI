package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Exercise represents a specific exercise, including details for logging and video interaction.
type Exercise struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string             `bson:"name" json:"name" binding:"required" validate:"required,alpha,min=5,max=25"`
	Description     string             `bson:"description" json:"description" binding:"required" validate:"required,alpha,min=5,max=100"`
	VideoURL        string             `bson:"videoUrl" json:"videoUrl" binding:"required" validate:"required,url"`
	TargetMuscles   []string           `bson:"targetMuscles" json:"targetMuscles" binding:"required" validate:"required,dive,required"`
	EquipmentNeeded []string           `bson:"equipmentNeeded" json:"equipmentNeeded" binding:"required" validate:"required,dive,required"`
	Instructions    []string           `bson:"instructions" json:"instructions" binding:"required" validate:"required,dive,required"`
	Time            int                `bson:"time" json:"time" binding:"required" validate:"required"` // Time in seconds.
	Log             ExerciseLog      `bson:"log" json:"log" validate:"required,dive,required"`
}

type UserExerciseStatus struct {
	UserID        primitive.ObjectID `bson:"userId" json:"userId" binding:"required" validate:"required"`
	ExerciseID    primitive.ObjectID `bson:"exerciseId" json:"exerciseId" binding:"required" validate:"required"`
	Completed     bool               `bson:"completed" json:"completed" binding:"required" validate:"omitempty"`
	CompletedLogs []ExerciseLog      `bson:"completedLogs" json:"completedLogs" binding:"required" validate:"required,dive,required"`
}

// ExerciseLog represents the logging of exercises, including proposed logs and actual user logs.
type ExerciseLog struct {
	SetNumber      *int      `bson:"setNumber,omitempty" json:"setNumber,omitempty" validate:"omitempty"`
	ProposedReps   int      `bson:"proposedReps" json:"proposedReps" binding:"required" validate:"required"`
	ActualReps     *int      `bson:"actualReps,omitempty" json:"actualReps,omitempty" validate:"omitempty"`
	ProposedWeight *float64 `bson:"proposedWeight,omitempty" json:"proposedWeight,omitempty" validate:"omitempty"` // Can be empty if no equipment.
	ActualWeight   *float64 `bson:"actualWeight,omitempty" json:"actualWeight,omitempty" validate:"omitempty"`     // Can be empty if no equipment.
	// ... Add other fields like duration, rate of perceived exertion, etc.
}

