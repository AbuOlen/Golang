package documentstore

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	collection := Collection{
		docs: map[string]Document{
			"key1": {
				Fields: map[string]DocumentField{
					"name": {Type: DocumentFieldTypeString, Value: "doc1"},
				},
			},
		},
		config: CollectionConfig{PrimaryKey: "name"},
		index: map[string]CollectionIndex{},
	}

	data, err := json.Marshal(collection)

	assert.NoError(t, err, "unexpected error during marshalling")

	expectedJSON := `{"docs":{"key1":{"Fields":{"name":{"Type":"string","Value":"doc1"}}}},"config":{"PrimaryKey":"name"}, "index":{}}`
	assert.JSONEq(t, expectedJSON, string(data), "marshalled JSON does not match")
}

func TestUnmarshalJSON(t *testing.T) {
	jsonData := `{"docs":{"key1":{"Fields":{"name":{"Type":"string","Value":"doc1"}}}},"config":{"PrimaryKey":"name"}}`

	var collection Collection
	err := json.Unmarshal([]byte(jsonData), &collection)

	assert.NoError(t, err, "unexpected error during unmarshalling")

	expectedCollection := Collection{
		docs: map[string]Document{
			"key1": {
				Fields: map[string]DocumentField{
					"name": {Type: DocumentFieldTypeString, Value: "doc1"},
				},
			},
		},
		config: CollectionConfig{PrimaryKey: "name"},
	}

	assert.Equal(t, expectedCollection, collection, "unmarshalled collection does not match the expected result")
}

func TestPut(t *testing.T) {
	collection := Collection{
		docs:   make(map[string]Document),
		config: CollectionConfig{PrimaryKey: "id"},
	}

	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "key1"},
			"name": {Type: DocumentFieldTypeString, Value: "Example"},
		},
	}

	collection.Put(doc)

	storedDoc, ok := collection.docs["key1"]
	assert.True(t, ok, "document with key 'key1' was not added to the collection")
	assert.Equal(t, doc, storedDoc, "stored document does not match the original")
}

func TestGet(t *testing.T) {
	collection := Collection{
		docs:   make(map[string]Document),
		config: CollectionConfig{PrimaryKey: "id"},
	}

	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "key1"},
			"name": {Type: DocumentFieldTypeString, Value: "Example"},
		},
	}

	collection.Put(doc)

	// Test existing key
	getDoc, ok := collection.Get("key1")
	assert.True(t, ok, "document with key 'key1' should exist")
	assert.NotNil(t, getDoc, "document pointer should not be nil")

	assert.Equal(t, doc, *getDoc, "retrieved document does not match the expected one")

	// Test non-existing key
	_, ok = collection.Get("key2")
	assert.False(t, ok, "document with key 'key2' should not exist")
}

func TestDelete(t *testing.T) {
	collection := Collection{
		docs: map[string]Document{
			"key1": {
				Fields: map[string]DocumentField{
					"name": {Type: DocumentFieldTypeString, Value: "doc1"},
				},
			},
		},
	}

	// Test deletion of an existing key
	success := collection.Delete("key1")
	assert.True(t, success, "expected key 'key1' to be deleted successfully")

	_, exists := collection.docs["key1"]
	assert.False(t, exists, "key 'key1' should no longer exist in the collection")

	// Test deletion of a non-existing key
	success = collection.Delete("key2")
	assert.False(t, success, "deleting a non-existing key should return false")
}

func TestList(t *testing.T) {
	collection := Collection{
		docs: map[string]Document{
			"key1": {
				Fields: map[string]DocumentField{
					"name": {Type: DocumentFieldTypeString, Value: "doc1"},
				},
			},
			"key2": {
				Fields: map[string]DocumentField{
					"name": {Type: DocumentFieldTypeString, Value: "doc2"},
				},
			},
		},
	}

	docs := collection.List()
	expectedDocs := []Document{
		{
			Fields: map[string]DocumentField{
				"name": {Type: DocumentFieldTypeString, Value: "doc1"},
			},
		},
		{
			Fields: map[string]DocumentField{
				"name": {Type: DocumentFieldTypeString, Value: "doc2"},
			},
		},
	}

	assert.Equal(t, len(expectedDocs), len(docs), "number of documents does not match the expected count")

	// Check if all expected documents are present in the result
	for _, doc := range expectedDocs {
		assert.Contains(t, docs, doc, "expected document %+v not found in the result", doc)
	}
}