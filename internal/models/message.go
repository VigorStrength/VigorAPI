package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a message sent from one user to another or to a group.
type Message struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	SenderID    primitive.ObjectID  `bson:"senderId" json:"senderId" binding:"required"`
	RecipientID *primitive.ObjectID `bson:"recipientId,omitempty" json:"recipientId,omitempty"`
	GroupID     *primitive.ObjectID `bson:"groupId,omitempty" json:"groupId,omitempty"`
	Content     string              `bson:"content" json:"content" binding:"required"`
	SentAt      time.Time           `bson:"sentAt" json:"sentAt" binding:"required"`
	ReceivedAt  *time.Time          `bson:"receivedAt,omitempty" json:"receivedAt,omitempty"`
	ReadAt      *time.Time          `bson:"readAt,omitempty" json:"readAt,omitempty"`
	// Note: AllowReadReceipt is now handled at the user's system preferences level.
}
