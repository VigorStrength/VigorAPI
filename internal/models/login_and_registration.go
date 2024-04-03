package models

import "time"

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminRegistrationInput struct {
    Email        string             `json:"email" binding:"required" validate:"required,email,endswith=@vigor.com"`
    Password string             `json:"password" binding:"required" validate:"required,min=8,max=12"`
}

type UserRegistrationInput struct {
    FirstName          string    `json:"firstName" validate:"required,alpha"`
    LastName           string    `json:"lastName" validate:"required,alpha"`
    Email              string    `json:"email" validate:"required,email"`
    Password           string    `json:"password" validate:"required,min=8,max=12"`
    BirthDate          time.Time `json:"birthDate" validate:"required"`
    Gender             string    `json:"gender" validate:"required,oneof=male female"`
    Height             int       `json:"height" validate:"required,gt=0"`
    Weight             int       `json:"weight" validate:"required,gt=0"`
    Subscription       SubscriptionInput   `json:"subscription" validate:"required"`
    TrialEndsAt        time.Time `json:"trialEndsAt" validate:"required"`
    ProfileInformation UserProfileInput    `json:"profileInformation" validate:"required"`
    SystemPreferences  *SystemPreferencesInput  `json:"systemPreferences,omitempty" validate:"omitempty"`
}

type SubscriptionInput struct {
    Type            string    `json:"type" validate:"required,oneof=basic premium"`
    Status          string    `json:"status" validate:"required,oneof=active inactive"`
    StartDate       time.Time `json:"startDate" validate:"required"`
    EndDate         *time.Time `json:"endDate,omitempty" validate:"omitempty"`
    NextRenewalDate *time.Time `json:"nextRenewalDate,omitempty" validate:"omitempty"`
    IsActive        bool      `json:"isActive" validate:"required"`
}

type UserProfileInput struct {
    Username         string            `json:"username" validate:"required,alphanum"`
    ProfilePicture   string            `json:"profilePicture" validate:"omitempty,url"`
    MainGoal         string            `json:"mainGoal" validate:"required"`
    SecondaryGoal    string            `json:"secondaryGoal" validate:"omitempty"`
    BodyInformation  BodyInformationInput `json:"bodyInformation" validate:"required"`
    ActivityLevel    string            `json:"activityLevel" validate:"required,oneof=low medium high"`
    PhysicalActivity PhysicalActivityInput `json:"physicalActivity" validate:"required"`
    Lifestyle        LifestyleInput    `json:"lifestyle" validate:"required"`
    CycleInformation *CycleInformationInput `json:"cycleInformation,omitempty" validate:"omitempty"`
}

type BodyInformationInput struct {
    BodyType           string   `json:"bodyType" validate:"required,oneof=ectomorph mesomorph endomorph"`
    BodyGoal           *string  `json:"bodyGoal,omitempty" validate:"omitempty"`
    HealthRestrictions []string `json:"healthRestrictions,omitempty" validate:"omitempty,dive,required"`
    FocusArea          []string `json:"focusArea" validate:"omitempty,dive,required"`
}

type PhysicalActivityInput struct {
    FitnessLevel string   `json:"fitnessLevel" validate:"required,oneof=beginner intermediate advanced"`
    Activities   []string `json:"activities" validate:"required,dive,required"`
}

type LifestyleInput struct {
    Diet                     string   `json:"diet" validate:"required,oneof=vegan vegetarian omnivore keto paleo"`
    WaterIntake              *int     `json:"waterIntake,omitempty" validate:"omitempty,gt=0"`
    SleepDuration            *int     `json:"sleepDuration,omitempty" validate:"omitempty,gt=0"`
    TypicalDay               string   `json:"typicalDay" validate:"required,oneof=sedentary moderate active"`
    TrainingLocation         string   `json:"trainingLocation" validate:"required,oneof=home gym outdoor"`
    WorkoutTime              string   `json:"workoutTime" validate:"required"`
    WorkoutFrequency         *int     `json:"workoutFrequency,omitempty" validate:"omitempty,gt=0"`
    WorkoutDuration          string   `json:"workoutDuration" validate:"required"`
    DiscoveryMethod          *string  `json:"discoveryMethod,omitempty" validate:"omitempty"`
    IntolerancesAndAllergies []string `json:"intolerancesAndAllergies,omitempty" validate:"omitempty,dive,required"`
}

type CycleInformationInput struct {
    ReproductiveStage string `json:"reproductiveStage" validate:"required,oneof=menstrual_cycle pregnancy recently_gave_birth menopause no_period rather_not_reply"`
}

type SystemPreferencesInput struct {
    Language          string `json:"language,omitempty" validate:"omitempty"`
    TimeZone          string `json:"timeZone,omitempty" validate:"omitempty"`
    DisplayMode       string `json:"displayMode,omitempty" validate:"omitempty,oneof=light dark"`
    MeasurementSystem string `json:"measurementSystem,omitempty" validate:"omitempty,oneof=metric imperial"`
    AllowReadReceipt  bool   `json:"allowReadReceipt,omitempty" validate:"required,omitempty"`
}
