package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/gin-gonic/gin"
)

func (ac *AdminController) CreateMealPlan(c *gin.Context) {
	var mealPlan models.MealPlan

	if err := c.ShouldBindJSON(&mealPlan); err != nil {
		log.Printf("Error binding meal plan: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(mealPlan); err != nil {
		log.Printf("Error validating meal plan: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.AdminService.CreateMealPlan(c.Request.Context(), mealPlan); err != nil {
		if errors.Is(err, services.ErrMealPlanAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Meal plan already exists"})
			return
		}

		log.Printf("Error creating meal plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal plan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Meal plan created successfully"})
}