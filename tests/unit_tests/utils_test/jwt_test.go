package utils_test

import (
	"testing"
	"time"

	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	userId := primitive.NewObjectID()
	email := "test@example.com"
	jwtSecretkey := "lilsecret"

	accessToken, refreshToken, err := utils.GenerateToken(userId, email, jwtSecretkey)
	assert.Nil(t, err, "Generating tokens should not produce an error")
	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	assert.NotEmpty(t, refreshToken, "Refresh token should not be empty")

	// Verify Access Token
	verifiedClaims, err := utils.VerifyToken(accessToken, jwtSecretkey)
	assert.Nil(t, err, "Verifying access token should not produce an error")
	assert.Equal(t, userId, verifiedClaims.UserId, "User ID should match")
	assert.Equal(t, email, verifiedClaims.Email, "Email should match")
}

func TestVerifyUnexpectedSigningMethod(t *testing.T) {

	claims := &jwt.RegisteredClaims{}
    token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
    tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.Nil(t, err, "Creating unsigned token should not produce an error")

	_, err = utils.VerifyToken(tokenString, "anySecretKey")
    assert.NotNil(t, err, "Expected an error for unexpected signing method")
	assert.Contains(t, err.Error(), "unexpected signing method", "Error message should indicate unexpected signing method")
}

func TestVerifyTokenExpired(t *testing.T) {
	userId := primitive.NewObjectID()
	email := "test@example.com"
	jwtSecretKey := "secret"

	// Generate a token with a past expiration
	accessToken, _, err := GenerateTokenWithCustomExpiration(userId, email, jwtSecretKey, -1*time.Hour)
	assert.Nil(t, err, "Generating token should not produce an error")

	_, err = utils.VerifyToken(accessToken, jwtSecretKey)
	assert.NotNil(t, err, "Expected an error for expired token")
	assert.Contains(t, err.Error(), "token is expired", "Error message should indicate token expiration")
}

// GenerateTokenWithCustomExpiration is a modified version of GenerateToken for testing
// It allows setting custom expiration times to simulate different conditions
func GenerateTokenWithCustomExpiration(userId primitive.ObjectID, email, jwtSecretKey string, expDuration time.Duration) (string, string, error) {
	// Custom Expiration Time
	accesTokenExp := time.Now().Add(expDuration)

	// Access Token with custom expiration
	accessTokenClaims := &utils.Claims{
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

	// This example only modifies the access token for simplicity
	return accessTokenStr, "", nil
}

// func TestVerifyTokenInvalidSignature(t *testing.T) {
//     userId := primitive.NewObjectID()
//     email := "test@example.com"
//     jwtSecretKey := "secret"

//     // Generate a valid token
//     accessToken, _, err := utils.GenerateToken(userId, email, jwtSecretKey)
//     assert.Nil(t, err, "Generating token should not produce an error")

//     // Simulate an invalid signature by altering the token's signature part
//     parts := strings.Split(accessToken, ".")
//     if len(parts) != 3 {
//         t.Fatalf("Generated token has an unexpected format")
//     }
//     // Alter the signature (e.g., by changing the last character)
//     if parts[2][len(parts[2])-1] == 'a' {
//         parts[2] = parts[2][:len(parts[2])-1] + "b"
//     } else {
//         parts[2] = parts[2][:len(parts[2])-1] + "a"
//     }
//     alteredToken := strings.Join(parts, ".")

//     // Verify the altered token
//     _, err = utils.VerifyToken(alteredToken, jwtSecretKey)
//     assert.NotNil(t, err, "Expected an error due to invalid signature")
//     assert.Contains(t, err.Error(), "invalid token", "Error message should indicate invalid token")
// }

