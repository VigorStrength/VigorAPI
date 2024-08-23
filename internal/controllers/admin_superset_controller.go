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

func (ac *AdminController) CreateSuperset(c *gin.Context) {
	var superset models.Superset

	if err := c.ShouldBindJSON(&superset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(superset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.CreateSuperset(c.Request.Context(), superset); err != nil {
		if errors.Is(err, services.ErrSupersetAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Superset already exists"})
			return
		}

		log.Printf("Error creating superset: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create superset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Superset created successfully"})
}

func (ac *AdminController) UpdateSuperset(c *gin.Context) {
	supersetID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var supersetInput models.SupersetUpdateInput
	if err := c.ShouldBindJSON(&supersetInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(supersetInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.UpdateSuperset(c, supersetID, supersetInput); err != nil {
		if errors.Is(err, services.ErrSupersetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Superset not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Superset updated successfully"})
}

func (ac *AdminController) GetSupersets(c *gin.Context) {
	supersets, err := ac.AdminService.GetSupersets(c.Request.Context())
	if err != nil {
		log.Printf("Error getting supersets: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get supersets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"supersets": supersets})
}