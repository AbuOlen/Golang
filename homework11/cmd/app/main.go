package main

import (
	"fmt"
	ds "hw11/internal/documentstore"
	"math/rand"
	"sync"
	"time"
	"sync/atomic"
)

func testGoRoutines(col *ds.Collection) {
	var wg sync.WaitGroup
	rand.Seed(time.Now().Unix())
	startTime := time.Now()

	var putCounter int64 = 0
	var listCounter int64 = 0

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rand.Intn(10) % 2 == 0 {
				d := ds.Document{Fields: make(map[string]ds.DocumentField)}
				d.Fields["key1"] = ds.DocumentField{Type: ds.DocumentFieldTypeString,
					Value: fmt.Sprintf("key%d", rand.Intn(1000))}
				d.Fields["val"] = ds.DocumentField{Type: ds.DocumentFieldTypeString,
					Value: fmt.Sprintf("val%d", rand.Intn(1000))}
				col.Put(d)
				atomic.AddInt64(&putCounter, 1)
			} else {
				col.List()
				atomic.AddInt64(&listCounter, 1)

			}
		}()
	}

	wg.Wait()
	elapsedTime := time.Since(startTime)

	// Print performance metrics
	fmt.Printf("Performance Metrics:\n")
	fmt.Printf("  Total Put Calls: %d\n", atomic.LoadInt64(&putCounter))
	fmt.Printf("  Total List Calls: %d\n", atomic.LoadInt64(&listCounter))
	fmt.Printf("  Total Execution Time: %v\n", elapsedTime)
}

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
		testGoRoutines(col)
		//fmt.Println(col.List())
	}

	fmt.Println("Done")

}
