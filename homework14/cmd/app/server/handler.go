package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	coll *mongo.Collection
	db   *mongo.Database
}

type Document struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

func NewDocumentHandler(c *mongo.Client) *Handler {
	db := c.Database("app")
	coll := db.Collection("kv")
	return &Handler{coll: coll, db: db}
}

type PutReqBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PutRespBody struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handlePut(w http.ResponseWriter, r *http.Request) {
	reqBody := PutReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	doc := &Document{Key: reqBody.Key, Value: reqBody.Value}
	filter := bson.M{"key": doc.Key}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err = h.coll.UpdateOne(r.Context(), filter, update, opts)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to update document: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := PutRespBody{Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type GetReqBody struct {
	Key string `json:"key"`
}

type GetRespBody struct {
	Value string `json:"value"`
	Ok    bool   `json:"ok"`
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	reqBody := GetReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	filter := bson.M{"key": reqBody.Key}
	doc := &Document{}
	err = h.coll.FindOne(r.Context(), filter).Decode(doc)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to find document: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := GetRespBody{Value: doc.Value, Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type DeleteReqBody struct {
	Key string `json:"key"`
}

type DeleteRespBody struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	reqBody := DeleteReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	var res map[string]interface{}
	filter := bson.M{"key": reqBody.Key}
	err = h.coll.FindOneAndDelete(r.Context(), filter).Decode(&res)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to delete document: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	//respBody := DeleteRespBody{Ok: res.DeletedCount > 0}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type ListReqBody struct{}

type ListRespBody struct {
	Docs []Document `json:"items"`
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	cur, err := h.coll.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, fmt.Errorf("failed to find documents: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	docs := []Document{}
	for cur.Next(r.Context()) {
		var doc Document
		err = cur.Decode(&doc)
		if err != nil {
			http.Error(w, fmt.Errorf("failed to decode document: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		docs = append(docs, doc)
	}
	//err = cur.All(r.Context(), &docs)
	//if err != nil {
	//	http.Error(w, fmt.Errorf("failed to decode documents: %w", err).Error(), http.StatusInternalServerError)
	//	return
	//}

	respBody := ListRespBody{Docs: docs}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type CreateCollReqBody struct {
	Name string `json:"name"`
}

type CreateCollRespBody struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handleCreateCollection(w http.ResponseWriter, r *http.Request) {
	reqBody := CreateCollReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	err = h.db.CreateCollection(r.Context(), reqBody.Name)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to create collection: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := CreateCollRespBody{Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type DeleteCollReqBody struct {
	Name string `json:"name"`
}

type DeleteCollRespBody struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handleDeleteCollection(w http.ResponseWriter, r *http.Request) {
	reqBody := DeleteCollReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	coll := h.db.Collection(reqBody.Name)
	err = coll.Drop(r.Context())
	if err != nil {
		http.Error(w, fmt.Errorf("failed to delete collection: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := DeleteCollRespBody{Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type ListCollReqBody struct {
}

type ListCollRespBody struct {
	Collections []string `json:"items"` // Changed to string since we only need collection names
}

func (h *Handler) handleListCollections(w http.ResponseWriter, r *http.Request) {
	// ListCollections needs a filter (can be empty)
	cursor, err := h.db.ListCollections(r.Context(), bson.D{})
	if err != nil {
		http.Error(w, fmt.Errorf("failed to list collections: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	// Get all collection names
	var collections []string
	for cursor.Next(r.Context()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			http.Error(w, fmt.Errorf("failed to decode collection info: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		if name, ok := result["name"].(string); ok {
			collections = append(collections, name)
		}
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, fmt.Errorf("cursor error: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := ListCollRespBody{Collections: collections}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type CreateIndexReqBody struct {
	CollectionName string `json:"collection_name"`
	Field          string `json:"field"`      // Field to index
	IndexType      string `json:"index_type"` // "asc" or "desc"
	Unique         bool   `json:"unique"`     // Whether the index should be unique
}

type CreateIndexRespBody struct {
	IndexName string `json:"index_name"`
	Ok        bool   `json:"ok"`
}

func (h *Handler) handleCreateIndex(w http.ResponseWriter, r *http.Request) {
	reqBody := CreateIndexReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	// Validate inputs
	if reqBody.CollectionName == "" || reqBody.Field == "" {
		http.Error(w, "collection name and field are required", http.StatusBadRequest)
		return
	}

	// Get the collection
	coll := h.db.Collection(reqBody.CollectionName)

	// Create index model
	indexValue := 1 // ascending by default
	if reqBody.IndexType == "desc" {
		indexValue = -1
	}

	model := mongo.IndexModel{
		Keys:    bson.D{{Key: reqBody.Field, Value: indexValue}},
		Options: options.Index().SetUnique(reqBody.Unique),
	}

	// Create the index
	indexName, err := coll.Indexes().CreateOne(r.Context(), model)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to create index: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := CreateIndexRespBody{
		IndexName: indexName,
		Ok:        true,
	}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type DeleteIndexReqBody struct {
	CollectionName string `json:"collection_name"`
	IndexName      string `json:"index_name"`
}

type DeleteIndexRespBody struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handleDeleteIndex(w http.ResponseWriter, r *http.Request) {
	reqBody := DeleteIndexReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	// Validate inputs
	if reqBody.CollectionName == "" || reqBody.IndexName == "" {
		http.Error(w, "collection name and index name are required", http.StatusBadRequest)
		return
	}

	// Get the collection
	coll := h.db.Collection(reqBody.CollectionName)

	// Drop the index
	_, err = coll.Indexes().DropOne(r.Context(), reqBody.IndexName)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to delete index: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := DeleteIndexRespBody{Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}
