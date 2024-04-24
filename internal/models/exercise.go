package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Exercise represents a specific exercise, including details for logging and video interaction.
type Exercise struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string             `bson:"name" json:"name" binding:"required" validate:"required,min=5,max=50"`
	Description     string             `bson:"description" json:"description" binding:"required" validate:"required,min=5,max=1000"`
	VideoURL        string             `bson:"videoURL" json:"videoURL" binding:"required" validate:"required,url"`
	TargetMuscles   []string           `bson:"targetMuscles" json:"targetMuscles" binding:"required" validate:"required,dive,required,gt=0"`
	EquipmentNeeded []string           `bson:"equipmentNeeded,omitempty" json:"equipmentNeeded,omitempty" validate:"omitempty,dive,required"`
	Instructions    []string           `bson:"instructions" json:"instructions" binding:"required" validate:"required,dive,required"`
	Time            int                `bson:"time" json:"time" binding:"required" validate:"required,gte=30,lte=300"` // Time in seconds.
	ProposedLog     ExerciseLog        `bson:"proposedLog" json:"proposedLog" validate:"required"`
}

// ExerciseLog represents the logging of exercises, including proposed logs and actual user logs.
type ExerciseLog struct {
	SetNumber      *int     `bson:"setNumber,omitempty" json:"setNumber,omitempty" validate:"omitempty"`
	ProposedReps   int      `bson:"proposedReps" json:"proposedReps" binding:"required" validate:"required,gt=0"`
	ActualReps     *int     `bson:"actualReps,omitempty" json:"actualReps,omitempty" validate:"omitempty"`
	ProposedWeight *float64 `bson:"proposedWeight,omitempty" json:"proposedWeight,omitempty" validate:"omitempty"` // Can be empty if no equipment.
	ActualWeight   *float64 `bson:"actualWeight,omitempty" json:"actualWeight,omitempty" validate:"omitempty"`     // Can be empty if no equipment.
	// ... Add other fields like duration, rate of perceived exertion, etc.
}

// ExerciseUpdateInput represents the input for updating an exercise.
type ExerciseUpdateInput struct {
	Name            *string                 `json:"name" validate:"omitempty,min=5,max=50"`
	Description     *string                 `json:"description" validate:"omitempty,min=5,max=1000"`
	VideoURL        *string                 `json:"videoURL" validate:"omitempty,url"`
	TargetMuscles   *[]string               `json:"targetMuscles" validate:"omitempty,dive,required,gt=0"`
	EquipmentNeeded *[]string               `json:"equipmentNeeded" validate:"omitempty,dive,required"`
	Instructions    *[]string               `json:"instructions" validate:"omitempty,dive,required"`
	Time            *int                    `json:"time" validate:"omitempty,gte=30,lte=300"`
	ProposedLog     *ExerciseLogUpdateInput `json:"proposedLog" validate:"omitempty"`
}

// ExerciseLogUpdateInput represents the input for updating an exercise log.
type ExerciseLogUpdateInput struct {
	SetNumber      *int     `json:"setNumber" validate:"omitempty"`
	ProposedReps   *int     `json:"proposedReps" validate:"omitempty,gt=0"`
	ActualReps     *int     `json:"actualReps" validate:"omitempty"`
	ProposedWeight *float64 `json:"proposedWeight" validate:"omitempty"`
	ActualWeight   *float64 `json:"actualWeight" validate:"omitempty"`
}
