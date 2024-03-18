package middlewares

import (
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(ts utils.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		claims, err := ts.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}

func RefreshHandler(ts utils.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken := ctx.GetHeader("Refresh-Token")
		if refreshToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No refresh token provided"})
			return
		}

		claims, err := ts.VerifyToken(refreshToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
			return
		}

		newAccessToken, err := ts.GenerateAccessToken(claims.UserId, claims.Email)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate new access token"})
			return
		}

		newRefreshToken, err := ts.GenerateRefreshToken(claims.UserId, claims.Email)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate new refresh token"})
			return
		}

		// Return new tokens to the client
		ctx.JSON(http.StatusOK, gin.H{
			"accessToken": newAccessToken,
			"refreshToken": newRefreshToken,
		})
	}
}
