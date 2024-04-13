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
)

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