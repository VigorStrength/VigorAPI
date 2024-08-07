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
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
}

// type MongoSession interface {
// 	StartTransaction(opts ...*options.TransactionOptions) error
// 	AbortTransaction(ctx context.Context) error
// 	CommitTransaction(ctx context.Context) error
// 	EndSession(ctx context.Context)
// 	WithTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error),
// 		opts ...*options.TransactionOptions) (interface{}, error)
// 	ClusterTime() bson.Raw
// 	OperationTime() *primitive.Timestamp
// 	Client() *mongo.Client
// 	ID() bson.Raw
// 	AdvanceClusterTime(clusterTime bson.Raw) error
// 	AdvanceOperationTime(operationTime *primitive.Timestamp) error
// }

type MongoDatabase interface {
	Client() MongoClient
	CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error
	Collection(name string) MongoCollection
	RunCommand(ctx context.Context, runCommand interface{}) MongoSingleResult
}

type MongoCollection interface {
	CountDocuments(ctx context.Context, filter interface{}) (int64, error)
	Indexes() MongoIndexView
	Find(ctx context.Context, filter interface{}) (MongoCursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) MongoSingleResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) MongoSingleResult
	InsertMany(ctx context.Context, documents []interface{}) (MongoInsertManyResult, error)
	InsertOne(ctx context.Context, document interface{}) (MongoInsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (MongoUpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (MongoDeleteResult, error)
}

type MongoSingleResult interface {
	Decode(v interface{}) error
	Err() error
}

type MongoIndexView interface {
	CreateOne(ctx context.Context, model mongo.IndexModel) (string, error)
}

type MongoCursor interface {
	All(ctx context.Context, results interface{}) error
	Next(ctx context.Context) bool
	Decode(v interface{}) error
	Close(ctx context.Context) error
	Err() error
}

type DBService interface {
	ConnectDB(ctx context.Context, cfg *config.Config) error
	EnsureIndexes(ctx context.Context, dbName string) error
	CreateIndexes(ctx context.Context, dbName, collectionName string, indexes []mongo.IndexModel) error
	InitializeCollections(ctx context.Context, db MongoDatabase) error
	ApplyCollectionValidation(ctx context.Context, db MongoDatabase, collectionName string, schemaBson bson.M) error
	DisconnectDB(ctx context.Context) 
}

type MongoInsertOneResult struct {
	InsertedID interface{}
}

type MongoInsertManyResult struct {
	InsertedIDs []interface{}
}

type MongoUpdateResult struct {
	MatchedCount  int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID    interface{}
}

type MongoDeleteResult struct {
	DeletedCount int64
}

