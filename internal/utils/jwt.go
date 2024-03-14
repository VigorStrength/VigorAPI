package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	UserId primitive.ObjectID `json:"userId"`
	Email string 			  `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userId primitive.ObjectID, email , jwtSecretKey string) (string, string, error) {
	// Expiration Times
	accesTokenExp := time.Now().Add(1 * time.Hour)
	refreshTokenExp := time.Now().Add(24 * time.Hour)

	// Access Token
	accessTokenClaims := &Claims{
		UserId: userId,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accesTokenExp),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, accessTokenClaims)
	accessTokenStr, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshTokenClaims := &Claims{
		UserId: userId,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, refreshTokenClaims)
	refreshTokenStr, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}
