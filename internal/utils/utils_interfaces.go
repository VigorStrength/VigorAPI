package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HashPasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// var _ TokenService = (*JWTService)(nil)
type TokenService interface {
    GenerateAccessToken(userId primitive.ObjectID, email, role string) (string, error)
    GenerateRefreshToken(userId primitive.ObjectID, email, role string) (string, error)
    VerifyToken(tokenString string) (*Claims, error)
}

type ParserService interface {
	StructToBson(v interface{}) bson.M
}

type JWTHandler interface {
	ParseWithClaims(tokenString string, claims *Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
}