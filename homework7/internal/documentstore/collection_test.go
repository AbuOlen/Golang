package documentstore

import (
	"testing"
	"reflect"
)

func TestPut(t *testing.T) {
	cfg := CollectionConfig{PrimaryKey: "key1"}
	col := Collection{Docs: make(map[string]Document), Config: cfg}
	d1 := Document{Fields: make(map[string]DocumentField)}
	d1.Fields["key1"] = DocumentField{Type: DocumentFieldTypeString, Value: "key1"}
	d1.Fields["val"] = DocumentField{Type: DocumentFieldTypeString, Value: "value"}
	col.Put(d1)
	if len(col.Docs) == 0 {
		t.Errorf("Expected len 1, got %d", len(col.Docs))
	}

}

func TestGet(t *testing.T) {
	cfg := CollectionConfig{PrimaryKey: "key1"}
	col := Collection{Docs: make(map[string]Document), Config: cfg}
	d1 := Document{Fields: make(map[string]DocumentField)}
	d1.Fields["key1"] = DocumentField{Type: DocumentFieldTypeString, Value: "key1"}
	d1.Fields["val"] = DocumentField{Type: DocumentFieldTypeString, Value: "value"}
	col.Put(d1)
	d2, ok := col.Get("key1")

	if !ok {
		t.Error("Expected Get to return true, got false")
	}

	if !reflect.DeepEqual(d1, *d2) {
		t.Error("Expected d1 to be equal to d2")
	}
}

func TestDelete(t *testing.T) {
	cfg := CollectionConfig{PrimaryKey: "key1"}
	col := Collection{Docs: make(map[string]Document), Config: cfg}
	d1 := Document{Fields: make(map[string]DocumentField)}
	d1.Fields["key1"] = DocumentField{Type: DocumentFieldTypeString, Value: "key1"}
	d1.Fields["val"] = DocumentField{Type: DocumentFieldTypeString, Value: "value"}
	col.Put(d1)
	_, ok := col.Get("key1")

	if !ok {
		t.Error("Expected Get to return true, got false")
	}

	if !col.Delete("key1") {
		t.Error("Deletion failed")
	}

	_, ok = col.Get("key1")

	if ok {
		t.Error("Expected Get to return false, got true")
	}
}