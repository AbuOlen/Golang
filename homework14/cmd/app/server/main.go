package server

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run() {
	ctx := context.Background()
	opts := options.Client()

	opts.ApplyURI("mongodb://root:root@localhost:27017")

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("Error connecting to MongoDB: %v", err))
	}

	handler := NewDocumentHandler(c)

	http.HandleFunc("/put_document", handler.handlePut)
	http.HandleFunc("/get_document", handler.handleGet)
	http.HandleFunc("/delete_document", handler.handleDelete)
	http.HandleFunc("/list_document", handler.handleList)
	http.HandleFunc("/create_collection", handler.handleCreateCollection)
	http.HandleFunc("/list_collections", handler.handleListCollections)
	http.HandleFunc("/delete_collection", handler.handleDeleteCollection)
	http.HandleFunc("/create_index", handler.handleCreateIndex)
	http.HandleFunc("/delete_index", handler.handleDeleteIndex)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		panic(fmt.Errorf("server listening failed: %v", err))
	}
}