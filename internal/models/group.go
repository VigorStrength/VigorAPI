package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group represents a collection of users, such as a workout group or a meal plan group who can communicate with each others.
type Group struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string               `bson:"name" json:"name" binding:"required"`
	Members   []primitive.ObjectID `bson:"members" json:"members" binding:"required"` // IDs of User documents
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt" binding:"required"`
	CreatedBy primitive.ObjectID   `bson:"createdBy" json:"createdBy" binding:"required"` // ID of the User who created the group
	// Optionally, you can include additional fields specific to the group's purpose.
	// For instance, if it's a workout group, you might include fields for common goals, scheduled sessions, etc.
}
