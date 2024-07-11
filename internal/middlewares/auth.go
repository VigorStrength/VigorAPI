package middlewares

import (
	"net/http"

	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func RequireRole(ts utils.TokenService, requiredRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		// Check if Authorization header is present
		authorizationHeader := ctx.GetHeader("Authorization")
		if authorizationHeader != "" { 
			token = authorizationHeader[7:] // Remove "Bearer " from the token
		} else {
			// Check if Refresh-Token header is present
			refreshTokenHeader := ctx.GetHeader("Refresh-Token")
			if refreshTokenHeader != "" {
				token := refreshTokenHeader
				if len(token) > 7 && token[:7] == "Bearer " {
					token = token[7:]
				}
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				return
			}
		}

		claims, err := ts.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Check if the user has the required role
		roleIsAllowed := false
		for _, role := range requiredRoles {
			if claims.Role == role {
				roleIsAllowed = true
				break
			}
		}

		if !roleIsAllowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			return
		}

		// Set user info and role in the context
		ctx.Set("userId", claims.UserId)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)
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

		newRefreshToken, err := ts.GenerateRefreshToken(claims.UserId, claims.Email, claims.Role)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate new refresh token"})
			return
		}

		// Return new refresh toke to the client
		ctx.JSON(http.StatusOK, gin.H{
			"refreshToken": newRefreshToken,
		})
	}
}
