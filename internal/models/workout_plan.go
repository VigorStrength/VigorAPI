package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Circuit represents a set of exercises performed in sequence, with optional rest and laps tracking.
type Circuit struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	ExerciseIDs  []primitive.ObjectID `bson:"exerciseIds" json:"exerciseIds" binding:"required" validate:"required,dive,required"`
	RestTime     *int                 `bson:"restTime,omitempty" json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Optional rest time in seconds.
	ProposedLaps int                  `bson:"proposedLaps" json:"proposedLaps" binding:"required" validate:"required,gte=1"`
	//To be added later once the generative AI is integrated, program will be distinct from each user
	// ActualLaps   int                  `bson:"actualLaps" json:"actualLaps"` // Actual laps completed by the user.
	// Completed    bool                 `bson:"completed" json:"completed"`
}

// WorkoutDay represents a complete day's workout plan, including warm-up, workout, and cool-down.
type WorkoutDay struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name 		     string             `bson:"name" json:"name" binding:"required" validate:"required,min=5,max=50"`
	ImageURL		 string             `bson:"imageURL" json:"imageURL" binding:"required" validate:"required,url"`
	WarmUps           []Circuit          `bson:"warmUps" json:"warmUps" binding:"required" validate:"required,dive"`
	Workouts          []Circuit          `bson:"workouts" json:"workouts" binding:"required" validate:"required,dive"`
	CoolDowns         []Circuit          `bson:"coolDowns" json:"coolDowns" binding:"required" validate:"required,dive"`
	WorkoutTimeRange [2]int             `bson:"workoutTimeRange" json:"workoutTimeRange" binding:"required" validate:"required,dive,gte=1,lte=7200"` // [minTime, maxTime] in seconds.
	// Removed TotalExercises and Equipment fields
	// Other fields...
}

// WorkoutWeek represents a single week within a workout plan, containing multiple workout days.
type WorkoutWeek struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Days       []WorkoutDay       `bson:"days" json:"days" binding:"required" validate:"required,dive"`
	WeekNumber int                `bson:"weekNumber" json:"weekNumber" binding:"required" validate:"required,gte=1"`
	//To be added later once the generative AI is integrated, program will be distinct from each user
	// CompletedDays int          `bson:"completedDays" json:"completedDays"` // Derived by counting days with CompletionStatus true.
}

type WorkoutPlan struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name" binding:"required" validate:"required,min=5,max=50"`
	ImageURL string             `bson:"imageURL" json:"imageURL" binding:"required" validate:"required,url"`
	Duration int                `bson:"duration" json:"duration" binding:"required" validate:"required,gt=0"`
	Weeks    []WorkoutWeek      `bson:"weeks" json:"weeks" binding:"required" validate:"required,dive"`
}

// AI WorkoutPlan
// // WorkoutPlan represents a user's workout plan, including weeks and progress tracking.
// type WorkoutPlan struct {
// 	ID        primitive.ObjectID `bson:"_id" json:"id"`
// 	UserID    primitive.ObjectID `bson:"userId" json:"userId" binding:"required"` // Reference to the User.
// 	StartDate time.Time          `bson:"startDate" json:"startDate" binding:"required"`
// 	EndDate   time.Time          `bson:"endDate" json:"endDate" binding:"required"`
// 	Weeks     []WorkoutWeek      `bson:"weeks" json:"weeks" binding:"required"`
// 	// Progress field is intentionally omitted as it will be calculated elsewhere.
// }
