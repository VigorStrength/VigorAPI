package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


//DailyPlan represents meals suggested to a user througout the day
type DailyPlan struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Breakfast      primitive.ObjectID `bson:"breakfast" json:"breakfast" validation:"required"`
	MorningSnack   primitive.ObjectID `bson:"morningSnack,omitempty" json:"morningSnack,omitempty" validation:"omitempty"` // Optional, not all diets may include a morning snack.
	Lunch          primitive.ObjectID `bson:"lunch" json:"lunch" validation:"required"`
	AfternoonSnack primitive.ObjectID `bson:"afternoonSnack,omitempty" json:"afternoonSnack,omitempty" validation:"omitempty"` // Optional.
	Dinner         primitive.ObjectID `bson:"dinner" json:"dinner" validation:"required"`
}

// WeeklyPlan represents a weekly grouping of meals by category, used for generic meal proposals.
type WeeklyPlan struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	WeekNumber     int                  `bson:"weekNumber,omitempty" json:"weekNumber,omitempty" validation:"omitempty"` // Optional; not required for chef's picks.
	DailyPlans     []DailyPlan          `bson:"dailyPlans" json:"dailyPlans" binding:"required" validate:"required,dive"`
}
type MealPlan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name" validation:"required"`
	Duration    int                `bson:"duration" json:"duration" validation:"required"`
	WeeklyPlans []WeeklyPlan       `bson:"weeklyPlans" json:"weeklyPlans" validation:"required"`
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
