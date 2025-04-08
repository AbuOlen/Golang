package documentstore

type Collection struct {
	docs map[string]Document
	config CollectionConfig
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s *Collection) Put(doc Document) {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	key_field, ok := doc.Fields[s.config.PrimaryKey]
	if ok && key_field.Type == "string" {
		s.docs[s.config.PrimaryKey] = doc
	}
}

func (s *Collection) Get(key string) (*Document, bool) {
	doc, ok := s.docs[key]
	return &doc, ok
}

func (s *Collection) Delete(key string) bool {
	_, ok := s.docs[key]
	if !ok {
		return false
	}
	delete(s.docs, key)
	return true
}

func (s *Collection) List() []Document {
	values := make([]Document, 0, len(s.docs))

	for _, doc := range s.docs {
		values = append(values, doc)
	}

	return values
}