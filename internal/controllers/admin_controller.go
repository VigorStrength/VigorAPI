package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	AdminService services.AdminService
	JWTService   utils.TokenService
}

func NewAdminController(adminService services.AdminService, jwtService utils.TokenService) *AdminController {
	return &AdminController{
		AdminService: adminService,
		JWTService:   jwtService,
	}
}

func (ac *AdminController) Register(c *gin.Context) {
	var input models.AdminRegistrationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := validate.Struct(input); err != nil {
		log.Printf("Error validating input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.AdminService.RegisterAdmin(c.Request.Context(), input); err != nil {
		if errors.Is(err, services.ErrAdminAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Admin already exists"})
			return
		}

		log.Printf("Error registering admin: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin registered successfully"})
}

func (ac *AdminController) Login(c *gin.Context) {
	var loginDetails models.LoginDetails
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	admin, err := ac.AdminService.GetAdminByEmail(c.Request.Context(), loginDetails.Email, loginDetails.Password)
	if err != nil {
		log.Printf("Error authenticating admin: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := ac.JWTService.GenerateAccessToken(admin.ID, admin.Email, admin.Role)
	if err != nil {	
		log.Printf("Error generating access token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := ac.JWTService.GenerateRefreshToken(admin.ID, admin.Email, admin.Role)
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