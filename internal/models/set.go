package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Set struct {
	// Superset `bson:",inline"` Use this in a different branch on code cleanup.
	Circuit     					 `bson:",inline"`
	RestTime    *int                 `bson:"restTime,omitempty" json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Rest time in seconds after the superset.
	ProposedLaps int                  `bson:"proposedLaps" json:"proposedLaps" binding:"required" validate:"required,gte=1"`
}

type SetUpdateInput struct {
	ExerciseIDs  *[]primitive.ObjectID `bson:"exerciseIds,omitempty" json:"exerciseIds,omitempty" validate:"omitempty,dive,required"`
	RestTime     *int                 `bson:"restTime,omitempty" json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Rest time in seconds after the superset.
	ProposedLaps *int                  `bson:"proposedLaps,omitempty" json:"proposedLaps,omitempty" validate:"omitempty,gte=1"`
}