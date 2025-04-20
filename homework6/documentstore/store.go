package documentstore

import (
	"bufio"
	"encoding/json"
	"os"
	"log/slog"
)

type Store struct {
	Collections map[string]Collection
	logger *slog.Logger
}

func NewStore() *Store {
	s := Store{Collections: make(map[string]Collection)}
	log := slog.Default()
	s.SetLogger(log)
	s.GetLogger().Info("Created")
	return &s
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створеня то повертаємо `false` та nil
	if cfg == nil {
		s.GetLogger().Warn("CollectionConfig is nil, cannot create collection", "name", name)
		return false, nil
	}
	col := Collection{Docs: make(map[string]Document), Config: *cfg}

	_, exists := s.Collections[name]
	if exists {
		s.GetLogger().Warn("Collection already exists", "name", name)
		return false, nil
	}

	s.Collections[name] = col
	s.GetLogger().Info("Collection created", "name", name)
	return true, &col
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	col, ok := s.Collections[name]
	if !ok {

		s.GetLogger().Warn("Collection not found", "name", name)
		return nil, false
	}
	s.GetLogger().Info("Collection retrieved", "name", name)
	return &col, true
}

func (s *Store) DeleteCollection(name string) bool {
	_, ok := s.Collections[name]
	if !ok {

		s.GetLogger().Warn("Collection not found for deletion", "name", name)
		return false
	}
	delete(s.Collections, name)
	s.GetLogger().Info("Collection deleted", "name", name)
	return true
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	// Функція повинна створити та проініціалізувати новий `Store`
	// зі всіма колекціями да даними з вхідного дампу.
	store := Store{Collections: make(map[string]Collection)}
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
		s.GetLogger().Error("Failed to marshal store", "error", err)
		return nil, err
	}
	s.GetLogger().Info("Store dumped to JSON")
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
			s.GetLogger().Error("Failed to open file for writing", "filename", filename, "error", err)
		}
	}()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func (s *Store) SetLogger(logger *slog.Logger) {
	s.logger = logger
}

func (s *Store) GetLogger() *slog.Logger {
	if s.logger == nil {
		s.logger = slog.Default()
	}
	return s.logger
}
