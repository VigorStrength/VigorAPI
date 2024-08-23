package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrSupersetAlreadyExists = fmt.Errorf("superset already exists")
	ErrSupersetNotFound = fmt.Errorf("superset not found")
)

func (as * AdminService) CreateSuperset(ctx context.Context, superset models.Superset) ( error) {
	// Get the superset collection
	supersetCollection := as.database.Collection("supersets")

	// check if the new superset already exists
	filter := bson.M{"_id": superset.ID}
	count, err := supersetCollection.CountDocuments(ctx, filter)
	if err != nil {
		return  fmt.Errorf("error checking if superset already exists: %w", err)
	}

	if count > 0 {
		return  ErrSupersetAlreadyExists
	}

	// Insert the new superset
	_, err = supersetCollection.InsertOne(ctx, superset)
	if err != nil {
		return  fmt.Errorf("error inserting new superset: %w", err)
	}

	return  nil
}

func (as *AdminService) UpdateSuperset(ctx context.Context, supersetID primitive.ObjectID, supersetInput models.SupersetUpdateInput) (error) {
	// Get the superset collection
	supersetCollection := as.database.Collection("supersets")

	// Find the superset by ID
	filter := bson.M{"_id": supersetID}
	var superset models.Superset
	err := supersetCollection.FindOne(ctx, filter).Decode(&superset)
	if err != nil {
		return  fmt.Errorf("error finding superset: %w", err)
	}

	// Update the superset
	update := bson.M{"$set": supersetInput}
	_, err = supersetCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating superset: %w", err)
	}

	return nil
}

func (as *AdminService) GetSupersets(ctx context.Context) ([]models.Superset, error) {
	// Get the superset collection
	supersetCollection := as.database.Collection("supersets")

	// Find all supersets
	cursor, err := supersetCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding supersets: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the supersets
	var supersets []models.Superset
	if err := cursor.All(ctx, &supersets); err != nil {
		return nil, fmt.Errorf("error decoding supersets: %w", err)
	}

	return supersets, nil
}