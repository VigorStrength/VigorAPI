package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group represents a group of users who can communicate with each other.
type Group struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string               `bson:"name" json:"name" binding:"required"`
	Members   []primitive.ObjectID `bson:"members" json:"members" binding:"required"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt" binding:"required"`
	CreatedBy primitive.ObjectID   `bson:"createdBy" json:"createdBy" binding:"required"`
}
