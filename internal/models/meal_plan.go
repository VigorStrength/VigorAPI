package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ingredient represents an ingredient used in a meal, including its quantity and measurement.
type Ingredient struct {
	Name     string  `bson:"name" json:"name" binding:"required"`
	Quantity *string `bson:"quantity,omitempty" json:"quantity,omitempty"` // Quantity can be nil, e.g., "to taste".
}

// NutritionalInfo represents the nutritional breakdown of a meal per serving.
type NutritionalInfo struct {
	Energy        float64 `bson:"energy" json:"energy" binding:"required"`               // In KJ or Cal.
	Protein       float64 `bson:"protein" json:"protein" binding:"required"`             // In grams.
	Fat           float64 `bson:"fat" json:"fat" binding:"required"`                     // In grams.
	SaturatedFat  float64 `bson:"saturatedFat" json:"saturatedFat" binding:"required"`   // In grams.
	Carbohydrates float64 `bson:"carbohydrates" json:"carbohydrates" binding:"required"` // In grams.
	Sugars        float64 `bson:"sugars" json:"sugars" binding:"required"`               // In grams.
	Fiber         float64 `bson:"fiber" json:"fiber" binding:"required"`                 // In grams.
}

// Meal represents a specific meal, including its ingredients, preparation, and nutritional information.
type Meal struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name              string             `bson:"name" json:"name" binding:"required"`
	MealType          string             `bson:"mealType" json:"mealType" binding:"required"` // E.g., Breakfast, Lunch, etc.
	Ingredients       []Ingredient       `bson:"ingredients" json:"ingredients" binding:"required"`
	Method            []string           `bson:"method" json:"method" binding:"required"`           // Cooking instructions.
	PrepTime          int                `bson:"prepTime" json:"prepTime" binding:"required"`       // Preparation time in minutes.
	CookingTime       int                `bson:"cookingTime" json:"cookingTime" binding:"required"` // Cooking time in minutes.
	NutritionalInfo   NutritionalInfo    `bson:"nutritionalInfo" json:"nutritionalInfo" binding:"required"`
	Description       string             `bson:"description,omitempty" json:"description,omitempty"` // Optional.
	NutritionalLabels []string           `bson:"nutritionalLabels" json:"nutritionalLabels"`         // E.g., GF, DF, etc.
	Completed         bool               `bson:"completed" json:"completed"`
	NumberOfServings  int                `bson:"numberOfServings" json:"numberOfServings" binding:"required"` // Default is 1; can be updated.
}

// UserNutritionalLog represents the user's daily log of nutritional intake.
type UserNutritionalLog struct {
	UserID      primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	Date        time.Time          `bson:"date" json:"date" binding:"required"`
	Calories    float64            `bson:"calories" json:"calories"`
	Proteins    float64            `bson:"proteins" json:"proteins"`
	Fats        float64            `bson:"fats" json:"fats"`
	Carbs       float64            `bson:"carbs" json:"carbs"`
	WaterIntake float64            `bson:"waterIntake" json:"waterIntake"` // In liters or ounces.
	// Monthly and daily nutritional goals can be added here later with AI features.
}

// WeeklyPlan represents a weekly grouping of meals by category, used for generic meal proposals.
type WeeklyPlan struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	WeekNumber     *int               `bson:"weekNumber,omitempty" json:"weekNumber,omitempty"` // Optional; not required for generic meal proposals.
	Breakfast      []Meal             `bson:"breakfast" json:"breakfast" binding:"required"`
	MorningSnack   []Meal             `bson:"morningSnack" json:"morningSnack"` // Optional, not all diets may include a morning snack.
	Lunch          []Meal             `bson:"lunch" json:"lunch" binding:"required"`
	AfternoonSnack []Meal             `bson:"afternoonSnack" json:"afternoonSnack"` // Optional.
	Dinner         []Meal             `bson:"dinner" json:"dinner" binding:"required"`
	// Additional meal categories like PreWorkoutSnack and PostWorkoutSnack can be added here.
}

// MealPlan represents the structured plan of meals for a user over a month, organized into weekly plans.
type MealPlan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	StartDate   time.Time          `bson:"startDate" json:"startDate" binding:"required"`
	EndDate     time.Time          `bson:"endDate" json:"endDate" binding:"required"`
	WeeklyPlans []WeeklyPlan       `bson:"weeklyPlans" json:"weeklyPlans" binding:"required"` // Array of weekly meal plans, customized per user.
}
