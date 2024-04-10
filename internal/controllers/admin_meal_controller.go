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

func (ac *AdminController) GetMealByID(c *gin.Context) {
	mealID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	meal, err := ac.AdminService.GetMealByID(c.Request.Context(), mealID)
	if err != nil {
		if errors.Is(err, services.ErrMealNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meal not found"})
			return
		}

		log.Printf("Error getting meal: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meal"})
		return
	}

	c.JSON(http.StatusOK, meal)
}

func (ac *AdminController) GetMeals(c *gin.Context) {
	meals, err := ac.AdminService.GetMeals(c.Request.Context())
	if err != nil {
		log.Printf("Error getting meals: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meals"})
		return
	}

	c.JSON(http.StatusOK, meals)
}


func (ac *AdminController) CreateMeal(c *gin.Context) {
	var meal models.Meal

	if err := c.ShouldBindJSON(&meal); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(meal); err != nil {
		log.Printf("Error validating input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.CreateMeal(c.Request.Context(), meal); err != nil {
		if errors.Is(err, services.ErrMealAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Meal already exists"})
			return
		}

		log.Printf("Error creating meal: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Meal created successfully"})
}

func (ac *AdminController) DeleteMeal(c *gin.Context) {
	mealID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := ac.AdminService.DeleteMeal(c.Request.Context(), mealID); err != nil {
		if errors.Is(err, services.ErrMealNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meal not found"})
			return
		}

		log.Printf("Error deleting meal: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meal deleted successfully"})
}

func (ac *AdminController) SearchMealsByName(c *gin.Context) {
	name := c.Query("name")

	if name == "" {
		log.Printf("Meal name query parameter must be specified: %v\n", name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

	meals, err := ac.AdminService.SearchMealsByName(c.Request.Context(), name)
	if err != nil {
		log.Printf("Error searching meals by name: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search meals"})
		return
	}

	c.JSON(http.StatusOK, meals)
}