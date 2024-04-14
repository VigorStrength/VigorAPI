package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrMealPlanAlreadyExists = fmt.Errorf("meal plan already exists")
	ErrMealPlanNotFound = fmt.Errorf("meal plan not found")
)

func (as *AdminService) GetMealPlanByID(ctx context.Context, mealPlanID primitive.ObjectID) (models.MealPlan, error) {
	mealPlanCollection := as.database.Collection("mealPlans")

	filter := bson.M{"_id": mealPlanID}
	var mealPlan models.MealPlan
	err := mealPlanCollection.FindOne(ctx, filter).Decode(&mealPlan)
	if err != nil {
		return models.MealPlan{}, fmt.Errorf("error finding meal plan: %w", err)
	}

	return mealPlan, nil
}

func (as *AdminService) GetMealPlans(ctx context.Context) ([]models.MealPlan, error) {
	mealPlanCollection := as.database.Collection("mealPlans")

	cursor, err := mealPlanCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding meal plans: %w", err)
	}
	defer cursor.Close(ctx)

	var mealPlans []models.MealPlan
	if err := cursor.All(ctx, &mealPlans); err != nil {
		return nil, fmt.Errorf("error decoding meal plans: %w", err)
	}

	return mealPlans, nil
}

func (as *AdminService) SearchMealPlansByName(ctx context.Context, name string) ([]models.MealPlan, error) {
	mealPlanCollection := as.database.Collection("mealPlans")

	filter := bson.M{"name": primitive.Regex{Pattern: name, Options: "i"}}
	cursor, err := mealPlanCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding meal plans by name: %w", err)
	}
	defer cursor.Close(ctx)

	var mealPlans []models.MealPlan
	if err := cursor.All(ctx, &mealPlans); err != nil {
		return nil, fmt.Errorf("error decoding meal plans: %w", err)
	}

	return mealPlans, nil
}

func (as *AdminService) CreateMealPlan(ctx context.Context, mealPlanInput models.MealPlan) error {
	mealPlanCollection := as.database.Collection("mealPlans")

	filter := bson.M{"name": mealPlanInput.Name}
	count, err := mealPlanCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error counting meal plans: %w", err)
	}

	if count > 0 {
		return ErrMealPlanAlreadyExists
	}

	mealPlanInput.ID = primitive.NewObjectID()
	for i := range mealPlanInput.WeeklyPlans {
		mealPlanInput.WeeklyPlans[i].ID = primitive.NewObjectID()
		for j := range mealPlanInput.WeeklyPlans[i].DailyPlans {
			mealPlanInput.WeeklyPlans[i].DailyPlans[j].ID = primitive.NewObjectID()
		}
	}

	_, err = mealPlanCollection.InsertOne(ctx, mealPlanInput)
	if err != nil {
		return fmt.Errorf("error inserting meal plan: %w", err)
	}

	return nil
}

func (as *AdminService) DeleteMealPlan(ctx context.Context, mealPlanID primitive.ObjectID) error {
	mealPlanCollection := as.database.Collection("mealPlans")

	filter := bson.M{"_id": mealPlanID}
	result, err := mealPlanCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting meal plan: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrMealPlanNotFound
	}

	return nil
}