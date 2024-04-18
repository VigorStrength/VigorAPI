package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (us *UserService) GetUserProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	userCollection := us.database.Collection("users")
	var user models.User

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("error converting user ID to ObjectID: %w", err)
	}

	projection := bson.M{"profileInformation": 1}
	filter := bson.M{"_id": objID}
	opts := options.FindOne().SetProjection(projection)
	if err := userCollection.FindOne(ctx, filter, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("error returning user profile informations: %w", err)
	}

	return &user.ProfileInformation, nil
}