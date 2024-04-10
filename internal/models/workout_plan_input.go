package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type WorkoutPlanInput struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	Duration    *int    `json:"duration,omitempty" validate:"omitempty,gte=1,lte=10"`
	Weeks       *[]WorkoutWeekInput `json:"weeks,omitempty" validate:"omitempty,dive"`
}

type WorkoutWeekInput struct {
	Days *[]WorkoutDayInput `json:"days,omitempty" validate:"omitempty,dive"`
	WeekNumber *int `json:"weekNumber,omitempty" validate:"omitempty,gte=1"`
}

type WorkoutDayInput struct {
	WarmUps 		 *[]CircuitInput `json:"warmUps,omitempty"  validate:"omitempty,dive"`
	Workouts 	 *[]CircuitInput `json:"workouts,omitempty"  validate:"omitempty,dive"`
	CoolDowns 	 *[]CircuitInput `json:"coolDowns,omitempty"  validate:"omitempty,dive"`
	WorkoutTimeRange *[2]int `json:"workoutTimeRange,omitempty" validate:"omitempty,dive,gte=1,lte=7200"` // [minTime, maxTime] in seconds.
}

type CircuitInput struct {
	ExerciseIDs  *[]primitive.ObjectID `json:"exerciseIds,omitempty" validate:"omitempty,dive"`
	RestTime     *int                 `json:"restTime,omitempty" validate:"omitempty,gte=5,lte=240"` // Optional rest time in seconds.
	ProposedLaps *int                  `json:"proposedLaps,omitempty" validate:"omitempty,gte=1"`
}