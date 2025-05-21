package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx := context.Background()
	opts := options.Client()
	opts.ApplyURI("mongodb://root:root@localhost:27017")
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("mongodb connect failed: %v", err))
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			panic(fmt.Errorf("mongodb disconnect failed: %v", err))
		}
	}(client, ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(fmt.Errorf("mongodb ping failed: %v", err))
	}
	db := client.Database("app")
	if db == nil {
		panic(fmt.Errorf("mongodb database failed: %v", err))
	}
	collection := db.Collection("kv")
	result := collection.FindOne(ctx, bson.M{"user_id": "123"})
	if result == nil {
		panic(fmt.Errorf("mongodb find failed: %v", err))
	}
	var res map[string]interface{}
	err = result.Decode(&res)
	if err != nil {
		panic(fmt.Errorf("mongodb decode failed: %v", err))
	}
	fmt.Println(res)
}
