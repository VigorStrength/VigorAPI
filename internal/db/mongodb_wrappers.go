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

func (msr *mongoSingleResultWrapper) Decode(v interface{}) error {
	return msr.singleResult.Decode(v)
}

type mongoCollectionWrapper struct{
	collection *mongo.Collection
}

func (mdc *mongoCollectionWrapper) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	return mdc.collection.CountDocuments(ctx, filter)
}

func (mdc *mongoCollectionWrapper) Indexes() MongoIndexView {
	return &mongoIndexViewWrapper{indexView: mdc.collection.Indexes()}
}

func (mdc *mongoCollectionWrapper) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) MongoSingleResult {
	return &mongoSingleResultWrapper{singleResult: mdc.collection.FindOne(ctx, filter, opts...)}
}

func (mdc *mongoCollectionWrapper) Find(ctx context.Context, filter interface{}) (MongoCursor, error) {
	cursor, err := mdc.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &mongoCursorWrapper{cursor: cursor}, nil
}

func (mdc *mongoCollectionWrapper) InsertMany(ctx context.Context, documents []interface{}) (MongoInsertManyResult, error) {
	result, err := mdc.collection.InsertMany(ctx, documents)
	if err != nil {
		return MongoInsertManyResult{}, err
	}
	return MongoInsertManyResult{InsertedIDs: result.InsertedIDs}, nil
}

func (mdc *mongoCollectionWrapper) InsertOne(ctx context.Context, document interface{}) (MongoInsertOneResult, error) {
	result, err := mdc.collection.InsertOne(ctx, document)
	if err != nil {
		return MongoInsertOneResult{}, err
	}
	return MongoInsertOneResult{InsertedID: result.InsertedID}, nil
}

func (mdc *mongoCollectionWrapper) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (MongoUpdateResult, error) {
	result, err := mdc.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return MongoUpdateResult{}, err
	}
	return MongoUpdateResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}, nil
}

func (mdc *mongoCollectionWrapper) DeleteOne(ctx context.Context, filter interface{}) (MongoDeleteResult, error) {
	result, err := mdc.collection.DeleteOne(ctx, filter)
	if err != nil {
		return MongoDeleteResult{}, err
	}
	return MongoDeleteResult{DeletedCount: result.DeletedCount}, nil
}

type mongoIndexViewWrapper struct {
	indexView mongo.IndexView
}

func (miv *mongoIndexViewWrapper) CreateOne(ctx context.Context, model mongo.IndexModel) (string, error) {
	return miv.indexView.CreateOne(ctx, model)
}

type mongoCursorWrapper struct {
	cursor *mongo.Cursor
}

func (mcw *mongoCursorWrapper) All(ctx context.Context, results interface{}) error {
	return mcw.cursor.All(ctx, results)
}

func (mcw *mongoCursorWrapper) Next(ctx context.Context) bool {
	return mcw.cursor.Next(ctx)
}

func (mcw *mongoCursorWrapper) Decode(v interface{}) error {
	return mcw.cursor.Decode(v)
}

func (mcw *mongoCursorWrapper) Close(ctx context.Context) error {
	return mcw.cursor.Close(ctx)
}

func (mcw *mongoCursorWrapper) Err() error {
	return mcw.cursor.Err()
}
