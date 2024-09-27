package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StandAlone struct {
	ID         	primitive.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
	ExerciseID  primitive.ObjectID    `bson:"exerciseId" json:"exerciseId" binding:"required" validate:"required"`
	RestTime   int                  `bson:"restTime" json:"restTime"  binding:"required" validate:"required,gte=5,lte=240"` // Rest time in seconds after the superset.
	ProposedLaps int                  `bson:"proposedLaps" json:"proposedLaps" binding:"required" validate:"required,gte=1"`
}

type StandAloneUpdateInput struct {
	ExerciseID *primitive.ObjectID `bson:"exerciseId,omitempty" json:"exerciseId,omitempty" validate:"omitempty"`
	RestTime   *int                 `bson:"restTime,omitempty" json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Rest time in seconds after the superset.
	ProposedLaps *int                  `bson:"proposedLaps,omitempty" json:"proposedLaps,omitempty" validate:"omitempty,gte=1"`
}