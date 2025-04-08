package main

import (
	"fmt"
	ds "hw4/documentstore"
)

func main() {
	d1 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d1.Fields["key1"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key1"}
	d1.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "value"}

	d2 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d2.Fields["key2"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key2"}
	d2.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeBool, Value: true}

	d3 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d3.Fields["key3"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key3"}
	d3.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeNumber, Value: 73}

	d4 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d4.Fields["key3"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key3"}
	d4.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeArray, Value: []interface{}{1, 2, 3}}

	d5 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d5.Fields["key3"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key3"}
	d5.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeObject, Value: map[string]interface{}{}}


	cfg1 := ds.CollectionConfig{PrimaryKey: "key1"}
	cfg2 := ds.CollectionConfig{PrimaryKey: "key2"}
	cfg3 := ds.CollectionConfig{PrimaryKey: "key3"}

	store := ds.NewStore()
	ok, col1 := store.CreateCollection("key1", &cfg1)
	if !ok {
		fmt.Println("Collection creation failed")
		return
	}
	col1.Put(d1)

	ok, col2 := store.CreateCollection("key2", &cfg2)
	if !ok {
		fmt.Println("Collection creation failed")
		return
	}
	col2.Put(d2)

	ok, col3 := store.CreateCollection("key3", &cfg3)
	if !ok {
		fmt.Println("Collection creation failed")
		return
	}
	col3.Put(d3)
	col3.Put(d4)
	col3.Put(d5)

	c, ok := store.GetCollection("key1")
	if !ok {
		fmt.Println("Collection with name 'key1' not found")
		return
	}
	docs := c.List()
	for _, doc := range docs {
		for k, v := range doc.Fields {
			fmt.Printf("list => key: %s value: %s\n", k, v.Value)
		}
	}

	deleted := store.DeleteCollection("key1")
	if !deleted {
		fmt.Println("Could not delete collection with name 'key1'")
		return
	}
	c, exists := store.GetCollection("key1")
	if exists {
		fmt.Println("Collection with name 'key1' was not deleted")
		return
	}
	fmt.Println("Done!")
}

//list => key: key1 value: key1
//list => key: val value: value
//Done!

