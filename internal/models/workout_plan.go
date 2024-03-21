package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Exercise represents a specific exercise, including details for logging and video interaction.
type Exercise struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string             `bson:"name" json:"name" binding:"required"`
	Description     string             `bson:"description" json:"description" binding:"required"`
	VideoURL        string             `bson:"videoUrl" json:"videoUrl" binding:"required"`
	TargetMuscles   []string           `bson:"targetMuscles" json:"targetMuscles" binding:"required"`
	EquipmentNeeded []string           `bson:"equipmentNeeded" json:"equipmentNeeded" binding:"required"`
	Instructions    []string           `bson:"instructions" json:"instructions" binding:"required"`
	Time            int                `bson:"time" json:"time" binding:"required"` // Time in seconds.
	Log             ExerciseLog        `bson:"log" json:"log"`
}

// ExerciseLog represents the logging of exercises, including proposed logs and actual user logs.
type ExerciseLog struct {
	ProposedReps   int      `bson:"proposedReps" json:"proposedReps" binding:"required"`
	ActualReps     int      `bson:"actualReps" json:"actualReps"`
	ProposedWeight float64  `bson:"proposedWeight" json:"proposedWeight"` // Can be empty if no equipment.
	ActualWeight   float64  `bson:"actualWeight" json:"actualWeight"`
}

// Circuit represents a set of exercises performed in sequence, with optional rest and laps tracking.
type Circuit struct {
	ID 			 primitive.ObjectID	  `bson:"_id,omitempty" json:"id,omitempty"`
	ExerciseIDs  []primitive.ObjectID `bson:"exerciseIds" json:"exerciseIds" binding:"required"`
	RestTime     *int                 `bson:"restTime,omitempty" json:"restTime,omitempty"` // Optional rest time in seconds.
	ProposedLaps int                  `bson:"proposedLaps" json:"proposedLaps" binding:"required"`
	ActualLaps   int                  `bson:"actualLaps" json:"actualLaps"`
	//To be added later once the generative AI is integrated, program will be distinct from each user
	// Completed    bool                 `bson:"completed" json:"completed"`
}

// UserCircuitStatus tracks the completion status of a circuit for a specific user.
type UserCircuitStatus struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
	CircuitID  primitive.ObjectID `bson:"circuitId" json:"circuitId" binding:"required"`
	Completed  bool               `bson:"completed" json:"completed"`
}

// WorkoutDay represents a complete day's workout plan, including warm-up, workout, and cool-down.
type WorkoutDay struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	WarmUp           []Circuit `bson:"warmUp" json:"warmUp" binding:"required"`
	Workout          []Circuit `bson:"workout" json:"workout" binding:"required"`
	CoolDown         []Circuit `bson:"coolDown" json:"coolDown" binding:"required"`
	TotalExercises   int       `bson:"totalExercises" json:"totalExercises"`
	TotalLaps        int       `bson:"totalLaps" json:"totalLaps"`
	WorkoutTimeRange [2]int    `bson:"workoutTimeRange" json:"workoutTimeRange"` // [minTime, maxTime] in seconds.
	Equipment        []string  `bson:"equipment" json:"equipment"`
	//To be added later once the generative AI is integrated, program will be distinct from each user
	// CompletionStatus bool      `bson:"completionStatus" json:"completionStatus"` // Indicates if the WorkoutDay is completed.
}

// UserWorkoutDayStatus tracks the completion status of a workout day for a specific user.
type UserWorkoutDayStatus struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID       primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
    WorkoutDayID primitive.ObjectID `bson:"workoutDayId" json:"workoutDayId" binding:"required"`
    Completed    bool               `bson:"completed" json:"completed"`
}

// WorkoutWeek represents a single week within a workout plan, containing multiple workout days.
type WorkoutWeek struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Days          []WorkoutDay `bson:"days" json:"days" binding:"required"`
	WeekNumber    int          `bson:"weekNumber" json:"weekNumber" binding:"required"`
	//To be added later once the generative AI is integrated, program will be distinct from each user
	// CompletedDays int          `bson:"completedDays" json:"completedDays"` // Derived by counting days with CompletionStatus true.
}

// UserWorkoutWeekStatus tracks the completion status of a workout week for a specific user.
type UserWorkoutWeekStatus struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID          primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
    WorkoutWeekID   primitive.ObjectID `bson:"workoutWeekId" json:"workoutWeekId" binding:"required"`
    CompletedDays   int                `bson:"completedDays" json:"completedDays"`
}

// WorkoutPlan represents a user's workout plan, including weeks and progress tracking.
type WorkoutPlan struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId" binding:"required"` // Reference to the User.
	StartDate time.Time          `bson:"startDate" json:"startDate" binding:"required"`
	EndDate   time.Time          `bson:"endDate" json:"endDate" binding:"required"`
	Weeks     []WorkoutWeek      `bson:"weeks" json:"weeks" binding:"required"`
	// Progress field is intentionally omitted as it will be calculated elsewhere.
}
