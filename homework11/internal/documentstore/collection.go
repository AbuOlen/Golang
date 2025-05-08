package documentstore

import (
	"encoding/json"
	"sync"
)

type Collection struct {
	docs map[string]Document
	config CollectionConfig
	mx sync.RWMutex
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s Collection) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Docs map[string]Document `json:"docs"`
		Config CollectionConfig `json:"config"`
	}
	alias := Alias{
		Docs: s.docs,
		Config: s.config,
	}

	return json.Marshal(alias)
}

func (s *Collection) UnmarshalJSON(data []byte) error {
	// Create an alias or temporary struct for unmarshalling
	alias := struct {
		Docs map[string]Document `json:"docs"`
		Config CollectionConfig `json:"config"`
	}{}

	// Unmarshal into the alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Set private field manually
	s.docs = alias.Docs
	s.config = alias.Config
	s.mx = sync.RWMutex{}
	return nil
}


func (s *Collection) Put(doc Document) {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	keyField, ok := doc.Fields[s.config.PrimaryKey]
	if !ok {
		return
	}
	s.mx.Lock()
	defer func() {
		s.mx.Unlock()
	}()
	if keyField.Type != DocumentFieldTypeString {
		return
	}
	key, isString := keyField.Value.(string)
	if isString && len(key) > 0 {
		s.docs[key] = doc
	}
}

func (s *Collection) Get(key string) (*Document, bool) {
	s.mx.RLock()
	defer func() {
		s.mx.RUnlock()
	}()
	doc, ok := s.docs[key]
	return &doc, ok
}

func (s *Collection) Delete(key string) bool {
	s.mx.Lock()
	defer func() {
		s.mx.Unlock()
	}()
	_, ok := s.docs[key]
	if !ok {
		return false
	}
	delete(s.docs, key)
	return true
}

func (s *Collection) List() []Document {
	values := make([]Document, 0, len(s.docs))
	s.mx.RLock()
	defer func() {
		s.mx.RUnlock()
	}()
	for _, doc := range s.docs {
		values = append(values, doc)
	}

	return values
}
