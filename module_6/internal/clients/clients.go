package clients

import (
	"context"
	app_mongo "module6/internal/clients/mongo"
	"module6/internal/config"
)

type Clients struct {
	Mongo *app_mongo.Client
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	mongo, err := app_mongo.NewMongo(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Clients{
		Mongo: mongo,
	}, nil
}
