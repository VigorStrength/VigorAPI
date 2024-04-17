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

func (ac *AdminController) GetMealPlanByID(c *gin.Context) {
	mealPlanID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing meal plan ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	mealPlan, err := ac.AdminService.GetMealPlanByID(c.Request.Context(), mealPlanID)
	if err != nil {
		if errors.Is(err, services.ErrMealPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
			return
		}

		log.Printf("Error getting meal plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meal plan"})
		return
	}

	c.JSON(http.StatusOK, mealPlan)
}

func (ac *AdminController) GetMealPlans(c *gin.Context) {
	mealPlans, err := ac.AdminService.GetMealPlans(c.Request.Context())
	if err != nil {
		log.Printf("Error getting meal plans: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meal plans"})
		return
	}

	c.JSON(http.StatusOK, mealPlans)
}

func (ac *AdminController) SearchMealPlansByName(c *gin.Context) {
	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name query parameter is required"})
		return
	}

	mealPlans, err := ac.AdminService.SearchMealPlansByName(c.Request.Context(), name)
	if err != nil {
		log.Printf("Error searching meal plans by name: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search meal plans"})
		return
	}

	c.JSON(http.StatusOK, mealPlans)
}

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

func (ac *AdminController) UpdateMealPlan(c *gin.Context) {
	mealPlanID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing meal plan ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	var mealPlan models.MealPlanUpdateInput

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

	if err := ac.AdminService.UpdateMealPlan(c.Request.Context(), mealPlanID, mealPlan); err != nil {
		if errors.Is(err, services.ErrMealPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
			return
		}

		log.Printf("Error updating meal plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meal plan updated successfully"})
}


func (ac *AdminController) DeleteMealPlan(c *gin.Context) {
	mealPlanID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing meal plan ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	if err := ac.AdminService.DeleteMealPlan(c.Request.Context(), mealPlanID); err != nil {
		if errors.Is(err, services.ErrMealPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
			return
		}

		log.Printf("Error deleting meal plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meal plan deleted successfully"})
}