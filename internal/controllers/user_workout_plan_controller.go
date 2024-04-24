package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *UserController) JoinWorkoutPlan(c *gin.Context) {
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

	workoutPlanID, err := primitive.ObjectIDFromHex(c.Param("workoutPlanId"))
	if err != nil {
		log.Printf("Error parsing workout plan ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout plan ID"})
		return
	}

	if err := uc.UserService.JoinWorkoutPlan(c.Request.Context(), objID, workoutPlanID); err != nil {
		if errors.Is(err, services.ErrWorkoutPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
			return
		}
		if errors.Is(err, services.ErrAlreadyJoinded) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User has already joined this workout plan"})
			return
		}

		log.Printf("Error joining workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join workout plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined workout plan"})
}