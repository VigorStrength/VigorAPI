package db

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBService interface {
	ConnectDB(ctx context.Context, cfg *config.Config) error
	EnsureIndexes(ctx context.Context, dbName string) error
	CreateIndexes(ctx context.Context, dbName, collectionName string, indexes []mongo.IndexModel) error
	InitializeCollections(ctx context.Context, dbName string) error
	ApplyCollectionValidation(ctx context.Context, db *mongo.Database, collectionName string, schemaBson bson.M) error
	DisconnectDB(ctx context.Context) 
}
