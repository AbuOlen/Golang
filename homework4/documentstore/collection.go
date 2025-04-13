package documentstore

type Collection struct {
	docs   map[string]Document
	config CollectionConfig
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s *Collection) Put(doc Document) {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	keyField, ok := doc.Fields[s.config.PrimaryKey]
	if ok && keyField.Type == DocumentFieldTypeString {
		key, isString := keyField.Value.(string)
		if isString && len(key) > 0 {
			s.docs[key] = doc
		}
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
