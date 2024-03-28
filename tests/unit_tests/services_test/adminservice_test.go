package s_test

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

type MockMongoCollection struct {
	mock.Mock
	db.MongoCollection
}

type MockMongoSingleResult struct {
	mock.Mock
	db.MongoSingleResult
}

type MockHasher struct {
	mock.Mock
	utils.DefaultHasher
}

func (mc *MockMongoCollection) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	args := mc.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (mc *MockMongoCollection) InsertOne(ctx context.Context, document interface{}) (db.MongoInsertOneResult, error) {
	args := mc.Called(ctx, document)
	res, _ := args.Get(0).(db.MongoInsertOneResult)
	return res, args.Error(1)
}

func (mc *MockMongoCollection) FindOne(ctx context.Context, filter interface{}) db.MongoSingleResult {
	args := mc.Called(ctx, filter)
	res, _ := args.Get(0).(db.MongoSingleResult)
	return res
}

func (msr *MockMongoSingleResult) Err() error {
	args := msr.Called()
	return args.Error(0)
}

func (msr *MockMongoSingleResult) Decode(v interface{}) error {
	args := msr.Called(v)
	return args.Error(0)
}

func (mh *MockHasher) CheckPasswordHash(password, hash string) bool {
	args := mh.Called(password, hash)
	return args.Bool(0)
}

func TestRegisterAdminSuccess(t *testing.T){
	ctx := context.Background()
	 
	mockCollection := new(MockMongoCollection)
	hasher := &utils.DefaultHasher{}
	adminService := services.NewAdminService(mockCollection, hasher)

	input := models.AdminRegistrationInput{
		Email: "admin@vigor.com",
		Password: "securepassword",
	}
	filter := bson.M{"email": input.Email}

	mockCollection.On("CountDocuments", ctx, filter).Return(int64(0), nil)
	mockCollection.On("InsertOne", ctx, mock.AnythingOfType("models.Admin")).Return(nil, nil)

	err := adminService.RegisterAdmin(ctx, input)

	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
}

