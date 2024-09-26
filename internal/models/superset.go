package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Superset represents a set of exercises performed in sequence, with no rest in between.
type Superset struct {
	Circuit    					 	 `bson:",inline"`
	RestTime    *int                 `bson:"restTime,omitempty" json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Rest time in seconds after the superset.
	ProposedLaps int                  `bson:"proposedLaps" json:"proposedLaps" binding:"required" validate:"required,gte=1"`
}

type SupersetUpdateInput struct {
	ExerciseIDs  *[]primitive.ObjectID `bson:"exerciseIds,omitempty" json:"exerciseIds,omitempty" validate:"omitempty,dive,required"`
	RestTime     *int                 `bson:"restTime,omitempty" json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Rest time in seconds after the superset.
	ProposedLaps *int                  `bson:"proposedLaps,omitempty" json:"proposedLaps,omitempty" validate:"omitempty,gte=1"`
}

