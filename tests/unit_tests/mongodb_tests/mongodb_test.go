package mongodb_test

import (
	"context"
	"embed"
	"testing"

	"github.com/GhostDrew11/vigor-api/internal/config"
	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MockMongoClient struct {
	mock.Mock
	db.MongoClient
}

func (m *MockMongoClient) Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*mongo.Client), args.Error(1)
}

func (m *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error { 
	args := m.Called(ctx, rp)
	return args.Error(0)
}

func (m *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	args := m.Called(name, opts)
	return args.Get(0).(*mongo.Database)
}

func (m *MockMongoClient) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestConnectDB(t *testing.T) {

	var mockSchemaFiles embed.FS

	ctx := context.Background()
	cfg := &config.Config{
		MongoDBURI:   "mongodb://localhost:27017",
        DatabaseName: "Vigor_test",
	}

	mockClient := new(MockMongoClient)
	service := &db.MongoDBService{
		Client: mockClient,
		SchemaFiles: mockSchemaFiles,
	}
	

	mockClient.On("Connect", ctx, mock.AnythingOfType("[]*options.ClientOptions")).Return(&mongo.Client{}, nil)
	mockClient.On("Ping", ctx, mock.AnythingOfType("*readpref.ReadPref")).Return(nil)

	err := service.ConnectDB(ctx, cfg)
	mockClient.AssertExpectations(t)
	assert.NoError(t, err)
}