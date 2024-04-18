package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *UserController) GetUserSubsctiption(c *gin.Context) {
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

	userSubscription, err := uc.UserService.GetUserSubsctiption(c.Request.Context(), objID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		log.Printf("Error getting user subscription: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user subscription"})
		return
	}

	c.JSON(http.StatusOK, userSubscription)
}