func TestRegisterAdminFailure_CheckingAdmin(t *testing.T){
	ctx := context.Background()
	input := models.AdminRegistrationInput{
		Email: "admin@vigor.com",
		Password: "securepassword",
	}
	filter := bson.M{"email": input.Email}
	mockCollection := new(MockMongoCollection)
	mockHasher := new(MockHasher)

	adminService := services.NewAdminService(mockCollection, mockHasher)

	mockCollection.On("CountDocuments", ctx, filter).Return(int64(0), errors.New("error checking if admin already exists"))

	err := adminService.RegisterAdmin(ctx, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if admin already exists")
	mockCollection.AssertExpectations(t)
}

func TestRegisterAdminFailure_AdminAlreadyExists(t *testing.T){
	ctx := context.Background()
	input := models.AdminRegistrationInput{
		Email: "admin@vigor.com",
		Password: "securepassword",
	}
	filter := bson.M{"email": input.Email}
	mockCollection := new(MockMongoCollection)
	mockHasher := new(MockHasher)

	adminService := services.NewAdminService(mockCollection, mockHasher)

	mockCollection.On("CountDocuments", ctx, filter).Return(int64(1), nil)

	err := adminService.RegisterAdmin(ctx, input)
	assert.Error(t, err)
	assert.Equal(t, services.ErrAdminAlreadyExists, err)
	mockCollection.AssertExpectations(t)
}

func TestRegisterAdminFailure_InsertingOne(t *testing.T){
	ctx := context.Background()
	input := models.AdminRegistrationInput{
		Email: "admin@vigor.com",
		Password: "securepassword",
	}
	filter := bson.M{"email": input.Email}
	mockCollection := new(MockMongoCollection)
	mockHasher := new(MockHasher)

	adminService := services.NewAdminService(mockCollection, mockHasher)

	mockCollection.On("CountDocuments", ctx, filter).Return(int64(0), nil)
	mockCollection.On("InsertOne", ctx, mock.AnythingOfType("models.Admin")).Return(nil, errors.New("error inserting admin"))

	err := adminService.RegisterAdmin(ctx, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error inserting admin")
	mockCollection.AssertExpectations(t)
}

func TestGetAdminByEmailSuccess(t *testing.T){
	ctx := context.Background()
	 
	mockCollection := new(MockMongoCollection)
	mockMongoSingleResult := new(MockMongoSingleResult)
	mockHasher := new(MockHasher)
	adminService := services.NewAdminService(mockCollection, mockHasher)

	email := "testadmin@vigor.com"
	password := "securepassword"
	filter := bson.M{"email": email}

	mockCollection.On("FindOne", ctx, filter).Return(mockMongoSingleResult)
	mockMongoSingleResult.On("Err").Return(nil)
	mockMongoSingleResult.On("Decode", mock.AnythingOfType("*models.Admin")).Return(nil)
	mockHasher.On("CheckPasswordHash", password, mock.AnythingOfType("string")).Return(true)


	_, err := adminService.GetAdminByEmail(ctx, email, password)

	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
	mockMongoSingleResult.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestGetAdminByEmailFailure_FindingAdmin(t *testing.T){
	ctx := context.Background()
	 
	mockCollection := new(MockMongoCollection)
	mockMongoSingleResult := new(MockMongoSingleResult)
	mockHasher := new(MockHasher)
	adminService := services.NewAdminService(mockCollection, mockHasher)

	email := "testadmin@vigor.com"
	password := "securepassword"
	filter := bson.M{"email": email}

	mockCollection.On("FindOne", ctx, filter).Return(mockMongoSingleResult)
	mockMongoSingleResult.On("Err").Return(errors.New("an error"))

	_, err := adminService.GetAdminByEmail(ctx, email, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching admin with email")
	mockCollection.AssertExpectations(t)
	mockMongoSingleResult.AssertExpectations(t)
}

func TestGetAdminByEmailFailure_DecodingResult(t *testing.T){
	ctx := context.Background()
	 
	mockCollection := new(MockMongoCollection)
	mockMongoSingleResult := new(MockMongoSingleResult)
	mockHasher := new(MockHasher)
	adminService := services.NewAdminService(mockCollection, mockHasher)

	email := "testadmin@vigor.com"
	password := "securepassword"
	filter := bson.M{"email": email}

	mockCollection.On("FindOne", ctx, filter).Return(mockMongoSingleResult)
	mockMongoSingleResult.On("Err").Return(nil)
	mockMongoSingleResult.On("Decode", mock.AnythingOfType("*models.Admin")).Return(errors.New("another error"))

	_, err := adminService.GetAdminByEmail(ctx, email, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error decoding admin data for email")
	mockCollection.AssertExpectations(t)
	mockMongoSingleResult.AssertExpectations(t)
}

func TestGetAdminByEmailFailure_CheckingPassword(t *testing.T){
	ctx := context.Background()
	 
	mockCollection := new(MockMongoCollection)
	mockMongoSingleResult := new(MockMongoSingleResult)
	mockHasher := new(MockHasher)
	adminService := services.NewAdminService(mockCollection, mockHasher)

	email := "testadmin@vigor.com"
	password := "securepassword"
	filter := bson.M{"email": email}

	mockCollection.On("FindOne", ctx, filter).Return(mockMongoSingleResult)
	mockMongoSingleResult.On("Err").Return(nil)
	mockMongoSingleResult.On("Decode", mock.AnythingOfType("*models.Admin")).Return(nil)
	mockHasher.On("CheckPasswordHash", password, mock.AnythingOfType("string")).Return(false)

	_, err := adminService.GetAdminByEmail(ctx, email, password)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidAdminCredentials, err)
	mockCollection.AssertExpectations(t)
	mockMongoSingleResult.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}