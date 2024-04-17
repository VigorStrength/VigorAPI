package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (as *AdminService) GetUsers(ctx context.Context) ([]models.User, error) {
	userCollection := as.database.Collection("users")

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error decoding users: %w", err)
	}

	return users, nil
}