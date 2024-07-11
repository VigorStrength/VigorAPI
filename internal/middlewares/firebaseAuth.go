package middlewares

import (
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func FirebaseRequireRole(client *auth.Client, requiredRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return 
		}

		idToken := authHeader[7:] // Remove "Bearer " from the token
		verifiedIDToken, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Check if the user has the required role
		roleIsAllowed := false
		for _, role := range requiredRoles {
			if verifiedIDToken.Claims["role"] == role {
				roleIsAllowed = true
				break
			}
		}

		if !roleIsAllowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			return
		}

		// Set user info and role in the context
		ctx.Set("userId", verifiedIDToken.UID)
		ctx.Set("email", verifiedIDToken.Claims["email"])
		ctx.Set("role", verifiedIDToken.Claims["role"])
		ctx.Next()
	}
 }