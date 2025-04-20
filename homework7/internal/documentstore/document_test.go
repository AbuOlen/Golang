package documentstore

import "testing"

func TestDocumentFieldTypeConstants(t *testing.T) {
	if DocumentFieldTypeString != "string" {
		t.Errorf("Expected DocumentFieldTypeString to equal 'string', got '%s'", DocumentFieldTypeString)
	}
	if DocumentFieldTypeNumber != "number" {
		t.Errorf("Expected DocumentFieldTypeNumber to equal 'number', got '%s'", DocumentFieldTypeNumber)
	}
	if DocumentFieldTypeBool != "bool" {
		t.Errorf("Expected DocumentFieldTypeBool to equal 'bool', got '%s'", DocumentFieldTypeBool)
	}
	if DocumentFieldTypeArray != "array" {
		t.Errorf("Expected DocumentFieldTypeArray to equal 'array', got '%s'", DocumentFieldTypeArray)
	}
	if DocumentFieldTypeObject != "object" {
		t.Errorf("Expected DocumentFieldTypeObject to equal 'object', got '%s'", DocumentFieldTypeObject)
	}
}

func TestDocumentInitialization(t *testing.T) {
	doc := Document{
		Fields: make(map[string]DocumentField),
	}

	if len(doc.Fields) != 0 {
		t.Errorf("Expected Fields to be empty upon initialization, got len=%d", len(doc.Fields))
	}
}

func TestDocumentAddAndRetrieveField(t *testing.T) {
	doc := Document{
		Fields: make(map[string]DocumentField),
	}

	doc.Fields["field1"] = DocumentField{
		Type:  DocumentFieldTypeNumber,
		Value: 73,
	}

	field, exists := doc.Fields["field1"]
	if !exists {
		t.Error("Expected field 'field1' to exist in document")
	}

	if field.Type != DocumentFieldTypeNumber {
		t.Errorf("Expected Type to equal DocumentFieldTypeNumber, got '%s'", field.Type)
	}

	if field.Value != 73 {
		t.Errorf("Expected Value to equal 42, got '%v'", field.Value)
	}
}

func TestDocumentMissingField(t *testing.T) {
	doc := Document{
		Fields: make(map[string]DocumentField),
	}

	_, exists := doc.Fields["nonexistent"]
	if exists {
		t.Error("Did not expect 'nonexistent' field to exist")
	}
}