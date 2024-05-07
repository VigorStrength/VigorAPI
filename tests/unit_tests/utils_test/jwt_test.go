package u

import (
	"testing"
	"time"

	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockJWTHandler struct {
	mock.Mock
}

func (m *MockJWTHandler) ParseWithClaims(tokenString string, claims *utils.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
    args := m.Called(tokenString, claims, keyFunc)
	token, _ := args.Get(0).(*jwt.Token)
    return token, args.Error(1)
}

func TestGenerateToken(t *testing.T) {
	signingMethod := jwt.SigningMethodHS256
	tokenClaims := &utils.Claims{
		UserId: primitive.NewObjectID(),
		Email: "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	secretKey := []byte("lilsecret")

	token, err := utils.GenerateToken(signingMethod, *tokenClaims, secretKey)
	assert.Nil(t, err, "Generating token should not return an error")
	assert.NotNil(t, token, "The generated token should not be nil")
}

func TestGenerateTokenFailure(t *testing.T) {
	signingMethod := jwt.SigningMethodNone
	tokenClaims := &utils.Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	secretKey := []byte("supersecret")
	token, err := utils.GenerateToken(signingMethod, *tokenClaims, secretKey)
	assert.NotNil(t, err, "Generating token should return an error")
	assert.Empty(t, token, "The token string must be an empty string")
}

// Same as GenerateRefresh
// Failure case really not needed cause the main
// function is calling GenerateToken for which the failure case is already tested above
func TestGenerateAccessTokenSuccess(t *testing.T) {
	userId := primitive.NewObjectID()
    email := "test@example.com"
	role := "admin"
    jwtSecretKey := "supersecret"
    mockHandler := new(MockJWTHandler)

    jwtService := utils.NewJWTService(jwtSecretKey, mockHandler)

    accessTokenStr, err := jwtService.GenerateAccessToken(userId, email, role)

    assert.NoError(t, err, "Generating access token should not produce an error")
    assert.NotEmpty(t, accessTokenStr, "Access token should not be empty")
}

func TestVerifyTokenSuccess(t *testing.T) {
	userId := primitive.NewObjectID()
	email := "test@example.com"
	role := "admin"
	jwtSecretKey := "supersecret"
	mockHandler := new(MockJWTHandler)

	mockToken := &jwt.Token{Valid: true}
	mockHandler.On("ParseWithClaims",
    mock.AnythingOfType("string"),
    mock.AnythingOfType("*utils.Claims"),
    mock.AnythingOfType("jwt.Keyfunc"),
).Run(func(args mock.Arguments) {
    claims := args.Get(1).(*utils.Claims)
    claims.UserId = userId
    claims.Email = email
	claims.Role = role
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
}).Return(mockToken, nil)

	jwtService := utils.NewJWTService(jwtSecretKey, mockHandler)

	accessTokenStr, err := jwtService.GenerateAccessToken(userId, email, role)
	assert.NoError(t, err, "Generating all tokens should not produce an error")

	accessTokenClaims, err := jwtService.VerifyToken(accessTokenStr)
	assert.Nil(t, err, "Verifying access token should not produce an error")

	assert.Equal(t, userId, accessTokenClaims.UserId, "UserId should match")
	assert.Equal(t, email, accessTokenClaims.Email, "Email should match")

	mockHandler.AssertExpectations(t)
}

func TestVerifyTokenFailureInvalidSigningMethod(t *testing.T) {
    jwtSecretKey := "secret"
	mockHandler := new(MockJWTHandler)

	mockHandler.On("ParseWithClaims", mock.AnythingOfType("string"), mock.AnythingOfType("*utils.Claims"), mock.AnythingOfType("jwt.Keyfunc")).Return(nil, utils.ErrInvalidSigningMethod)

	jwtService := utils.NewJWTService(jwtSecretKey,mockHandler)

	_, err := jwtService.VerifyToken("dummyToken")
	assert.Error(t, err, "Verifying wrongly signed token should return an error")
	assert.Equal(t, utils.ErrInvalidSigningMethod, err)

	mockHandler.AssertExpectations(t)
}

func TestVerifyTokenFailureInvalid(t *testing.T) {
    jwtSecretKey := "secret"
	mockHandler := new(MockJWTHandler)

	mockToken := &jwt.Token{Valid: false}
	mockHandler.On("ParseWithClaims", mock.AnythingOfType("string"), mock.AnythingOfType("*utils.Claims"), mock.AnythingOfType("jwt.Keyfunc")).Return(mockToken, utils.ErrInvalidToken)

	jwtService := utils.NewJWTService(jwtSecretKey, mockHandler)

	_, err := jwtService.VerifyToken("tamperedToken")
	assert.Error(t, err, "Veryfying invalid token should return an error")
	assert.Equal(t, utils.ErrInvalidToken, err)
	// assert.Contains(t, utils.ErrInvalidToken, err.Error())
}

func TestVerifyTokenFailureExpired(t *testing.T) {
	jwtSecretKey := "supersecret"
	mockHandler := new(MockJWTHandler)

	mockHandler.On("ParseWithClaims", mock.AnythingOfType("string"), mock.AnythingOfType("*utils.Claims"), mock.AnythingOfType("jwt.Keyfunc")).Run(func(args mock.Arguments) {
		claims := args.Get(1).(*utils.Claims)
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))
	}).Return(nil, utils.ErrTokenExpired)

	jwtService := utils.NewJWTService(jwtSecretKey, mockHandler)

	_, err := jwtService.VerifyToken("dummytoken")
	assert.Error(t, err, "Veryfying expired access token should return an error")
	assert.Equal(t, utils.ErrTokenExpired, err)

	mockHandler.AssertExpectations(t)
}




