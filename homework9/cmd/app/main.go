package main

import (
	"fmt"
	ds "hw9/internal/documentstore"
)

func documentWithValue(key string, val string) ds.Document {
	d := ds.Document{Fields: make(map[string]ds.DocumentField)}
	d.Fields["key1"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: key}
	d.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeString, Value: val}
	return d
}

func strPtr(s string) *string {
	return &s
}

func main() {
	cfg1 := ds.CollectionConfig{PrimaryKey: "key1"}
	store := ds.NewStore()
	ok, col1 := store.CreateCollection("key1", &cfg1)
	if !ok {
		fmt.Println("Collection creation failed")
		return
	}
	d0 := documentWithValue("key0", "val0")
	col1.Put(d0)
	d1 := documentWithValue("key1", "val1")
	col1.Put(d1)
	d2 := documentWithValue("key2", "val2")
	col1.Put(d2)
	d3 := documentWithValue("key3", "val3")
	col1.Put(d3)
	d4 := documentWithValue("key4", "val4")
	col1.Put(d4)
	d5 := documentWithValue("key5", "val5")
	col1.Put(d5)

	col1.CreateIndex("val")

	docs, qerr := col1.Query("val", ds.QueryParams{Desc: true, MinValue: strPtr("val2"), MaxValue: strPtr("val4")})
	if qerr != nil {
		fmt.Println(qerr)
		return
	}
	fmt.Println(docs)

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
		docs, qerr := col.Query("val", ds.QueryParams{Desc: false, MinValue: strPtr("val2"), MaxValue: strPtr("val4")})
		if qerr != nil {
			fmt.Println(qerr)
			return
		}
		fmt.Println(docs)

		col.Delete("key3")
		d6 := documentWithValue("key6", "val6")
		col1.Put(d6)

		docs, qerr = col.Query("val", ds.QueryParams{Desc: false, MinValue: strPtr("val2"), MaxValue: strPtr("val6")})
		if qerr != nil {
			fmt.Println(qerr)
			return
		}
		fmt.Println(docs)
	}
	fmt.Println("Done")

}
