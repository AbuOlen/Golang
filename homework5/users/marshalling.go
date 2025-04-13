package users

import (
	"errors"
	"fmt"
	ds "hw5/documentstore"
	"reflect"
)

var ErrCollectionAlreadyExists = errors.New("collection already exists")
var ErrUnsupportedDocumentField = errors.New("unsupported document field")
var ErrDocumentIsNil = errors.New("document pointer is nil")
var ErrWrongDataType = errors.New("wrong data type")

func MarshalDocument(input any) (*ds.Document, error) {
	var v = ds.Document{Fields: make(map[string]ds.DocumentField)}
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		var fieldType ds.DocumentFieldType
		switch field.Type.Name() {
		case "string":
			fieldType = ds.DocumentFieldTypeString
		case "array":
			fieldType = ds.DocumentFieldTypeArray
		case "bool":
			fieldType = ds.DocumentFieldTypeBool
		case "number":
			fieldType = ds.DocumentFieldTypeNumber
		case "object":
			fieldType = ds.DocumentFieldTypeObject
		default:
			return nil, ErrUnsupportedDocumentField
		}
		v.Fields[field.Name] = ds.DocumentField{Value: val.Field(i).Interface(), Type: fieldType}
	}
	return &v, nil

}

func UnmarshalDocument(doc *ds.Document, output any) error {
	if doc == nil {
		return fmt.Errorf("document is nil: %w", ErrDocumentIsNil)
	}
	val := reflect.ValueOf(output)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("output must be a pointer to a struct")
	}
	structValue := val.Elem()
	structType := structValue.Type()

	// Iterate over the keys in the map
	for key, value := range doc.Fields {
		if field, _ := structType.FieldByName(key); field.IsExported() {
			fieldValue := structValue.FieldByName(key)

			// Check if the fieldValue is settable (can assign to it)
			if fieldValue.CanSet() {
				// Convert the map value to the appropriate type
				mapValue := reflect.ValueOf(value.Value)

				// Ensure type compatibility
				if mapValue.Type().AssignableTo(fieldValue.Type()) {
					fieldValue.Set(mapValue) // Assign the value
				} else if mapValue.Type().ConvertibleTo(fieldValue.Type()) {
					// Convert and set if types are convertible
					fieldValue.Set(mapValue.Convert(fieldValue.Type()))
				} else {
					return fmt.Errorf("cannot assign value of type %s to field %s of type %s",
						mapValue.Type(), key, fieldValue.Type())
				}
			}
		}
	}
	return nil
}
