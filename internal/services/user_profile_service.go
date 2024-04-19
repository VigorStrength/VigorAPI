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

func (us *UserService) GetUserProfile(ctx context.Context, userID primitive.ObjectID) (*models.UserProfile, error) {
	userCollection := us.database.Collection("users")
	var user models.User

	projection := bson.M{"profileInformation": 1}
	filter := bson.M{"_id": userID}
	opts := options.FindOne().SetProjection(projection)
	if err := userCollection.FindOne(ctx, filter, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("error returning user profile informations: %w", err)
	}

	return &user.ProfileInformation, nil
}

func (us *UserService) UpdateUserProfile(ctx context.Context, userID primitive.ObjectID, updateInput models.UserProfileUpdateInput) (*models.UserProfile, error) {
	userCollection := us.database.Collection("users")
	filter := bson.M{"_id": userID}

	var existingUser models.User
	if err := userCollection.FindOne(ctx, filter).Decode(&existingUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	updatedProfile := mergeUserProfileUpdates(existingUser.ProfileInformation, updateInput)
	updatedDoc := bson.M{"$set": bson.M{"profileInformation": updatedProfile}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.User
	if err := userCollection.FindOneAndUpdate(ctx, filter, updatedDoc, opts).Decode(&updatedUser); err != nil {
		return nil, fmt.Errorf("error updating user profile informations: %w", err)
	}

	return &updatedUser.ProfileInformation, nil
}

func mergeUserProfileUpdates(existingProfile models.UserProfile, updateInput models.UserProfileUpdateInput) models.UserProfile {
	if updateInput.Username != nil {
		existingProfile.Username = *updateInput.Username
	}
	if updateInput.ProfilePicture != nil {
		existingProfile.ProfilePicture = *updateInput.ProfilePicture
	}
	if updateInput.MainGoal != nil {
		existingProfile.MainGoal = *updateInput.MainGoal
	}
	if updateInput.SecondaryGoal != nil {
		existingProfile.SecondaryGoal = updateInput.SecondaryGoal
	}

	// Merge nested BodyInformation
	if updateInput.BodyInformation != nil {
		existingProfile = mergeBodyInformationUpdates(existingProfile, *updateInput.BodyInformation)
	}

	// Merge PhysicalActivity 
	if updateInput.PhysicalActivity != nil {
		existingProfile = mergePhysicalActivityUpdates(existingProfile, *updateInput.PhysicalActivity)
	}

	// Merge Lifestyle
	if updateInput.Lifestyle != nil {
		existingProfile = mergeLifestyleUpdates(existingProfile, *updateInput.Lifestyle)
	}

	// Handle CycleInformation if applicable
	if updateInput.CycleInformation != nil {
		if existingProfile.CycleInformation == nil {
			existingProfile.CycleInformation = &models.CycleInformation{}
		} 
		if updateInput.CycleInformation.ReproductiveStage != nil {
			existingProfile.CycleInformation.ReproductiveStage = *updateInput.CycleInformation.ReproductiveStage
		}
	}

	return existingProfile
}

func mergeBodyInformationUpdates(existingProfile models.UserProfile, updateInput models.BodyInformationUpdateInput) models.UserProfile {
	if updateInput.BodyType != nil {
		existingProfile.BodyInformation.BodyType = *updateInput.BodyType
	}
	if updateInput.BodyGoal != nil {
		existingProfile.BodyInformation.BodyGoal = updateInput.BodyGoal
	}
	if updateInput.HealthRestrictions != nil {
		existingProfile.BodyInformation.HealthRestrictions = *updateInput.HealthRestrictions
	}
	if updateInput.FocusArea != nil {
		existingProfile.BodyInformation.FocusArea = *updateInput.FocusArea
	}

	return existingProfile
}

func mergePhysicalActivityUpdates(existingProfile models.UserProfile, updateInput models.PhysicalActivityUpdateInput) models.UserProfile {
	if updateInput.FitnessLevel != nil {
		existingProfile.PhysicalActivity.FitnessLevel = *updateInput.FitnessLevel
	}
	if updateInput.Activities != nil {
		existingProfile.PhysicalActivity.Activities = *updateInput.Activities
	}

	return existingProfile
}

func mergeLifestyleUpdates(existingProfile models.UserProfile, updateInput models.LifestyleUpdateInput) models.UserProfile {
	if updateInput.Diet != nil {
		existingProfile.Lifestyle.Diet = *updateInput.Diet
	}
	if updateInput.WaterIntake != nil {
		existingProfile.Lifestyle.WaterIntake = updateInput.WaterIntake
	}
	if updateInput.SleepDuration != nil {
		existingProfile.Lifestyle.SleepDuration = updateInput.SleepDuration
	}
	if updateInput.TypicalDay != nil {
		existingProfile.Lifestyle.TypicalDay = *updateInput.TypicalDay
	}
	if updateInput.TrainingLocation != nil {
		existingProfile.Lifestyle.TrainingLocation = *updateInput.TrainingLocation
	}
	if updateInput.WorkoutTime != nil {
		existingProfile.Lifestyle.WorkoutTime = *updateInput.WorkoutTime
	}
	if updateInput.WorkoutFrequency != nil {
		existingProfile.Lifestyle.WorkoutFrequency = updateInput.WorkoutFrequency
	}
	if updateInput.WorkoutDuration != nil {
		existingProfile.Lifestyle.WorkoutDuration = *updateInput.WorkoutDuration
	}
	if updateInput.DiscoveryMethod != nil {
		existingProfile.Lifestyle.DiscoveryMethod = updateInput.DiscoveryMethod
	}
	if updateInput.IntolerancesAndAllergies != nil {
		existingProfile.Lifestyle.IntolerancesAndAllergies = *updateInput.IntolerancesAndAllergies
	}

	return existingProfile
}