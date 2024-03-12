package mongodb_test

import (
	"context"
	"testing"

	"github.com/GhostDrew11/vigor-api/internal/config"
	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConnectDB(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		MongoDBURI:   "mongodb://localhost:27017",
		DatabaseName: "Vigor_test",
	}

	mockClient := new(MockMongoClient)
	mockDB := new(MockMongoDatabase)
	mockSingleResult := new(MockMongoSingleResult)
	mockCollection := new(MockMongoCollection)
	mockIndexView := new(MockMongoIndexView)
	mockIndexName := "mockIndexName"
	service := db.NewMongoDBService(mockClient)
	

	mockClient.On("Connect", ctx, mock.AnythingOfType("[]*options.ClientOptions")).Return(mockClient, nil)
	mockClient.On("Ping", ctx, mock.AnythingOfType("*readpref.ReadPref")).Return(nil)
	mockClient.On("Database", cfg.DatabaseName, mock.AnythingOfType("[]*options.DatabaseOptions")).Return(mockDB)
	
	mockDB.On("CreateCollection", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]*options.CreateCollectionOptions")).Return(nil)
	// mockDB.On("RunCommand", ctx, mock.AnythingOfType("bson.D")).Return(mockSingleResult)
	// mockSingleResult.On("Err").Return(nil)
	mockDB.On("Collection", mock.AnythingOfType("string")).Return(mockCollection)

	

	mockCollection.On("Indexes").Return(mockIndexView)

	mockIndexView.On("CreateOne", ctx, mock.AnythingOfType("mongo.IndexModel")).Return(mockIndexName, nil)

	err := service.ConnectDB(ctx, cfg)
	mockClient.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockSingleResult.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
	mockIndexView.AssertExpectations(t)
	assert.NoError(t, err)
}
