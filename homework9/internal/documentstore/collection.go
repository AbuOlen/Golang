package documentstore

import (
	"encoding/json"
	"fmt"
	"github.com/google/btree"
)

type Collection struct {
	docs map[string]Document
	config CollectionConfig
	index map[string]CollectionIndex
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s Collection) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Docs map[string]Document `json:"docs"`
		Config CollectionConfig `json:"config"`
		Index map[string]CollectionIndex `json:"index"`
	}
	alias := Alias{
		Docs: s.docs,
		Config: s.config,
		Index: s.index,
	}

	return json.Marshal(alias)
}

func (s *Collection) UnmarshalJSON(data []byte) error {
	// Create an alias or temporary struct for unmarshalling
	alias := struct {
		Docs map[string]Document `json:"docs"`
		Config CollectionConfig `json:"config"`
		Index map[string]CollectionIndex `json:"index"`
	}{}

	// Unmarshal into the alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Set private field manually
	s.docs = alias.Docs
	s.config = alias.Config
	s.index = alias.Index
	return nil
}


func (s *Collection) Put(doc Document) {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	keyField, ok := doc.Fields[s.config.PrimaryKey]
	if !ok {
		return
	}
	if keyField.Type != DocumentFieldTypeString {
		return
	}
	key, isString := keyField.Value.(string)
	if isString && len(key) > 0 {
		s.docs[key] = doc
		s.updateIndex(key, doc, false)
	}
}

func (s *Collection) Get(key string) (*Document, bool) {
	doc, ok := s.docs[key]
	return &doc, ok
}

func (s *Collection) Delete(key string) bool {
	doc, ok := s.docs[key]
	if !ok {
		return false
	}
	delete(s.docs, key)
	s.updateIndex(key, doc, true)
	return true
}

func (s *Collection) updateIndex(key string, doc Document, del bool) {
	idx := s.findIndexForDocument(doc)
	if idx == nil {
		return
	}
	val := doc.Fields[key].Value
	// Індексуватись мають тільки поля типу string. Якщо у документа поле має інший тип або взагалі поле відсутнє - воно не попадає в індекс
	if strVal, ok := val.(string); ok {
		if del {
			idx.tree.Delete(StringItem(strVal))
			delete(idx.lut, strVal)
		} else {
			idx.tree.ReplaceOrInsert(StringItem(strVal))
			idx.lut[strVal] = key
		}
	}
}

func (s *Collection) findIndexForDocument(doc Document) *CollectionIndex {
	keys := getKeysFromDoc(doc)
	for _, key := range keys {
		for indexField := range s.index {
			if indexField == key {
				idx := s.index[indexField]
				return &idx
			}
		}
	}
	return nil
}

func (s *Collection) List() []Document {
	values := make([]Document, 0, len(s.docs))

	for _, doc := range s.docs {
		values = append(values, doc)
	}

	return values
}

func getKeysFromDoc(doc Document) []string {
	keys := make([]string, 0) // Prepare a slice for the keys
	for key := range doc.Fields {
		keys = append(keys, key) // Append each key to the slice
	}
	return keys
}

type CollectionIndex struct {
	tree *btree.BTree `json:"tree"`
	lut map[string]string `json:"lut"`
}

func (ci CollectionIndex) MarshalJSON() ([]byte, error) {
	type BTreeData []string
	type Alias struct {
		Tree BTreeData `json:"tree"`
		Lut map[string]string `json:"lut"`
	}
	var btreeData BTreeData
	ci.tree.Ascend(func(item btree.Item) bool {
		key := item.(StringItem)
		btreeData = append(btreeData, string(key)) // Convert BTree into a slice of keys
		return true
	})
	alias := Alias{
		Tree: btreeData,
		Lut: ci.lut,
	}
	return json.Marshal(alias)
}

func (ci *CollectionIndex) UnmarshalJSON(data []byte) error {
	type BTreeData []string
	// Create an alias or temporary struct for unmarshalling
	alias := struct {
		Tree BTreeData `json:"tree"`
		Lut map[string]string `json:"lut"`
	}{}

	// Unmarshal into the alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Set private field manually
	ci.lut = alias.Lut
	ci.tree = btree.New(32)
	for _, key := range alias.Tree {
		ci.tree.ReplaceOrInsert(StringItem(key))
	}
	return nil
}

type StringItem string

// Implement the Less method for B-Tree ordering.
// This determines the order of elements in the tree (lexical sorting for strings).
func (a StringItem) Less(b btree.Item) bool {
	return a < b.(StringItem)
}

func (s *Collection) CreateIndex(fieldName string) error {
	// Якщо індекс вже існує - повертаємо помилку
	if _, ok := s.index[fieldName]; ok {
		return fmt.Errorf( "Index for field %s already exists", fieldName)
	}
	idx := CollectionIndex{ tree: btree.New(32), lut: make(map[string]string) }
	for key, doc := range s.docs {
		v, ok := doc.Fields[fieldName]
		if !ok {
			continue
		}
		val := v.Value
		// Індексуватись мають тільки поля типу string. Якщо у документа поле має інший тип або взагалі поле відсутнє - воно не попадає в індекс
		if strVal, ok := val.(string); ok {
			idx.tree.ReplaceOrInsert(StringItem(strVal))
			idx.lut[strVal] = key
		}
	}
	s.index[fieldName] = idx
	return nil
}

func (s *Collection) DeleteIndex(fieldName string) error {
	if _, ok := s.index[fieldName]; !ok {
		return fmt.Errorf("Index for field %s doesn't exist", fieldName)
	}
	delete(s.index, fieldName)
	return nil
}

type QueryParams struct {
	// TODO: Implement
	Desc bool // Визначає в якому порядку повертати дані
	MinValue *string // Визначає мінімальне значення поля для фільтрації
	MaxValue *string // Визначає максимальне значення поля для фільтрації
}

func (s *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	// TODO: Implement
	// Якщо для даного поля не існує індекса - повертаємо помилку
	idx, ok := s.index[fieldName]
	if !ok {
		return nil, fmt.Errorf("Index for field %s doesn't exist", fieldName)
	}

	// Prepare the result slice
	var result []Document

	// Add filtering logic based on QueryParams
	var min, max btree.Item
	if params.MinValue != nil {
		min = StringItem(*params.MinValue)
	}
	if params.MaxValue != nil {
		max = StringItem(*params.MaxValue)
	}

	// Iterate through the B-Tree with given bounds
	iterator := func(item btree.Item) bool {
		strVal := item.(StringItem)
		// Lookup document key in LUT
		key, exists := idx.lut[string(strVal)]
		if exists {
			if doc, found := s.Get(key); found {
				result = append(result, *doc)
			}
		}
		return true
	}

	if params.Desc {
		// Descending iteration with bounds
		if max != nil && min != nil {
			idx.tree.DescendRange(max, min, iterator)
		} else if max != nil {
			idx.tree.DescendLessOrEqual(max, iterator)
		} else if min != nil {
			idx.tree.Descend(iterator)
		} else {
			idx.tree.Descend(iterator)
		}
	} else {
		// Ascending iteration with bounds
		if min != nil && max != nil {
			idx.tree.AscendRange(min, max, iterator)
		} else if min != nil {
			idx.tree.AscendGreaterOrEqual(min, iterator)
		} else {
			idx.tree.Ascend(iterator)
		}
	}

	return result, nil
}