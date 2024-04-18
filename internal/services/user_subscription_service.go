package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (us *UserService) GetUserSubsctiption(ctx context.Context, userID primitive.ObjectID) (*models.UserSubscription, error) {
	userCollection := us.database.Collection("users")
	var user models.User

	projection := bson.M{"subscription": 1}
	filter := bson.M{"_id": userID}
	opts := options.FindOne().SetProjection(projection)
	if err := userCollection.FindOne(ctx, filter, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("error returning user subscription informations: %w", err)
	}

	return &user.Subscription, nil
}