package controllers

import (
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/gin-gonic/gin"
)

type UserFireBaseController struct {
	UserService services.UserService
	FirebaseAuth *auth.Client
}

func NewUserFireBaseController(userService services.UserService, firebaseAuth *auth.Client) *UserFireBaseController {
	return &UserFireBaseController{
		UserService: userService,
		FirebaseAuth: firebaseAuth,
	}
}

func (ufc *UserFireBaseController) RegisterUser() {
	// Register user using Firebase Auth
}

func (ufc *UserFireBaseController) LoginUser(c *gin.Context) {
	var loginDetails models.UserFirebaseLoginDetails
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	token, err := ufc.FirebaseAuth.VerifyIDToken(c, loginDetails.IDToken)
	if err != nil {
		log.Printf("Error verifying ID token: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Extract user info from the token
	email := token.Claims["email"].(string)
	password := token.Claims["password"].(string)

	// Find the user in the database
	_, err = ufc.UserService.GetUserByEmail(c, email, password)
	if err != nil {
		log.Printf("Error finding user: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the firebase ID token
	c.JSON(http.StatusOK, gin.H{"firebaseToken": loginDetails.IDToken})
}