package middlewares

import (
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate( jwtSecretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		claims, err := utils.VerifyToken(token, jwtSecretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}

// RefreshHandler to handle refresh token logic
func RefreshHandler(ctx *gin.Context, jwtSecretKey string) {
	// Assuming the refresh token is sent via headers or cookies
	refreshToken := ctx.GetHeader("Refresh-Token")
	if refreshToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No refresh token provided"})
		return
	}

	// Verify the refresh token
	claims, err := utils.VerifyToken(refreshToken, jwtSecretKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
		return
	}

	// Generate new access and refresh tokens
	newAccessToken, newRefreshToken, err := utils.GenerateToken(claims.UserId, claims.Email, jwtSecretKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate new tokens"})
		return
	}

	// Return new tokens to the client
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": newAccessToken,
		"refreshToken": newRefreshToken,
	})
}