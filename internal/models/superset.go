package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Superset represents a set of exercises performed in sequence, with no rest in between.
type Superset struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	ExerciseIDs []primitive.ObjectID `bson:"exerciseIds" json:"exerciseIds" binding:"required" validate:"required,dive,required"`
	RestTime    *int                 `bson:"restTime,omitempty" json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Rest time in seconds after the superset.
	ProposedLaps int                  `bson:"proposedLaps" json:"proposedLaps" binding:"required" validate:"required,gte=1"`
}

type SupersetUpdateInput struct {
	ExerciseIDs  *[]primitive.ObjectID `json:"exerciseIds,omitempty" validate:"omitempty,dive,required"`
	RestTime     *int                 `json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Rest time in seconds after the superset.
	ProposedLaps *int                  `json:"proposedLaps,omitempty" validate:"omitempty,gte=1"`
}

