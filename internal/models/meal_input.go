package models

type MealUpdateInput struct {
	Name              *string    `json:"name" validate:"omitempty"`
	MealType          *string    `json:"mealType" validate:"omitempty"`
	Ingredients       *[]IngredientUpdateInput `json:"ingredients" validate:"omitempty,dive,omitempty"`
	Method            *[]string  `json:"method" validate:"omitempty"`
	PrepTime          *int       `json:"prepTime" validate:"omitempty,min=5"`
	CookingTime       *int       `json:"cookingTime" validate:"omitempty,min=5"`
	NutritionalInfo   *NutritionalInfoUpdateInput `json:"nutritionalInfo" validate:"omitempty"`
	Description       *string    `json:"description" validate:"omitempty"`
	NutritionalLabels *[]string  `json:"nutritionalLabels" validate:"omitempty"`
	NumberOfServings  *int       `json:"numberOfServings" validate:"omitempty,min=1"`
}

type IngredientUpdateInput struct {
	Name     *string `json:"name" validate:"omitempty"`
	Quantity *string `json:"quantity,omitempty" validate:"omitempty,min=0"`
}

type NutritionalInfoUpdateInput struct {
	Energy        *float64 `json:"energy" validate:"omitempty,min=0"`
	Protein       *float64 `json:"protein" validate:"omitempty,min=0"`
	Fat           *float64 `json:"fat" validate:"omitempty,min=0"`
	SaturatedFat  *float64 `json:"saturatedFat" validate:"omitempty,min=0"`
	Carbohydrates *float64 `json:"carbohydrates" validate:"omitempty,min=0"`
	Sugar         *float64 `json:"sugar" validate:"omitempty,min=0"`
	DietaryFiber  *float64 `json:"dietaryFiber" validate:"omitempty,min=0"`
	Sodium        *float64 `json:"sodium" validate:"omitempty,min=0"`
	Cholesterol   *float64 `json:"cholesterol" validate:"omitempty,min=0"`
}