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

func (us *UserService) GetUserPreferences(ctx context.Context, userID primitive.ObjectID) (*models.SystemPreferences, error) {
	userCollection := us.database.Collection("users")
	var user models.User

	projection := bson.M{"preferences": 1}
	filter := bson.M{"_id": userID}
	opts := options.FindOne().SetProjection(projection)
	if err := userCollection.FindOne(ctx, filter, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("error returning user preferences: %w", err)
	}

	return user.SystemPreferences, nil
}

func (us *UserService) UpdateUserSystemPreferences(ctx context.Context, userID primitive.ObjectID, updateInput models.SystemPreferencesUpdateInput) (*models.SystemPreferences, error) {
	userCollection := us.database.Collection("users")
	filter := bson.M{"_id": userID}

	var existingUser models.User
	if err := userCollection.FindOne(ctx, filter).Decode(&existingUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	updatedPreferences := mergeUserSystemPrefencesUpdates(*existingUser.SystemPreferences, updateInput)
	updatedDoc := bson.M{"$set": bson.M{"preferences": updatedPreferences}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.User
	if err := userCollection.FindOneAndUpdate(ctx, filter, updatedDoc, opts).Decode(&updatedUser); err != nil {
		return nil, fmt.Errorf("error updating user preferences: %w", err)
	}

	return updatedUser.SystemPreferences, nil
}

func mergeUserSystemPrefencesUpdates(existingPreferences models.SystemPreferences, updateInput models.SystemPreferencesUpdateInput) models.SystemPreferences {
	if updateInput.Language != nil {
		existingPreferences.Language = *updateInput.Language
	}
	if updateInput.TimeZone != nil {
		existingPreferences.TimeZone = *updateInput.TimeZone
	}
	if updateInput.DisplayMode != nil {
		existingPreferences.DisplayMode = *updateInput.DisplayMode
	}
	if updateInput.MeasurementSystem != nil {
		existingPreferences.MeasurementSystem = *updateInput.MeasurementSystem
	}
	if updateInput.AllowReadReceipt != nil {
		existingPreferences.AllowReadReceipt = *updateInput.AllowReadReceipt
	}

	return existingPreferences
}