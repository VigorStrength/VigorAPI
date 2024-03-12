package mongodb_test

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MockMongoClient struct {
	mock.Mock
	db.MongoClient
}

func (mc *MockMongoClient) Connect(ctx context.Context, opts ...*options.ClientOptions) (db.MongoClient, error) {
	args := mc.Called(ctx, opts)
	return args.Get(0).(db.MongoClient), args.Error(1)
}

func (mc *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	args := mc.Called(ctx, rp)
	return args.Error(0)
}

func (mc *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) db.MongoDatabase {
	args := mc.Called(name, opts)
	return args.Get(0).(db.MongoDatabase)
}

func (mc *MockMongoClient) Disconnect(ctx context.Context) error {
	args := mc.Called(ctx)
	return args.Error(0)
}

type MockMongoDatabase struct {
	mock.Mock
	db.MongoDatabase
}

func (md *MockMongoDatabase) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	args := md.Called(ctx, name, opts)
	return args.Error(0)
}

func (md *MockMongoDatabase) Collection(name string) db.MongoCollection {
	args := md.Called(name)
	return args.Get(0).(db.MongoCollection)
}

func (md *MockMongoDatabase) RunCommand(ctx context.Context, runCmd interface{}) db.MongoSingleResult {
	args := md.Called(ctx, runCmd)
	return args.Get(0).(db.MongoSingleResult)
}

type MockMongoCollection struct {
	mock.Mock
	db.MongoCollection
}

func (mdc *MockMongoCollection) Indexes() db.MongoIndexView {
	args := mdc.Called()
	return args.Get(0).(db.MongoIndexView)
}

type MockMongoSingleResult struct {
	mock.Mock
	db.MongoSingleResult
}

func (msr *MockMongoSingleResult) Err() error {
	args := msr.Called()
	return args.Error(0)
}

type MockMongoIndexView struct {
	mock.Mock
	db.MongoIndexView
}

func (mi *MockMongoIndexView) CreateOne(ctx context.Context, model mongo.IndexModel) (string, error) {
	args := mi.Called(ctx, model)
	return args.String(0), args.Error(1)
}