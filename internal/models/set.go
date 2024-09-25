package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Set struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	ExerciseIDs []primitive.ObjectID `bson:"exerciseIds" json:"exerciseIds" binding:"required" validate:"required,dive,required"`
}

type SetUpdateInput struct {
	ExerciseIDs  *[]primitive.ObjectID `json:"exerciseIds,omitempty" validate:"omitempty,dive,required"`
}