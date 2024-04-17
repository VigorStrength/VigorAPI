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

func (as *AdminService) UpdateMealPlan(ctx context.Context, mealPlanID primitive.ObjectID, updateInput models.MealPlanUpdateInput) error {
	mealPlanCollection := as.database.Collection("mealPlans")
	filter := bson.M{"_id": mealPlanID}

	var existingMealPlan models.MealPlan
	if err := mealPlanCollection.FindOne(ctx, filter).Decode(&existingMealPlan); err != nil {
		return fmt.Errorf("error finding meal plan: %w", err)
	}

	updatedMealPlanDoc := mergeUpdatesIntoExistingMealPlan(existingMealPlan, updateInput)

	update := bson.M{"$set": updatedMealPlanDoc}
	result, err := mealPlanCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating meal plan: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrMealPlanNotFound
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

func mergeUpdatesIntoExistingMealPlan(existingPlan models.MealPlan, updateInput models.MealPlanUpdateInput) models.MealPlan {
	if updateInput.Name != nil {
		existingPlan.Name = *updateInput.Name
	}
	if updateInput.Duration != nil {
		existingPlan.Duration = *updateInput.Duration
	}

	existingWeeklyPlansMap := make(map[int]*models.WeeklyPlan)
	for i := range existingPlan.WeeklyPlans {
		weekNumber := existingPlan.WeeklyPlans[i].WeekNumber
		existingWeeklyPlansMap[weekNumber] = &existingPlan.WeeklyPlans[i]
	}

	updatedWeeklyPlans := []models.WeeklyPlan{}
	if updateInput.WeeklyPlans != nil {
		for _, updatedWeekInput := range *updateInput.WeeklyPlans {
			if existingWeek, exists := existingWeeklyPlansMap[*updatedWeekInput.WeekNumber]; exists {
				updatedWeek := mergeWeeklyPlanUpdates(*existingWeek, updatedWeekInput)
				updatedWeeklyPlans = append(updatedWeeklyPlans, updatedWeek)
				delete(existingWeeklyPlansMap, *updatedWeekInput.WeekNumber)
			} else {
				newWeeklyPlan := createNewWeeklyPlanFromInput(updatedWeekInput)
				updatedWeeklyPlans = append(updatedWeeklyPlans, newWeeklyPlan)
			}
		}
	}

	for _, existingWeek := range existingWeeklyPlansMap {
		updatedWeeklyPlans = append(updatedWeeklyPlans, *existingWeek)
	}

	existingPlan.WeeklyPlans = updatedWeeklyPlans

	return existingPlan
}

func mergeWeeklyPlanUpdates(existingWeeklyPlan models.WeeklyPlan, updateWeeklyInput models.WeeklyPlanUpdateInput) models.WeeklyPlan {
	if updateWeeklyInput.DailyPlans != nil {
		updatedDailyPlans := []models.DailyPlan{}
		for _, updatedDailyPlanInput := range *updateWeeklyInput.DailyPlans {
			newDailyPlan := createNewDailyPlanFromInput(updatedDailyPlanInput)
			updatedDailyPlans = append(updatedDailyPlans, newDailyPlan)
		}
		existingWeeklyPlan.DailyPlans = updatedDailyPlans
	}

	return existingWeeklyPlan
}

func createNewWeeklyPlanFromInput(input models.WeeklyPlanUpdateInput) models.WeeklyPlan {
	newWeeklyPlan := models.WeeklyPlan{
		ID:         primitive.NewObjectID(),
		WeekNumber: *input.WeekNumber,
		DailyPlans: []models.DailyPlan{},
	}

	if input.DailyPlans != nil {
		for _, dailyPlanInput := range *input.DailyPlans {
			newDailyPlan := createNewDailyPlanFromInput(dailyPlanInput)
			newWeeklyPlan.DailyPlans = append(newWeeklyPlan.DailyPlans, newDailyPlan)
		}
	}

	return newWeeklyPlan
}

func createNewDailyPlanFromInput(input models.DailyPlanUpdateInput) models.DailyPlan {
	newDailyPlan := models.DailyPlan{
		ID:             primitive.NewObjectID(),
		Breakfast:     *input.Breakfast,
		Lunch:         *input.Lunch,
		Dinner:        *input.Dinner,
	}

	if input.MorningSnack != nil {
		newDailyPlan.MorningSnack = *input.MorningSnack
	}
	if input.AfternoonSnack != nil {
		newDailyPlan.AfternoonSnack = *input.AfternoonSnack
	}

	return newDailyPlan
}