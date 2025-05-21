package documentstore

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	assert.NotNil(t, store, "NewStore should create a non-nil Store instance")
	assert.NotNil(t, store.collections, "NewStore should initialize the collections map")
}

func TestCreateCollection(t *testing.T) {
	store := NewStore()

	cfg := &CollectionConfig{PrimaryKey: "id"}
	created, col := store.CreateCollection("test_collection", cfg)

	assert.True(t, created, "the collection should be created successfully")
	assert.NotNil(t, col, "the returned collection pointer should not be nil")
	assert.Contains(t, store.collections, "test_collection", "the collection should exist in the store")

	// Attempt to create a collection with the same name
	created, col = store.CreateCollection("test_collection", cfg)
	assert.False(t, created, "creating a collection with the same name should fail")
	assert.Nil(t, col, "the returned collection should be nil when creation fails")
}

func TestGetCollection(t *testing.T) {
	store := NewStore()

	cfg := &CollectionConfig{PrimaryKey: "id"}
	store.CreateCollection("test_collection", cfg)

	col, found := store.GetCollection("test_collection")
	assert.True(t, found, "the collection should exist in the store")
	assert.NotNil(t, col, "the returned collection pointer should not be nil")

	// Test retrieval of a non-existent collection
	col, found = store.GetCollection("non_existent")
	assert.False(t, found, "the collection should not exist in the store")
	assert.Nil(t, col, "the returned value for a non-existent collection should be nil")
}

func TestDeleteCollection(t *testing.T) {
	store := NewStore()

	cfg := &CollectionConfig{PrimaryKey: "id"}
	store.CreateCollection("test_collection", cfg)

	deleted := store.DeleteCollection("test_collection")
	assert.True(t, deleted, "the collection should be deleted successfully")
	assert.NotContains(t, store.collections, "test_collection", "the collection should no longer exist in the store")

	// Test attempting to delete a non-existent collection
	deleted = store.DeleteCollection("non_existent")
	assert.False(t, deleted, "deleting a non-existent collection should return false")
}

func TestDumpAndNewStoreFromDump(t *testing.T) {
	store := NewStore()

	// Add some collections
	cfg := &CollectionConfig{PrimaryKey: "id"}
	store.CreateCollection("collection1", cfg)
	store.CreateCollection("collection2", cfg)

	// Dump the store to JSON
	dump, err := store.Dump()
	assert.NoError(t, err, "Dump should not return an error")
	assert.NotNil(t, dump, "Dump should return non-nil data")

	// Create a new store from the dump
	newStore, err := NewStoreFromDump(dump)
	assert.NoError(t, err, "NewStoreFromDump should not return an error")
	assert.NotNil(t, newStore, "NewStoreFromDump should return a valid Store instance")
	assert.Len(t, newStore.collections, len(store.collections), "The new store should have the same number of collections")

	// Verify collections exist in the new store
	for name := range store.collections {
		_, exists := newStore.collections[name]
		assert.True(t, exists, "Collection '%s' should exist in the new store", name)
	}
}

func TestDumpToFileAndNewStoreFromFile(t *testing.T) {
	store := NewStore()

	// Add a collection
	cfg := &CollectionConfig{PrimaryKey: "id"}
	store.CreateCollection("test_collection", cfg)

	// Dump the store to a file
	fileName := "test_store_dump.json"
	defer os.Remove(fileName)

	err := store.DumpToFile(fileName)
	assert.NoError(t, err, "DumpToFile should not return an error")

	// Load the store back from the file
	newStore, err := NewStoreFromFile(fileName)
	assert.NoError(t, err, "NewStoreFromFile should not return an error")
	assert.NotNil(t, newStore, "NewStoreFromFile should return a valid Store instance")
	assert.Len(t, newStore.collections, 1, "The new store should have exactly one collection")

	// Verify the collection exists in the new store
	_, exists := newStore.collections["test_collection"]
	assert.True(t, exists, "The collection 'test_collection' should exist in the new store")
}

func TestMarshalJSONAndUnmarshalJSON(t *testing.T) {
	store := NewStore()

	// Add a collection
	cfg := &CollectionConfig{PrimaryKey: "id"}
	store.CreateCollection("test_collection", cfg)

	// Marshal the store
	data, err := json.Marshal(store)
	assert.NoError(t, err, "MarshalJSON should not return an error")

	// Unmarshal the JSON into a new store
	var newStore Store
	err = json.Unmarshal(data, &newStore)
	assert.NoError(t, err, "UnmarshalJSON should not return an error")
	assert.Len(t, newStore.collections, 1, "The unmarshalled store should have exactly one collection")

	// Verify the collections are the same in both stores
	_, exists := newStore.collections["test_collection"]
	assert.True(t, exists, "The collection 'test_collection' should exist in the unmarshalled store")
}

func strPtr(s string) *string {
	return &s
}

func BenchmarkCollectionCreateIndex(b *testing.B) {
	cfg1 := CollectionConfig{PrimaryKey: "key1"}
	store := NewStore()
	_, col1 := store.CreateCollection("key1", &cfg1)
	for i := 0; i < b.N; i++ {
		d := Document{Fields: make(map[string]DocumentField)}
		d.Fields["key1"] = DocumentField{Type: DocumentFieldTypeString, Value: "key"+strconv.Itoa(i)}
		d.Fields["val"] = DocumentField{Type: DocumentFieldTypeString, Value: "val"+strconv.Itoa(i)}
		col1.Put(d)
	}

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		col1.CreateIndex("val")
		col1.DeleteIndex("val")
	}
}

func BenchmarkCollectionQuery(b *testing.B) {
	cfg1 := CollectionConfig{PrimaryKey: "key1"}
	store := NewStore()
	_, col1 := store.CreateCollection("key1", &cfg1)
	for i := 0; i < b.N; i++ {
		d := Document{Fields: make(map[string]DocumentField)}
		d.Fields["key1"] = DocumentField{Type: DocumentFieldTypeString, Value: "key"+strconv.Itoa(i)}
		d.Fields["val"] = DocumentField{Type: DocumentFieldTypeString, Value: "val"+strconv.Itoa(i)}
		col1.Put(d)
	}
	col1.CreateIndex("val")
	b.ResetTimer()
	for i := 0; i < 100; i++ {
		min := rand.Intn(b.N) + 1
		max := rand.Intn(b.N) + 1
		if(max < min) {
			min, max = max, min
		}
		col1.Query("val", QueryParams{Desc: max % 2 == 0, MinValue: strPtr("val"+strconv.Itoa(min)), MaxValue: strPtr("val"+strconv.Itoa(max))})
	}
}