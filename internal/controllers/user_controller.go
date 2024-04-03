package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	UserService services.UserService
	JWTService   utils.TokenService
}

func NewUserController(userService services.UserService, jwtService utils.TokenService) *UserController {
	return &UserController{
		UserService: userService,
		JWTService: jwtService,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	var input models.UserRegistrationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			detailedErrors := make(map[string]string)
			for _, valErr := range validationErrors {
				field := valErr.StructField()
				detailedErrors[field] = valErr.Tag()
			}
			log.Printf("Validation errors: %v\n", detailedErrors)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": detailedErrors})
			return
		} else {
			log.Printf("Error validating input: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
	}

	if err := uc.UserService.RegisterUser(c.Request.Context(), input); err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}

		log.Printf("Error registering user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (uc *UserController) Login(c *gin.Context) {
	var loginDetails models.LoginDetails
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	user, err := uc.UserService.GetUserByEmail(c.Request.Context(), loginDetails.Email, loginDetails.Password)
	if err != nil {
		log.Printf("Error authenticating user: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	
	accessToken, err := uc.JWTService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating access token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := uc.JWTService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating refresh token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}