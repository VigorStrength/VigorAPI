package s

import (
	"context"
	"errors"
	"testing"

	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

func TestRegisterUserSuccess(t *testing.T){
	ctx := context.Background()

	mockDB := new(MockMongoDatabase)
	mockCollection := new(MockMongoCollection)
	hasher := &utils.DefaultHasher{}
	parser := &utils.DefaultParser{}
	userService := services.NewUserService(mockDB, hasher, parser)
	input := new(models.UserRegistrationInput)
	filter := bson.M{"email": input.Email}
	usernameFilter := bson.M{"profileInformation.username": input.ProfileInformation.Username}

	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)
	mockCollection.On("CountDocuments", ctx, filter).Return(int64(0), nil)
	mockCollection.On("CountDocuments", ctx, usernameFilter).Return(int64(0), nil)
	mockCollection.On("InsertOne", ctx, mock.AnythingOfType("models.User")).Return(*new(db.MongoInsertOneResult), nil)

	err := userService.RegisterUser(ctx, *input)

	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
}

func TestRegisterUserFailure_CheckingUser(t *testing.T){
	ctx := context.Background()
	input := new(models.UserRegistrationInput)
	filter := bson.M{"email": input.Email}
	mockDB := new(MockMongoDatabase)
	mockCollection := new(MockMongoCollection)
	mockHasher := new(MockHasher)
	mockParser := new(MockParser)

	userService := services.NewUserService(mockDB, mockHasher, mockParser)

	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)
	mockCollection.On("CountDocuments", ctx, filter).Return(int64(0), errors.New("error checking if user already exists"))

	err := userService.RegisterUser(ctx, *input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if user already exists")
	mockCollection.AssertExpectations(t)
}

func TestRegisterUserFailure_UserAlreadyExists(t *testing.T){
	ctx := context.Background()
	input := new(models.UserRegistrationInput)
	filter := bson.M{"email": input.Email}
	mockDB := new(MockMongoDatabase)
	mockCollection := new(MockMongoCollection)
	mockHasher := new(MockHasher)
	mockParser := new(MockParser)

	userService := services.NewUserService(mockDB, mockHasher, mockParser)

	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)
	mockCollection.On("CountDocuments", ctx, filter).Return(int64(1), nil)

	err := userService.RegisterUser(ctx, *input)
	assert.Error(t, err)
	assert.Equal(t, services.ErrUserAlreadyExists, err)
	mockCollection.AssertExpectations(t)
}

func TestRegisterUserFailure_InsertingOne(t *testing.T){
	ctx := context.Background()
	input := new(models.UserRegistrationInput)
	filter := bson.M{"email": input.Email}
	usernameFilter := bson.M{"profileInformation.username": input.ProfileInformation.Username}
	mockDB := new(MockMongoDatabase)
	mockCollection := new(MockMongoCollection)
	mockHasher := new(MockHasher)
	mockParser := new(MockParser)

	userService := services.NewUserService(mockDB, mockHasher, mockParser)

	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)
	mockCollection.On("CountDocuments", ctx, filter).Return(int64(0), nil)
	mockCollection.On("CountDocuments", ctx, usernameFilter).Return(int64(0), nil)
	mockCollection.On("InsertOne", ctx, mock.AnythingOfType("models.User")).Return(*new(db.MongoInsertOneResult), errors.New("error inserting user"))

	err := userService.RegisterUser(ctx, *input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error inserting user")
	mockCollection.AssertExpectations(t)
}

func TestGetUserByEmailSuccess(t *testing.T){
	ctx := context.Background()

	mockDB := new(MockMongoDatabase)
	mockCollection := new(MockMongoCollection)
	mockMongoSingleResult := new(MockMongoSingleResult)
	mockHasher := new(MockHasher)
	mockParser := new(MockParser)
	userService := services.NewUserService(mockDB, mockHasher, mockParser)

	email := "testuser@example.com"
	password := "securepassword"
	filter := bson.M{"email": email}

	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)
	mockCollection.On("FindOne", ctx, filter, mock.AnythingOfType("[]*options.FindOneOptions")).Return(mockMongoSingleResult)
	mockMongoSingleResult.On("Err").Return(nil)
	mockMongoSingleResult.On("Decode", mock.AnythingOfType("*models.User")).Return(nil)
	mockHasher.On("CheckPasswordHash", password, mock.AnythingOfType("string")).Return(true)

	_, err := userService.GetUserByEmail(ctx, email, password)

	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
	mockMongoSingleResult.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestGetUserByEmailFailure_FecthingUser(t *testing.T){
	ctx := context.Background()

	mockDB := new(MockMongoDatabase)
	mockCollection := new(MockMongoCollection)
	mockMongoSingleResult := new(MockMongoSingleResult)
	mockHasher := new(MockHasher)
	mockParser := new(MockParser)
	userService := services.NewUserService(mockDB, mockHasher, mockParser)

	email := "testuser@example.com"
	password := "securepassword"
	filter := bson.M{"email": email}

	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)
	mockCollection.On("FindOne", ctx, filter, mock.AnythingOfType("[]*options.FindOneOptions")).Return(mockMongoSingleResult)
	mockMongoSingleResult.On("Err").Return(errors.New("an error"))

	_, err := userService.GetUserByEmail(ctx, email, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching user with email")
	mockCollection.AssertExpectations(t)
	mockMongoSingleResult.AssertExpectations(t)
}

func TestGetUserByEmailFailure_DecodingResult(t *testing.T){
	ctx := context.Background()

	mockDB := new(MockMongoDatabase)
	mockCollection := new(MockMongoCollection)
	mockMongoSingleResult := new(MockMongoSingleResult)
	mockHasher := new(MockHasher)
	mockParser := new(MockParser)
	userService := services.NewUserService(mockDB, mockHasher, mockParser)

	email := "testuser@example.com"
	password := "securepassword"
	filter := bson.M{"email": email}

	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)
	mockCollection.On("FindOne", ctx, filter, mock.AnythingOfType("[]*options.FindOneOptions")).Return(mockMongoSingleResult)
	mockMongoSingleResult.On("Err").Return(nil)
	mockMongoSingleResult.On("Decode", mock.AnythingOfType("*models.User")).Return(errors.New("another error"))

	_, err := userService.GetUserByEmail(ctx, email, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error decoding user data for email")
	mockCollection.AssertExpectations(t)
	mockMongoSingleResult.AssertExpectations(t)
}

func TestGetUSerByEmailFailure_CheckingPassword(t *testing.T){
	ctx := context.Background()

	mockDB := new(MockMongoDatabase)
	mockCollection := new(MockMongoCollection)
	mockMongoSingleResult := new(MockMongoSingleResult)
	mockHasher := new(MockHasher)
	mockParser := new(MockParser)
	userService := services.NewUserService(mockDB, mockHasher, mockParser)

	email := "testuser@example.com"
	password := "securepassword"
	filter := bson.M{"email": email}

	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)
	mockCollection.On("FindOne", ctx, filter, mock.AnythingOfType("[]*options.FindOneOptions")).Return(mockMongoSingleResult)
	mockMongoSingleResult.On("Err").Return(nil)
	mockMongoSingleResult.On("Decode", mock.AnythingOfType("*models.User")).Return(nil)
	mockHasher.On("CheckPasswordHash", password, mock.AnythingOfType("string")).Return(false)

	_, err := userService.GetUserByEmail(ctx, email, password)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidUserCredentials, err)
	mockCollection.AssertExpectations(t)
	mockMongoSingleResult.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
