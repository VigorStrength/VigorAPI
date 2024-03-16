package utils

import "github.com/golang-jwt/jwt/v5"

type HashPasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type TokenService interface {
	GenerateAllTokens() (string, string, error)
	VerifyToken(tokenString string) (*Claims, error)
}

type JWTHandler interface {
	ParseWithClaims(tokenString string, claims *Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
}