package main

import (
	"fmt"
	ds "hw6/documentstore"
)

func main() {
	d1 := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d1.Fields["key1"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "key1"}
	d1.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: "value"}
	cfg1 := ds.CollectionConfig{PrimaryKey: "key1"}
	store := ds.NewStore()
	ok, col1 := store.CreateCollection("key1", &cfg1)
	if !ok {
		fmt.Println("Collection creation failed")
		return
	}
	col1.Put(d1)

	err := store.DumpToFile("store.json")
	if err != nil {
		return
	}

	restoredStore, err := ds.NewStoreFromFile("store.json")
	if err != nil {
		return
	}
	col, ok := restoredStore.GetCollection("key1")
	if ok {
		fmt.Println(col.List())
	}
	fmt.Println("Done")

}
