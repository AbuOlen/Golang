package documentstore

type Collection struct {
	Docs  map[string]Document
	Config CollectionConfig
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s *Collection) Put(doc Document) {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	keyField, ok := doc.Fields[s.Config.PrimaryKey]
	if ok && keyField.Type == DocumentFieldTypeString {
		key, isString := keyField.Value.(string)
		if isString && len(key) > 0 {
			s.Docs[key] = doc
		}
	}
}

func (s *Collection) Get(key string) (*Document, bool) {
	doc, ok := s.Docs[key]
	return &doc, ok
}

func (s *Collection) Delete(key string) bool {
	_, ok := s.Docs[key]
	if !ok {
		return false
	}
	delete(s.Docs, key)
	return true
}

func (s *Collection) List() []Document {
	values := make([]Document, 0, len(s.Docs))

	for _, doc := range s.Docs {
		values = append(values, doc)
	}

	return values
}
