package controllers

import (
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
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		log.Panic(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	if err := ac.AdminService.RegisterAdmin(c.Request.Context(), admin); err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin registered successfully"})
}

func (ac *AdminController) Login(c *gin.Context) {
	var loginDetails models.LoginDetails
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		log.Panic(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	admin, err := ac.AdminService.GetAdminByEmail(c.Request.Context(), loginDetails.Email)
	if err != nil || !utils.CheckPasswordHash(loginDetails.Password, admin.PasswordHash) {
		log.Panic(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := ac.JWTService.GenerateAccessToken(admin.ID, admin.Email)
	if err != nil {	
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := ac.JWTService.GenerateRefreshToken(admin.ID, admin.Email)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}