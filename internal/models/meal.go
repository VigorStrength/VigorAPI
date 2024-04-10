package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Ingredient represents an ingredient used in a meal, including its quantity and measurement.
type Ingredient struct {
	Name     string  `bson:"name" json:"name" binding:"required" validate:"required"`
	Quantity *string `bson:"quantity,omitempty" json:"quantity,omitempty" validate:"omitempty,min=0"` // Quantity can be nil, e.g., "to taste".
}

// NutritionalInfo represents the nutritional breakdown of a meal per serving.
type NutritionalInfo struct {
	Energy        float64 `bson:"energy" json:"energy" binding:"min=0" validate:"min=0"`               // In KJ or Cal.
	Protein       float64 `bson:"protein" json:"protein" binding:"min=0" validate:"min=0"`             // In grams.
	Fat           float64 `bson:"fat" json:"fat" binding:"min=0" validate:"min=0"`                     // In grams.
	SaturatedFat  float64 `bson:"saturatedFat" json:"saturatedFat" binding:"min=0" validate:"min=0"`   // In grams.
	Carbohydrates float64 `bson:"carbohydrates" json:"carbohydrates" binding:"min=0" validate:"min=0"` // In grams.
	Sugar        float64 `bson:"sugar" json:"sugar" binding:"min=0" validate:"min=0"`               // In grams.
	DietaryFiber         float64 `bson:"dietaryFiber" json:"dietaryFiber" binding:"min=0" validate:"min=0"`                 // In grams.
	Sodium               float64 `bson:"sodium" json:"sodium" binding:"min=0" validate:"min=0"`                         // In milligrams.
	Cholesterol          float64 `bson:"cholesterol" json:"cholesterol" binding:"min=0" validate:"min=0"`                   // In milligrams.
}

// Meal represents a specific meal, including its ingredients, preparation, and nutritional information.
type Meal struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name              string             `bson:"name" json:"name" binding:"required" validate:"required"`
	MealType          string             `bson:"mealType" json:"mealType" binding:"required" validate:"required"` // E.g., Breakfast, Lunch, etc.
	Ingredients       []Ingredient       `bson:"ingredients" json:"ingredients" binding:"required" validate:"required,dive,required"`
	Method            []string           `bson:"method" json:"method" binding:"required" validate:"required"`           // Cooking instructions.
	PrepTime          int                `bson:"prepTime" json:"prepTime" binding:"required" validate:"required,min=5"`       // Preparation time in minutes.
	CookingTime       int                `bson:"cookingTime" json:"cookingTime" binding:"required" validate:"required,min=5"` // Cooking time in minutes.
	NutritionalInfo   NutritionalInfo    `bson:"nutritionalInfo" json:"nutritionalInfo" binding:"required" validate:"required"`
	Description       string             `bson:"description,omitempty" json:"description,omitempty" validate:"omitempty"`          // Optional.
	NutritionalLabels []string           `bson:"nutritionalLabels" json:"nutritionalLabels" validate:"omitempty"`                  // E.g., GF, DF, etc.
	NumberOfServings  int                `bson:"numberOfServings" json:"numberOfServings" binding:"required" validate:"required,min=1"` // Default is 1; can be updated.
}