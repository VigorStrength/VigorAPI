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

// WeeklyPlan represents a weekly grouping of meals by category, used for generic meal proposals.
type WeeklyPlan struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	WeekNumber     int                  `bson:"weekNumber,omitempty" json:"weekNumber,omitempty"` // Optional; not required for chef's picks.
	Breakfast      []primitive.ObjectID `bson:"breakfast" json:"breakfast"`
	MorningSnack   []primitive.ObjectID `bson:"morningSnack" json:"morningSnack"` // Optional, not all diets may include a morning snack.
	Lunch          []primitive.ObjectID `bson:"lunch" json:"lunch"`
	AfternoonSnack []primitive.ObjectID `bson:"afternoonSnack" json:"afternoonSnack"` // Optional.
	Dinner         []primitive.ObjectID `bson:"dinner" json:"dinner"`
}

// UserWeeklyPlanStatus tracks the completion status of a meal plan week for a specific user
type UserWeeklyPlanStatus struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	WeeklyPlanId  primitive.ObjectID `bson:"weeklyPlanId" json:"weeklyPlanId" binding:"required"`
	CompletedDays int                `bson:"completedDays" json:"completedDays"`
}

type MealPlan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Duration    int                `bson:"duration" json:"duration" binding:"required"`
	WeeklyPlans []WeeklyPlan       `bson:"weeklyPlans" json:"weeklyPlans" binding:"required"`
}

type UserMealPlanStatus struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	MealPlanID     primitive.ObjectID `bson:"mealPlanId" json:"mealPlanId" binding:"required"`
	StartDate      time.Time          `bson:"startDate" json:"startDate" binding:"required"`
	CompletionDate *time.Time         `bson:"completionDate" json:"completionDate"`
	Completed      bool               `bson:"completed" json:"completed"`
}

// AI MealPlan
// MealPlan represents the structured plan of meals for a user over a month, organized into weekly plans.
// type MealPlan struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	UserID      primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
// 	StartDate   time.Time          `bson:"startDate" json:"startDate" binding:"required"`
// 	EndDate     time.Time          `bson:"endDate" json:"endDate" binding:"required"`
// 	WeeklyPlans []WeeklyPlan       `bson:"weeklyPlans" json:"weeklyPlans" binding:"required"` // Array of weekly meal plans, customized per user.
// }
