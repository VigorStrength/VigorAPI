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

func (m *mongoClientWrapper) Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	return mongo.Connect(ctx, opts...)
}

func (m *mongoClientWrapper) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return m.client.Ping(ctx, rp)
}

func (m *mongoClientWrapper) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return m.client.Database(name, opts...)
}

func (m *mongoClientWrapper) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}