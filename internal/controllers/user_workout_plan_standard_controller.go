package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *UserController) GetStandardWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("Error retrieving userID from context\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	objID, ok := userID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error converting userID from type interface {} to primitive.ObjectID\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert user ID to string"})
		return
	}

	//Get Active WorkoutPlan first to retriev it's ID
	activeWorkoutPlan, err := uc.UserService.GetActiveWorkoutPlan(c.Request.Context(), objID)
	if err != nil {
		if errors.Is(err, services.ErrActiveWorkoutPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User has no active workout plan"})
			return
		}


		log.Printf("Error getting active workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active workout plan"})
		return
	}

	//Get Standard WorkoutPlan
	standardWorkoutPlan, err := uc.UserService.GetWorkoutPlanByID(c.Request.Context(), activeWorkoutPlan.WorkoutPlanID)
	if err != nil {
		if errors.Is(err, services.ErrWorkoutPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
			return
		}

		log.Printf("Error getting workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workout plan"})
		return
	}

	c.JSON(http.StatusOK, standardWorkoutPlan)
}

func (uc *UserController) GetDailyExercisesByIDs(c *gin.Context) {
	//Get the daily exercises by ID's sent in the request body
	var requestBody struct {
		DailyExercisesIDs []primitive.ObjectID `json:"dailyExercisesIDs"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	dailyExercises, err := uc.UserService.GetDailyExercisesByIDs(c.Request.Context(), requestBody.DailyExercisesIDs)
	if err != nil {
		log.Printf("Error getting daily exercises: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily exercises"})
		return
	}

	c.JSON(http.StatusOK, dailyExercises)
}

func (uc *UserController) GetDailySetByID(c *gin.Context) {
	setID, err := primitive.ObjectIDFromHex(c.Param("setId"))
	if err != nil {
		log.Printf("Error parsing ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	dailySet, err := uc.UserService.GetDailySetByID(c.Request.Context(), setID)
	if err !=  nil {
		if errors.Is(err, services.ErrSetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Set not found"})
			return
		}

		log.Printf("Error getting daily set: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily set"})
		return
	}

	c.JSON(http.StatusOK, dailySet)
}

func (uc *UserController) GetDailySetsByIDs(c *gin.Context) {
	//Get the daily sets by ID's sent in the request body
	var requestBody struct {
		DailySetsIDs []primitive.ObjectID `json:"dailySetsIDs"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	dailySets, err := uc.UserService.GetDailySetsByIDs(c.Request.Context(), requestBody.DailySetsIDs)
	if err != nil {
		log.Printf("Error getting daily sets: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily sets"})
		return
	}

	c.JSON(http.StatusOK, dailySets)
}

func (uc *UserController) GetDailyStandAloneWorkoutItemByID(c *gin.Context) {
	standAloneWorkoutItemID, err := primitive.ObjectIDFromHex(c.Param("standaloneId"))
	if err != nil {
		log.Printf("Error parsing ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	standAloneWorkout, err := uc.UserService.GetDailyStandAloneWorkoutItemByID(c.Request.Context(), standAloneWorkoutItemID)
	if err != nil {
		if errors.Is(err, services.ErrStandAloneWorkoutNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stand alone workout not found"})
			return
		}

		log.Printf("Error getting stand alone workout: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stand alone workout"})
		return
	}

	c.JSON(http.StatusOK, standAloneWorkout)
}

func (uc *UserController) GetDailyStandAloneWorkoutItemsByIDs(c *gin.Context) {
	//Get the daily stand alone workouts by ID's sent in the request body
	var requestBody struct {
		DailyStandAloneWorkoutsIDs []primitive.ObjectID `json:"dailyStandAloneWorkoutsIDs"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	dailyStandAloneWorkouts, err := uc.UserService.GetDailyStandAloneWorkoutItemsByIDs(c.Request.Context(), requestBody.DailyStandAloneWorkoutsIDs)

	if err != nil {
		log.Printf("Error getting daily stand alone workouts: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily stand alone workouts"})
		return
	}

	c.JSON(http.StatusOK, dailyStandAloneWorkouts)
}


func (uc *UserController) GetDailySupertsetByID(c *gin.Context) {
	supersetID, err := primitive.ObjectIDFromHex(c.Param("supersetId"))
	if err != nil {
		log.Printf("Error parsing ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	dailySuperset, err := uc.UserService.GetDailySupersetByID(c.Request.Context(), supersetID)
	if err != nil {
		if errors.Is(err, services.ErrDailySupersetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Daily Superset not found"})
			return
		}

		log.Printf("Error getting daily superset: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily superset"})
		return
	}

	c.JSON(http.StatusOK, dailySuperset)
}

func (uc *UserController) GetDailySupersetsByIDs(c *gin.Context) {
	//Get the daily supersets by ID's sent in the request body

	var requestBody struct {
		DailySupersetsIDs []primitive.ObjectID `json:"dailySupersetsIDs"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	dailySuperSets, err := uc.UserService.GetDailySupersetsByIDs(c.Request.Context(), requestBody.DailySupersetsIDs)
	if err != nil {
		log.Printf("Error getting daily supersets: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily supersets"})
		return
	}

	c.JSON(http.StatusOK, dailySuperSets)
}

