package mocks

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCursorAdapter adapts mongo.Cursor to our Cursor interface
type MongoCursorAdapter struct {
	*mongo.Cursor
}

func (m *MongoCursorAdapter) All(ctx context.Context, results interface{}) error {
	return m.Cursor.All(ctx, results)
}

// MongoCollectionAdapter adapts mongo.Collection to our Collection interface
type MongoCollectionAdapter struct {
	*mongo.Collection
}

func NewMongoCollectionAdapter(coll *mongo.Collection) Collection {
	return &MongoCollectionAdapter{Collection: coll}
}

func (m *MongoCollectionAdapter) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	cursor, err := m.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &MongoCursorAdapter{Cursor: cursor}, nil
}
