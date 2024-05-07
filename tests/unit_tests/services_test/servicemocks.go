package s

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoDatabase struct {
	mock.Mock
	db.MongoDatabase
}

func (md *MockMongoDatabase) Client() db.MongoClient {
	args := md.Called()
	return args.Get(0).(db.MongoClient)
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

func (mc *MockMongoCollection) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	args := mc.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (mc *MockMongoCollection) Indexes() db.MongoIndexView {
	args := mc.Called()
	return args.Get(0).(db.MongoIndexView)
}

func (mc *MockMongoCollection) Find(ctx context.Context, filter interface{}) (db.MongoCursor, error) {
	args := mc.Called(ctx, filter)
	return args.Get(0).(db.MongoCursor), args.Error(1)
}

func (mc *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) db.MongoSingleResult {
	args := mc.Called(ctx, filter, opts)
	return args.Get(0).(db.MongoSingleResult)
}

func (mc *MockMongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) db.MongoSingleResult {
	args := mc.Called(ctx, filter, update, opts)
	return args.Get(0).(db.MongoSingleResult)
}

func (mc *MockMongoCollection) InsertMany(ctx context.Context, documents []interface{}) (db.MongoInsertManyResult, error) {
	args := mc.Called(ctx, documents)
	return args.Get(0).(db.MongoInsertManyResult), args.Error(1)
}

func (mc *MockMongoCollection) InsertOne(ctx context.Context, document interface{}) (db.MongoInsertOneResult, error) {
	args := mc.Called(ctx, document)
	return args.Get(0).(db.MongoInsertOneResult), args.Error(1)
}

func (mc *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (db.MongoUpdateResult, error) {
	args := mc.Called(ctx, filter, update)
	return args.Get(0).(db.MongoUpdateResult), args.Error(1)
}

func (mc *MockMongoCollection) DeleteOne(ctx context.Context, filter interface{}) (db.MongoDeleteResult, error) {
	args := mc.Called(ctx, filter)
	return args.Get(0).(db.MongoDeleteResult), args.Error(1)
}

type MockMongoSingleResult struct {
	mock.Mock
	db.MongoSingleResult
}

type MockHasher struct {
	mock.Mock
	utils.DefaultHasher
}

type MockParser struct {
	mock.Mock
	utils.DefaultParser
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

type MockMongoInsertOneResult struct {
	mock.Mock
	db.MongoInsertOneResult
}
