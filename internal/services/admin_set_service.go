package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrSetAlreadyExists = fmt.Errorf("set already exists")
	ErrSetNotFound 	= fmt.Errorf("set not found")
)

func (as *AdminService) CreateSet(ctx context.Context, set models.Set) error {
	setCollection := as.database.Collection("sets")

	filter := bson.M{"_id": set.ID}
	count, err := setCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking if set already exists: %w", err)
	}

	if count > 0 {
		return ErrSetAlreadyExists
	}

	_, err = setCollection.InsertOne(ctx, set)
	if err != nil {
		return fmt.Errorf("error inserting new set: %w", err)
	}

	return nil
}

func (as *AdminService) UpdateSet(ctx context.Context, setID primitive.ObjectID, setInput models.SetUpdateInput) error {
	// Check if setInput is empty
	if (models.SetUpdateInput{}) == setInput {
		return fmt.Errorf("setInput cannot be empty")
	}

	setCollection := as.database.Collection("sets")

	filter := bson.M{"_id": setID}
	var set models.Set
	err := setCollection.FindOne(ctx, filter).Decode(&set)
	if err != nil {
		return fmt.Errorf("error finding set: %w", err)
	}

	update := bson.M{"$set": setInput}
	_, err = setCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating set: %w", err)
	}

	return nil
}

func (as *AdminService) GetSets(ctx context.Context) ([]models.Set, error) {
	setCollection := as.database.Collection("sets")

	cursor, err := setCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting sets: %w", err)
	}

	var sets []models.Set
	if err := cursor.All(ctx, &sets); err != nil {
		return nil, fmt.Errorf("error decoding sets: %w", err)
	}

	return sets, nil
}