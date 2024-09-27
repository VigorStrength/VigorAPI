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

func (ac *AdminController) CreateStandAloneWorkoutItem(c *gin.Context) {
	var standAloneWorkout models.StandAlone

	if err := c.ShouldBindJSON(&standAloneWorkout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(standAloneWorkout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.CreateStandAloneWorkoutItem(c.Request.Context(), standAloneWorkout); err != nil {
		if errors.Is(err, services.ErrStandAloneWorkoutAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Stand alone workout already exists"})
			return
		}

		log.Printf("Error creating stand alone workout: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stand alone workout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stand alone workout created successfully"})
}

func (ac *AdminController) UpdateStandAloneWorkoutItem(c *gin.Context) {
	standAloneWorkoutID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var standAloneWorkoutInput models.StandAloneUpdateInput
	if err := c.ShouldBindJSON(&standAloneWorkoutInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(standAloneWorkoutInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.UpdateStandAloneWorkoutItem(c.Request.Context(), standAloneWorkoutID, standAloneWorkoutInput); err != nil {
		if errors.Is(err, services.ErrStandAloneWorkoutNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stand alone workout not found"})
			return
		}

		log.Printf("Error updating stand alone workout: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stand alone workout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stand alone workout updated successfully"})
}

func (ac *AdminController) GetStandAloneWorkoutItems(c *gin.Context) {
	standAloneWorkouts, err := ac.AdminService.GetStandAloneWorkoutItems(c.Request.Context())
	if err != nil {
		log.Printf("Error getting stand alone workouts: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stand alone workouts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"standAloneWorkouts": standAloneWorkouts})
}

