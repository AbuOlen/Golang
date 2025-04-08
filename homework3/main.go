package main

import (
	"fmt"
	ds "hw3/documentstore"
)

func main() {
	d1 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d1.Fields["key"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key1"}
	d1.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "value"}

	ds.Put(d1)

	d2 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d2.Fields["key"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key2"}
	d2.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeBool, Value: true}

	ds.Put(d2)

	d3 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d3.Fields["key"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key3"}
	d3.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeNumber, Value: 73}

	ds.Put(d3)

	d4 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d4.Fields["key"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key3"}
	d4.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeArray, Value: []interface{}{1, 2, 3}}

	ds.Put(d4)

	d5 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d5.Fields["key"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key3"}
	d5.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeObject, Value: map[string]interface{}{}}

	ds.Put(d5)

	docsGet, _ := ds.Get("key1")
	for k, v := range docsGet.Fields {
		fmt.Printf("get => key: %s value: %s\n", k, v.Value)
	}

	docs := ds.List()
	for _, doc := range docs {
		for k, v := range doc.Fields {
			fmt.Printf("list => key: %s value: %s\n", k, v.Value)
		}
	}

	ds.Delete("key1")
	docsAfterDelete := ds.List()
	for _, doc := range docsAfterDelete {
		for k, v := range doc.Fields {
			fmt.Printf("after del => key: %s value: %s\n", k, v.Value)
		}
	}
}

//get => key: key value: key1
//get => key: val value: value
//list => key: key value: key3
//list => key: val value: map[]
//list => key: key value: key1
//list => key: val value: value
//list => key: key value: key2
//list => key: val value: %!s(bool=true)
//after del => key: key value: key2
//after del => key: val value: %!s(bool=true)
//after del => key: key value: key3
//after del => key: val value: map[]
