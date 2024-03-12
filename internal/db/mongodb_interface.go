package db

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient interface {
	Connect(ctx context.Context, opts ...*options.ClientOptions) (MongoClient, error)
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string, opts ...*options.DatabaseOptions) MongoDatabase
	Disconnect(ctx context.Context) error
}

type MongoDatabase interface {
	CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error
	Collection(name string) MongoCollection
	RunCommand(ctx context.Context, runCommand interface{}) MongoSingleResult
}

type MongoCollection interface {
	Indexes() MongoIndexView
}

type MongoSingleResult interface {
	Err() error
}

type MongoIndexView interface {
	CreateOne(ctx context.Context, model mongo.IndexModel) (string, error)
}

type DBService interface {
	ConnectDB(ctx context.Context, cfg *config.Config) error
	EnsureIndexes(ctx context.Context, dbName string) error
	CreateIndexes(ctx context.Context, dbName, collectionName string, indexes []mongo.IndexModel) error
	InitializeCollections(ctx context.Context, db MongoDatabase) error
	ApplyCollectionValidation(ctx context.Context, db MongoDatabase, collectionName string, schemaBson bson.M) error
	DisconnectDB(ctx context.Context) 
}
