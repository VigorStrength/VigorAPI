package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ac *AdminController) GetUsers(c *gin.Context) {
	users, err := ac.AdminService.GetUsers(c.Request.Context())
	if err != nil {
		log.Printf("Error getting users: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}