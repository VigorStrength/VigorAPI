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