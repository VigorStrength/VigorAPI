package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserDailyNutritionalLog represents the user's daily log of nutritional intake.
type UserDailyNutritionalLog struct {
	UserID      primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	Date        time.Time          `bson:"date" json:"date" binding:"required"`
	Calories    float64            `bson:"calories" json:"calories"`
	Proteins    float64            `bson:"proteins" json:"proteins"`
	Fats        float64            `bson:"fats" json:"fats"`
	Carbs       float64            `bson:"carbs" json:"carbs"`
	WaterIntake float64            `bson:"waterIntake" json:"waterIntake"` // In liters or ounces.
	// Monthly and daily nutritional goals can be added here later with AI features.
}

type UserMealStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	MealID    primitive.ObjectID `bson:"mealId" json:"mealId" binding:"required"`
	Completed bool               `bson:"completed" json:"completed"`
}

// UserWeeklyPlanStatus tracks the completion status of a meal plan week for a specific user
type UserWeeklyPlanStatus struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	WeeklyPlanId  primitive.ObjectID `bson:"weeklyPlanId" json:"weeklyPlanId" binding:"required"`
	CompletedDays int                `bson:"completedDays" json:"completedDays"`
}

type UserMealPlanStatus struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	MealPlanID     primitive.ObjectID `bson:"mealPlanId" json:"mealPlanId" binding:"required"`
	StartDate      time.Time          `bson:"startDate" json:"startDate" binding:"required"`
	CompletionDate *time.Time         `bson:"completionDate" json:"completionDate"`
	Completed      bool               `bson:"completed" json:"completed"`
}