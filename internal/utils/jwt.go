package utils

import (
	"errors"
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
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(jwtSecretKey))
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
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func VerifyToken(tokenString string, jwtSecretkey string) (*Claims, error) {
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	//Parse the JWT string and store the result in `claims`
	// Note: the function returns the key for validating the token's signature
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token's algorithm corresponds to "SigningMethodMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(jwtSecretkey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token is expired ")
	}

	return claims, nil
}

// Code for using ECDSA encryption instead of HMAC encryption in case I need to switch
// func GenerateToken(userID, email string) (accessTokenStr, refreshTokenStr string, err error) {
//     privateKeyEnv := os.Getenv("ECDSA_PRIVATE_KEY")
//     if privateKeyEnv == "" {
//         return "", "", errors.New("ECDSA private key not set in environment variables")
//     }

//     privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyEnv)
//     if err != nil {
//         return "", "", err
//     }

//     privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
//     if err != nil {
//         return "", "", err
//     }

//     // Access Token
//     accessTokenClaims := &Claims{
//         UserID: userID,
//         Email:  email,
//         RegisteredClaims: jwt.RegisteredClaims{
//             ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Short lifespan for access token
//         },
//     }

//     accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, accessTokenClaims)
//     accessTokenStr, err = accessToken.SignedString(privateKey)
//     if err != nil {
//         return "", "", err
//     }

//     // Refresh Token
//     refreshTokenClaims := &Claims{
//         UserID: userID,
//         Email:  email,
//         RegisteredClaims: jwt.RegisteredClaims{
//             ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Longer lifespan for refresh token
//         },
//     }

//     refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, refreshTokenClaims)
//     refreshTokenStr, err = refreshToken.SignedString(privateKey)
//     if err != nil {
//         return "", "", err
//     }

//     return accessTokenStr, refreshTokenStr, nil
// }

// func VerifyToken(tokenString string) (*Claims, error) {
//     publicKeyEnv := os.Getenv("ECDSA_PUBLIC_KEY")
//     if publicKeyEnv == "" {
//         return nil, errors.New("ECDSA public key not set in environment variables")
//     }

//     publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyEnv)
//     if err != nil {
//         return nil, err
//     }

//     publicKey, err := jwt.ParseECPublicKeyFromPEM(publicKeyBytes)
//     if err != nil {
//         return nil, err
//     }

//     token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//         return publicKey, nil
//     })
//     if err != nil {
//         return nil, err
//     }

//     if claims, ok := token.Claims.(*Claims); ok && token.Valid {
//         return claims, nil
//     } else {
//         return nil, errors.New("invalid token")
//     }
// }