package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DailyPlanUpdateInput struct {
	Breakfast     *primitive.ObjectID `json:"breakfast,omitempty" validate:"omitempty"`
	MorningSnack  *primitive.ObjectID `json:"morningSnack,omitempty" validate:"omitempty"`
	Lunch         *primitive.ObjectID `json:"lunch,omitempty" validate:"omitempty"`
	AfternoonSnack *primitive.ObjectID `json:"afternoonSnack,omitempty" validate:"omitempty"`
	Dinner        *primitive.ObjectID `json:"dinner,omitempty" validate:"omitempty"`
}

type WeeklyPlanUpdateInput struct {
	WeekNumber *int `json:"weekNumber,omitempty" validate:"omitempty,gte=1"`
	DailyPlans *[]DailyPlanUpdateInput `json:"dailyPlans,omitempty" validate:"omitempty,dive"`
}

type MealPlanUpdateInput struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	Duration    *int    `json:"duration,omitempty" validate:"omitempty,gte=1,lte=10"`
	WeeklyPlans *[]WeeklyPlanUpdateInput `json:"weeklyPlans,omitempty" validate:"omitempty,dive"`
}