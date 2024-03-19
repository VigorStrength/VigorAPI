package middlewares_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GhostDrew11/vigor-api/internal/middlewares"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockJWTService struct {
	mock.Mock
	utils.JWTService
}

func (m *MockJWTService) GenerateAccessToken(userId primitive.ObjectID, email string) (string, error) {
	args := m.Called(userId, email)
	// tokenStr, _ := args.String(0)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userId primitive.ObjectID, email string) (string, error) {
	args := m.Called(userId, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) VerifyToken(tokenString string) (*utils.Claims, error) {
	args := m.Called(tokenString)
	claims, _ := args.Get(0).(*utils.Claims)
	return claims, args.Error(1)
}

func TestAuthenticateMiddlewareSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockJWTService := new(MockJWTService)
	tokenString := "dummyToken"
	userId := primitive.NewObjectID()
	email := "test@example.com"

	claims := &utils.Claims{
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	mockJWTService.On("VerifyToken", tokenString).Return(claims, nil)

	router.Use(middlewares.Authenticate(mockJWTService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Passed"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", tokenString)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockJWTService.AssertExpectations(t)
}

func TestAuthenticateMiddlewareWithEmptyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockJWTService := new(MockJWTService)

	router.Use(middlewares.Authenticate(mockJWTService))
	router.GET("/test", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Should not get here"})
    })

	w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticateMiddlewareFailureTokenVerification(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockJWTService := new(MockJWTService)
	invalidTokenString := "invalidtokenstring"

	mockJWTService.On("VerifyToken", invalidTokenString).Return(nil, utils.ErrInvalidToken)

	router.Use(middlewares.Authenticate(mockJWTService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Should not reach"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", invalidTokenString)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockJWTService.AssertExpectations(t)
}

func TestRefreshHandlerMiddlewareSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockJWTService := new(MockJWTService)
	userId := primitive.NewObjectID()
	email := "test@example.com"

	refreshToken := "refreshTokenDummy"
	newAccessToken := "newAccessTokenDummy"
	newRefreshToken := "newRefreshTokenDummy"

	claims := &utils.Claims{
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	mockJWTService.On("VerifyToken", refreshToken).Return(claims, nil)
	mockJWTService.On("GenerateAccessToken", userId, email).Return(newAccessToken, nil)
	mockJWTService.On("GenerateRefreshToken", userId, email).Return(newRefreshToken, nil)

	router.POST("/refresh", middlewares.RefreshHandler(mockJWTService))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Refresh-Token", refreshToken)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), newAccessToken)
	assert.Contains(t, w.Body.String(), newRefreshToken)
	mockJWTService.AssertExpectations(t)
}

func TestRefreshHandlerMiddlewareFailureEmptyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockJWTService := new(MockJWTService)

	router.POST("/refresh", middlewares.RefreshHandler(mockJWTService))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code) 
}

func TestRefreshHandlerMiddlewareFailureTokenVerification(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()

    mockJWTService := new(MockJWTService)
    invalidRefreshToken := "invalidRefreshToken"

    mockJWTService.On("VerifyToken", invalidRefreshToken).Return(nil, utils.ErrInvalidToken)

    router.POST("/refresh", middlewares.RefreshHandler(mockJWTService))

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/refresh", nil)
    req.Header.Set("Refresh-Token", invalidRefreshToken)

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code) 
    mockJWTService.AssertExpectations(t)
}

func TestRefreshHandlerMiddlewareFailureTokenGeneration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockJWTService := new(MockJWTService)
	userId := primitive.NewObjectID()
	email := "test@example.com"

	refreshToken := "refreshTokenDummy"
	claims := &utils.Claims{
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	mockJWTService.On("VerifyToken", refreshToken).Return(claims, nil)
	mockJWTService.On("GenerateAccessToken", userId, email).Return("", errors.New("Unable to generate access token"))

	router.POST("/refresh", middlewares.RefreshHandler(mockJWTService))

	w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/refresh", nil)
    req.Header.Set("Refresh-Token", refreshToken)

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusInternalServerError, w.Code) 
    mockJWTService.AssertExpectations(t)
}

