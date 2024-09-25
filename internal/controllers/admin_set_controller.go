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

func (ac *AdminController) CreateSet(c *gin.Context) {
	var set models.Set

	if err := c.ShouldBindJSON(&set); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(set); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.CreateSet(c.Request.Context(), set); err != nil {
		if errors.Is(err, services.ErrSetAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Set already exists"})
			return
		}

		log.Printf("Error creating set: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create set"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Set created successfully"})
}

func (ac *AdminController) UpdateSet(c *gin.Context) {
	setID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var setInput models.SetUpdateInput
	if err := c.ShouldBindJSON(&setInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(setInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.UpdateSet(c, setID, setInput); err != nil {
		if errors.Is(err, services.ErrSetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Set not found"})
			return
		}

		log.Printf("Error updating set: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update set"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Set updated successfully"})
}

func (ac *AdminController) GetSets(c *gin.Context) {
	sets, err := ac.AdminService.GetSets(c.Request.Context())
	if err != nil {
		log.Printf("Error getting sets: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sets": sets})
}