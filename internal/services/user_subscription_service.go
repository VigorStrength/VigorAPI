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

func (us *UserService) UpdateUserSubscription(ctx context.Context, userID primitive.ObjectID, updateInput models.UserSubscriptionUpdateInput) (*models.UserSubscription, error) {	
	userCollection := us.database.Collection("users")
	filter := bson.M{"_id": userID}

	var existingUser models.User
	if err := userCollection.FindOne(ctx, filter).Decode(&existingUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	updatedSubscription := mergeUserSubscriptionUpdates(existingUser.Subscription, updateInput)
	updatedDoc := bson.M{"$set": bson.M{"subscription": updatedSubscription}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.User
	if err := userCollection.FindOneAndUpdate(ctx, filter, updatedDoc, opts).Decode(&updatedUser); err != nil {
		return nil, fmt.Errorf("error updating user subscription informations: %w", err)
	}

	return &updatedUser.Subscription, nil
}

func mergeUserSubscriptionUpdates(existingSubscription models.UserSubscription, updateInput models.UserSubscriptionUpdateInput) models.UserSubscription {
	if updateInput.Type != nil {
		existingSubscription.Type = *updateInput.Type
	}
	if updateInput.Status != nil {
		existingSubscription.Status = *updateInput.Status
	}
	if updateInput.StartDate != nil {
		existingSubscription.StartDate = *updateInput.StartDate
	}
	if updateInput.EndDate != nil {
		existingSubscription.EndDate = updateInput.EndDate
	}
	if updateInput.IsActive != nil {
		existingSubscription.IsActive = *updateInput.IsActive
	}

	return existingSubscription
}