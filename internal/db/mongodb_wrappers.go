package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoClientWrapper struct {
	client *mongo.Client
}

func NewMongoClientWrapper(client *mongo.Client) MongoClient {
	return &mongoClientWrapper{client: client}
}

func (m *mongoClientWrapper) Connect(ctx context.Context, opts ...*options.ClientOptions) (MongoClient, error) {
	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		return nil, err
	}
	m.client = client
	return m, nil
}

func (m *mongoClientWrapper) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return m.client.Ping(ctx, rp)
}

func (m *mongoClientWrapper) Database(name string, opts ...*options.DatabaseOptions) MongoDatabase {
	return &mongoDatabaseWrapper{database: m.client.Database(name, opts...)}
}

func (m *mongoClientWrapper) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

type mongoDatabaseWrapper struct {
	database *mongo.Database
}

func (md *mongoDatabaseWrapper) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	return md.database.CreateCollection(ctx, name, opts...)
}

func (md *mongoDatabaseWrapper) Collection(name string) MongoCollection {
	return &mongoCollectionWrapper{collection: md.database.Collection(name)}
}

func (md *mongoDatabaseWrapper) RunCommand(ctx context.Context, runCommand interface{}) MongoSingleResult {
	return &mongoSingleResultWrapper{singleResult: md.database.RunCommand(ctx, runCommand)}
}

type mongoSingleResultWrapper struct {
	singleResult *mongo.SingleResult
}

func (msr *mongoSingleResultWrapper) Err() error {
	return msr.singleResult.Err()
}

type mongoCollectionWrapper struct{
	collection *mongo.Collection
}

func (mdc *mongoCollectionWrapper) Indexes() MongoIndexView {
	return &mongoIndexViewWrapper{indexView: mdc.collection.Indexes()}
}

type mongoIndexViewWrapper struct {
	indexView mongo.IndexView
}

func (miv *mongoIndexViewWrapper) CreateOne(ctx context.Context, model mongo.IndexModel) (string, error) {
	return miv.indexView.CreateOne(ctx, model)
}
