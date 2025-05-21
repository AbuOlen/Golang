package documentstore

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
)

var logger = slog.Default()

type Store struct {
	collections map[string]Collection
}

func (s *Store) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Collections map[string]Collection `json:"collections"`
	}
	alias := Alias{
		Collections: s.collections,
	}

	return json.Marshal(alias)
}

func (s *Store) UnmarshalJSON(data []byte) error {
	// Create an alias or temporary struct for unmarshalling
	alias := struct {
		Collections map[string]Collection `json:"collections"`
	}{}

	// Unmarshal into the alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Set private field manually
	s.collections = alias.Collections
	return nil
}

func NewStore() *Store {
	return &Store{collections: make(map[string]Collection)}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створеня то повертаємо `false` та nil
	if cfg == nil {
		logger.Warn("CollectionConfig is nil, cannot create collection", "name", name)
		return false, nil
	}
	col := Collection{docs: make(map[string]Document), config: *cfg}

	_, exists := s.collections[name]
	if exists {
		logger.Warn("Collection already exists", "name", name)
		return false, nil
	}

	s.collections[name] = col
	logger.Info("Collection created", "name", name)
	return true, &col
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	col, ok := s.collections[name]
	if !ok {

		logger.Warn("Collection not found", "name", name)
		return nil, false
	}
	logger.Info("Collection retrieved", "name", name)
	return &col, true
}

func (s *Store) DeleteCollection(name string) bool {
	_, ok := s.collections[name]
	if !ok {

		logger.Warn("Collection not found for deletion", "name", name)
		return false
	}
	delete(s.collections, name)
	logger.Info("Collection deleted", "name", name)
	return true
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	// Функція повинна створити та проініціалізувати новий `Store`
	// зі всіма колекціями да даними з вхідного дампу.
	var store Store
	err := json.Unmarshal(dump, &store)
	if err != nil {
		return nil, err
	}
	return  &store, nil
}

func (s *Store) Dump() ([]byte, error) {
	// Методи повинен віддати дамп нашого стору в який включені дані про колекції та документ
	data, err := json.Marshal(s)
	if err != nil {
		logger.Error("Failed to marshal store", "error", err)
		return nil, err
	}
	logger.Info("Store dumped to JSON")
	return data, nil
}

// Значення яке повертає метод `store.Dump()` має без помилок оброблятись функцією `NewStoreFromDump`

func NewStoreFromFile(filename string) (*Store, error) {
	// Робить те ж саме що і функція `NewStoreFromDump`, але сам дамп має діставатись з файлу
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = file.Close()
		if err != nil {}
	}()
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	data := make([]byte, fileInfo.Size())
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}

	return NewStoreFromDump(data)
}

func (s *Store) DumpToFile(filename string) error {
	// Робить те ж саме що і метод  `Dump`, але записує у файл замість того щоб повертати сам дамп
	data, err := s.Dump()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Error("Failed to open file for writing", "filename", filename, "error", err)
		}
	}()

	file.Truncate(0)
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
