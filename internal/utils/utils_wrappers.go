package utils

import "github.com/golang-jwt/jwt/v5"

type DefaultJWTHandler struct {}

func (h *DefaultJWTHandler) ParseWithClaims(tokenString string, claims *Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, keyFunc)
}