package s

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/stretchr/testify/mock"
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
