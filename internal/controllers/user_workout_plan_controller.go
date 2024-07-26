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

func (uc *UserController) GetActiveWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("Error retrieving userID from context\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	objID, ok := userID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error converting userID from type interface{} to primitive.ObjectID\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert user ID to string"})
		return
	}

	userActiveWorkoutPlan, err := uc.UserService.GetActiveWorkoutPlan(c.Request.Context(), objID);
	if err != nil {
		if errors.Is(err, services.ErrActiveWorkoutPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User has no active workout plan"})
			return
		}

		log.Printf("Error getting active workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active workout plan"})
		return
	}

	c.JSON(http.StatusOK, userActiveWorkoutPlan)	
}

func (uc *UserController) JoinWorkoutPlan(c *gin.Context) {
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

	workoutPlanID, err := primitive.ObjectIDFromHex(c.Param("workoutPlanId"))
	if err != nil {
		log.Printf("Error parsing workout plan ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout plan ID"})
		return
	}

	if err := uc.UserService.JoinWorkoutPlan(c.Request.Context(), objID, workoutPlanID); err != nil {
		if errors.Is(err, services.ErrWorkoutPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
			return
		}
		if errors.Is(err, services.ErrAlreadyJoinded) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User has already joined this workout plan"})
			return
		}

		log.Printf("Error joining workout plan: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join workout plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined workout plan"})
}

func (uc *UserController) CompleteExercise(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("Error retrieving userID from context\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	objID, ok := userID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error converting userID from type interface{} to primitive.ObjectID\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert user ID to string"})
		return
	}

	exerciseID, err := primitive.ObjectIDFromHex(c.Param("exerciseId"))
	if err != nil {
		log.Printf("Error parsing exercise ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	circuitID, err := primitive.ObjectIDFromHex(c.Param("circuitId"))
    if err != nil {
        log.Printf("Error parsing circuit ID: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid circuit ID"})
        return
    }

	var logs []models.UserExerciseLogInput
	if err := c.ShouldBindJSON(&logs); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	err = uc.UserService.MarkExerciseAsCompleted(c.Request.Context(), objID, exerciseID, circuitID, logs)
	if err != nil {
		if errors.Is(err, services.ErrExerciseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User exercise status not found"})
			return
		}

		if errors.Is(err, services.ErrExerciseAlreadyCompleted) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Exercise has already been completed"})
			return
		}

		if errors.Is(err, services.ErrCircuitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User circuit status not found"})
			return
		}

		log.Printf("Error marking user exercise as completed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark user exercise as completed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise marked as completed"})
}

// func(uc *UserController) GetWorkoutPlanProgress(c *gin.Context) {
// 	userID, exists := c.Get("userId")
// 	if !exists {
// 		log.Printf("Error retrieving userID from context\n")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to retrieve user ID from context"})
// 		return
// 	}

// 	objID, ok := userID.(primitive.ObjectID)
// 	if !ok {
// 		log.Printf("Error converting userID from type interface{} to primitive.ObjectID\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert user ID to string"})
// 		return
// 	}

// 	workoutPlanID, err := primitive.ObjectIDFromHex(c.Param("workoutPlanId"))
// 	if err != nil {
// 		log.Printf("Error parsing workout plan ID: %v\n", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout plan ID"})
// 		return
// 	}

// 	progress, err := uc.UserService.GetWorkoutPlanProgress(c.Request.Context(), objID, workoutPlanID)
// 	if err != nil {
// 		if errors.Is(err, services.ErrWorkoutPlanNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "User workout plan status not found"})
// 			return
// 		}

// 		log.Printf("Error getting workout plan progress: %v\n", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workout plan progress"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"progress": progress})
// }



// func (uc *UserController) CompleteExercise(c *gin.Context) {
// 	userID, exits := c.Get("userId")
// 	if !exits {
// 		log.Printf("Error retrieving userID from context\n")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to retrieve user ID from context"})
// 		return
// 	}

// 	objID, ok := userID.(primitive.ObjectID)
// 	if !ok {
// 		log.Printf("Error converting userID from type interface {} to primitive.ObjectID\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert user ID to string"})
// 		return
// 	}

// 	exerciseID, err := primitive.ObjectIDFromHex(c.Param("exerciseId"))
// 	if err != nil {
// 		log.Printf("Error parsing exercise ID: %v\n", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
// 		return
// 	}

// 	//Transaction integrity to ensure that all updates within our cascading update either comple successfully or rollback together avoiding data inconsistency(userExerciseStatus, userCircuitStatus, userWorkoutDayStatus, userWorkoutWeekStatus, userWorkoutPlanStatus)
// 	session, err := uc.UserService.StartSession()
// 	if err != nil {
// 		log.Printf("Error starting session: %v\n", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
// 		return
// 	}
// 	defer session.EndSession(c.Request.Context())

// 	var logs []models.UserExerciseLogInput
// 	if err := c.ShouldBindJSON(&logs); err != nil {
// 		log.Printf("Error parsing JSON: %v\n", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
// 		return
// 	}

// 	err = mongo.WithSession(c.Request.Context(), session, func(sc mongo.SessionContext) error {
// 		if err := session.StartTransaction(); err != nil {
// 			return err
// 		}

// 		if err := uc.UserService.MarkExerciseAsCompleted(sc, objID, exerciseID, logs); err != nil {
// 			return err
// 		}

// 		if err := session.CommitTransaction(sc); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		if errors.Is(err, services.ErrExerciseNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "User exercise not found"})
// 			return
// 		}

// 		log.Printf("Error marking user exercise as completed: %v\n", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark user exercise as completed"})
// 		return
// 	}
	
// 	c.JSON(http.StatusOK, gin.H{"message": "Exercise marked as completed"})
// }