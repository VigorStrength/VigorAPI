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

func (ac *AdminController) CreateExercise(c *gin.Context) {
	var exercise models.Exercise

	if err := c.ShouldBindJSON(&exercise); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(exercise); err != nil {
		log.Printf("Error validating input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.CreateExercise(c.Request.Context(), exercise); err != nil {
		if errors.Is(err, services.ErrExerciseAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Exercise already exists"})
			return
		}

		log.Printf("Error creating exercise: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise created successfully"})
}

func (ac *AdminController) UpdateExercise(c *gin.Context) {
	exerciseID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updateInput models.ExerciseUpdateInput

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

	if err := ac.AdminService.UpdateExercise(c.Request.Context(), exerciseID, updateInput); err != nil {
		if errors.Is(err, services.ErrExerciseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
			return
		}

		log.Printf("Error updating exercise: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise updated successfully"})
}