package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrMealAlreadyExists = fmt.Errorf("meal already exists")
	ErrMealNotFound      = fmt.Errorf("meal not found")
)

func (as *AdminService) GetMealByID(ctx context.Context, mealID primitive.ObjectID) (models.Meal, error) {
	//Get the meal collection
	mealCollection := as.database.Collection("meals")

	// Find the meal by ID
	filter := bson.M{"_id": mealID}
	var meal models.Meal
	err := mealCollection.FindOne(ctx, filter).Decode(&meal)
	if err != nil {
		return models.Meal{}, fmt.Errorf("error finding meal: %w", err)
	}

	return meal, nil
}

func (as *AdminService) GetMeals(ctx context.Context) ([]models.Meal, error) {
	mealCollection := as.database.Collection("meals")

	cursor, err := mealCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding meals: %w", err)
	}
	defer cursor.Close(ctx)

	var meals []models.Meal
	if err := cursor.All(ctx, &meals); err != nil {
		return nil, fmt.Errorf("error decoding meals: %w", err)
	}

	return meals, nil
}

func (as *AdminService) CreateMeal(ctx context.Context, meal models.Meal) error {
	//Get the meal collection
	mealCollection := as.database.Collection("meals")

	// check if the new meal already exists
	filter := bson.M{"name": meal.Name}
	count, err := mealCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking if meal already exists: %w", err)
	}

	if count > 0 {
		return ErrMealAlreadyExists
	}

	// Insert the new meal
	_, err = mealCollection.InsertOne(ctx, meal)
	if err != nil {
		return fmt.Errorf("error inserting meal: %w", err)
	}

	return nil
}

func (as *AdminService) DeleteMeal(ctx context.Context, mealID primitive.ObjectID) error {
	mealCollection := as.database.Collection("meals")

	filter := bson.M{"_id": mealID}
	result, err := mealCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting meal: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrMealNotFound
	}

	return nil
}

func (as *AdminService) SearchMealsByName(ctx context.Context, name string) ([]models.Meal, error) {
	mealCollection := as.database.Collection("meals")

	filter := bson.M{"name": primitive.Regex{Pattern: name, Options: "i"}}
	cursor, err := mealCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding meals: %w", err)
	}
	defer cursor.Close(ctx)

	var meals []models.Meal
	if err := cursor.All(ctx, &meals); err != nil {
		return nil, fmt.Errorf("error decoding meals: %w", err)
	}

	return meals, nil
}

