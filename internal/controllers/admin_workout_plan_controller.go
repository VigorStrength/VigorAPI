package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/gin-gonic/gin"
)


func (ac *AdminController) CreateWorkoutPlan(c *gin.Context) {
	var workoutPlan models.WorkoutPlan

	if err := c.ShouldBindJSON(&workoutPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(workoutPlan); err != nil {
		log.Printf("Error validating input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.CreateWorkoutPlan(c, workoutPlan); err != nil {
		if errors.Is(err, services.ErrWorkoutPlanAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Workout plan already exists"})
			return
		}

		log.Printf("Error creating workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Workout plan created successfully"})
}