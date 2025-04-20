package documentstore

import (
	"testing"
	"os"
)

func TestNewStore(t *testing.T) {
	store := NewStore()

	if store.Collections == nil {
		t.Error("Expected Collections map to be initialized, but got nil")
	}

	if len(store.Collections) != 0 {
		t.Errorf("Expected Collections to be empty, but got %d", len(store.Collections))
	}

	if store.GetLogger() == nil {
		t.Error("Expected default logger to be set, but got nil")
	}
}

func TestStoreCreateCollection(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{}

	created, col := store.CreateCollection("test", cfg)
	if !created || col == nil {
		t.Errorf("Expected collection to be created, but creation failed")
	}

	if len(store.Collections) != 1 {
		t.Errorf("Expected 1 collection in store, but got %d", len(store.Collections))
	}

	created, _ = store.CreateCollection("test", cfg)
	if created {
		t.Error("Expected duplicate collection creation to return false, but got true")
	}

	created, col = store.CreateCollection("nil_test", nil)
	if created || col != nil {
		t.Error("Expected collection creation with nil config to fail, but it succeeded")
	}
}

func TestStoreGetCollection(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{}
	store.CreateCollection("test", cfg)

	col, exists := store.GetCollection("test")
	if !exists || col == nil {
		t.Error("Expected to retrieve existing collection, but got nil or does not exist")
	}

	col, exists = store.GetCollection("nonexistent")
	if exists || col != nil {
		t.Error("Expected to not retrieve non-existent collection, but found one")
	}
}

func TestStoreDeleteCollection(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{}
	store.CreateCollection("test", cfg)

	deleted := store.DeleteCollection("test")
	if !deleted {
		t.Error("Expected collection to be deleted, but deletion failed")
	}

	if len(store.Collections) != 0 {
		t.Errorf("Expected 0 collections after deletion, but found %d", len(store.Collections))
	}

	deleted = store.DeleteCollection("nonexistent")
	if deleted {
		t.Error("Expected deletion of non-existent collection to return false, but it returned true")
	}
}

func TestStoreDumpToFileAndNewStoreFromFile(t *testing.T) {
	filename := "test_store.json"
	defer os.Remove(filename)

	store := NewStore()
	cfg := &CollectionConfig{}
	store.CreateCollection("test", cfg)

	err := store.DumpToFile(filename)
	if err != nil {
		t.Errorf("Expected DumpToFile to succeed, but got error: %v", err)
	}

	newStore, err := NewStoreFromFile(filename)
	if err != nil {
		t.Errorf("Expected NewStoreFromFile to succeed, but got error: %v", err)
	}

	if len(newStore.Collections) != 1 {
		t.Errorf("Expected 1 collection in restored store, but found %d", len(newStore.Collections))
	}
}