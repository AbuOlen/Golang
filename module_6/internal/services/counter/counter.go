package counter

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"module6/internal/clients"
	"module6/internal/constants"
)

type Service struct {
	counterColl *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		counterColl: clients.Mongo.Db.Collection(constants.CollectionCounter),
	}
}

func (s *Service) IncrementOne(ctx context.Context) (int, error) {
	filter := bson.M{}
	update := bson.M{
		"$inc": bson.M{"count": 1},
	}

	opts := options.FindOneAndUpdate()
	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	res := s.counterColl.FindOneAndUpdate(ctx, filter, update, opts)
	if err := res.Err(); err != nil {
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}
	var counter Document
	err := res.Decode(&counter)
	if err != nil {
		return 0, fmt.Errorf("failed to decode counter result: %w", err)
	}
	return counter.Count, nil
}
