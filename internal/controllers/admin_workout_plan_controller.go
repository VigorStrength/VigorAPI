package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ac *AdminController) GetWorkoutPlanByID(c *gin.Context) {
	workoutPlanID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing workout plan ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout plan ID"})
		return
	}

	workoutPlan, err := ac.AdminService.GetWorkoutPlanByID(c.Request.Context(), workoutPlanID)
	if err != nil {
		if errors.Is(err, services.ErrWorkoutPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
			return
		}

		log.Printf("Error getting workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workout plan"})
		return
	}

	c.JSON(http.StatusOK, workoutPlan)
}

func (ac *AdminController) GetWorkoutPlans(c *gin.Context) {
	workoutPlans, err := ac.AdminService.GetWorkoutPlans(c.Request.Context())
	if err != nil {
		log.Printf("Error getting workout plans: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workout plans"})
		return
	}

	c.JSON(http.StatusOK, workoutPlans)
}

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

func (ac *AdminController) UpdateWorkoutPlan(c *gin.Context) {
	workoutPlanID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing workout plan ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout plan ID"})
		return
	}

	var updateInput models.WorkoutPlanInput

	if err := c.ShouldBindJSON(&updateInput); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(updateInput); err != nil {
		log.Printf("Error validating input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.UpdateWorkoutPlan(c, workoutPlanID, updateInput); err != nil {
		if errors.Is(err, services.ErrWorkoutPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
			return
		}

		log.Printf("Error updating workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout plan updated successfully"})
}

func (ac *AdminController) DeleteWorkoutPlan(c *gin.Context) {
	workoutPlanID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing workout plan ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout plan ID"})
		return
	}

	if err := ac.AdminService.DeleteWorkoutPlan(c, workoutPlanID); err != nil {
		if errors.Is(err, services.ErrWorkoutPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
			return
		}

		log.Printf("Error deleting workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout plan deleted successfully"})
}