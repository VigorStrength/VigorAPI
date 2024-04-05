package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/gin-gonic/gin"
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