package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token is expired")
)

type Claims struct {
    jwt.RegisteredClaims
    UserId primitive.ObjectID `json:"userId"`
    Email  string             `json:"email"`
	Role   string             `json:"role"`
}

type JWTService struct {
    jwtSecretKey  []byte
    signingMethod jwt.SigningMethod
    handler       JWTHandler
}

func NewJWTService(key string, handler JWTHandler) *JWTService {
    return &JWTService{
        jwtSecretKey:  []byte(key),
        signingMethod: jwt.SigningMethodHS256, // Keep as configurable if needed
        handler:       handler,
    }
}

func (j *JWTService) GenerateAccessToken(userId primitive.ObjectID, email, role string) (string, error) {
    accessTokenExp := time.Now().Add(1 * time.Hour) // Or use a configuration
    claims := Claims{
        UserId: userId,
        Email:  email,
		Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(accessTokenExp),
        },
    }
    return GenerateToken(j.signingMethod, claims, j.jwtSecretKey)
}

func (j *JWTService) GenerateRefreshToken(userId primitive.ObjectID, email, role string) (string, error) {
    refreshTokenExp := time.Now().Add(168 * time.Hour) // Or use a configuration
    claims := Claims{
        UserId: userId,
        Email:  email,
		Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
        },
    }
    return GenerateToken(j.signingMethod, claims, j.jwtSecretKey)
}

func (j *JWTService) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := j.handler.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return j.jwtSecretKey, nil
	})

	if token == nil || err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrTokenExpired
	}

	return claims, nil
}

func GenerateToken(tokenSigningMethod jwt.SigningMethod, tokenClaims Claims, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(tokenSigningMethod, tokenClaims)
	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	
	return tokenStr, nil
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