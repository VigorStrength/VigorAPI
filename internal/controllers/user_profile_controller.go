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

func (uc *UserController) GetUserProfile(c *gin.Context) {
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

	userProfile, err := uc.UserService.GetUserProfile(c.Request.Context(), objID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		log.Printf("Error getting user profile: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

func (uc *UserController) UpdateUserProfile(c *gin.Context) {
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

	var updateInput models.UserProfileUpdateInput
	if err := c.BindJSON(&updateInput); err != nil {
		log.Printf("Error binding json to update input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind json to update input"})
		return
	}

	if err := validate.Struct(updateInput); err != nil {
		log.Printf("Error validating update input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update input"})
		return
	}

	_, err := uc.UserService.UpdateUserProfile(c.Request.Context(), objID, updateInput)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			log.Printf("Error updating user profile: %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		log.Printf("Error updating user profile: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully"})
}