package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Conversation represents a chat session between two users or a group.
type Conversation struct {
    ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
    Participants []primitive.ObjectID `bson:"participants" json:"participants" binding:"required"` // Array of UserIDs
    GroupID      *primitive.ObjectID  `bson:"groupId,omitempty" json:"groupId,omitempty"` // Optional, present if this is a group conversation
    CreatedAt    time.Time            `bson:"createdAt" json:"createdAt" binding:"required"`
    UpdatedAt    time.Time            `bson:"updatedAt" json:"updatedAt" binding:"required"`
}

// Message represents a message within a conversation.
type Message struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId" binding:"required"` // Link to Conversation
    SenderID       primitive.ObjectID `bson:"senderId" json:"senderId" binding:"required"`
    Content        string             `bson:"content" json:"content" binding:"required"`
    SentAt         time.Time          `bson:"sentAt" json:"sentAt" binding:"required"`
    ReceivedAt     *time.Time         `bson:"receivedAt,omitempty" json:"receivedAt,omitempty"`
    ReadAt         *time.Time         `bson:"readAt,omitempty" json:"readAt,omitempty"`
}