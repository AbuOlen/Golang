package documentstore

type Store struct {
	collections map[string]Collection
}

func NewStore() *Store {
	return &Store{collections: make(map[string]Collection)}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створеня то повертаємо `false` та nil
	col := Collection{ docs: make(map[string]Document), config: *cfg }
	if &col == nil {
		return false, nil
	}
	s.collections[name] = col
	return true, &col
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	col, ok := s.collections[name]
	if !ok { return nil, false; }
	return &col, true
}

func (s *Store) DeleteCollection(name string) bool {
	_, ok := s.collections[name]
	if !ok {
		return false
	}
	delete(s.collections, name)
	return true
}