package models

type UserProfileUpdateInput struct {
	Username *string `json:"username" validate:"omitempty"`
	ProfilePicture *string `json:"profilePicture" validate:"omitempty,url"`
	MainGoal *string `json:"mainGoal" validate:"omitempty"`
	SecondaryGoal *string `json:"secondaryGoal,omitempty" validate:"omitempty"`
	BodyInformation *BodyInformationUpdateInput `json:"bodyInformation" validate:"omitempty"`
	ActivityLevel *string `json:"activityLevel" validate:"omitempty"`
	PhysicalActivity *PhysicalActivityUpdateInput `json:"physicalActivity" validate:"omitempty"`
	Lifestyle *LifestyleUpdateInput `json:"lifestyle" validate:"omitempty"`
	CycleInformation *CycleInformationUpdateInput `json:"cycleInformation,omitempty" validate:"omitempty"`
}

type BodyInformationUpdateInput struct {
	BodyType *string `json:"bodyType" validate:"omitempty"`
	BodyGoal *string `json:"bodyGoal,omitempty" validate:"omitempty"`
	HealthRestrictions *[]string `json:"healthRestrictions" validate:"omitempty,dive"`
	FocusArea *[]string `json:"focusArea" validate:"omitempty,dive"`
}

type PhysicalActivityUpdateInput struct {
	FitnessLevel *string `json:"fitnessLevel" validate:"omitempty"`
	Activities *[]string `json:"activities" validate:"omitempty,dive"`
}

type LifestyleUpdateInput struct {
	Diet *string `json:"diet" validate:"omitempty"`
	WaterIntake *int `json:"waterIntake,omitempty" validate:"omitempty,min=0"`
	SleepDuration *int `json:"sleepDuration,omitempty" validate:"omitempty,min=0"`
	TypicalDay *string `json:"typicalDay" validate:"omitempty"`
	TrainingLocation *string `json:"trainingLocation" validate:"omitempty"`
	WorkoutTime *string `json:"workoutTime" validate:"omitempty"`
	WorkoutFrequency *int `json:"workoutFrequency,omitempty" validate:"omitempty,min=0"`
	WorkoutDuration *string `json:"workoutDuration" validate:"omitempty"`
	DiscoveryMethod *string `json:"discoveryMethod,omitempty" validate:"omitempty"`
	IntolerancesAndAllergies *[]string `json:"intolerancesAndAllergies,omitempty" validate:"omitempty,dive"`
}

type CycleInformationUpdateInput struct {
	ReproductiveStage *string `json:"reproductiveStage" validate:"omitempty"`
}