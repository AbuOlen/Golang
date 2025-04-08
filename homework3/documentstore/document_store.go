package documentstore

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = map[string]Document{}

func Put(doc Document) {
	// 1. Перевірити що документ містить в мапі поле `key` типу `string`
	// 2. Додати Document до локальної мапи з документами
	var key, ok = doc.Fields["key"].Value.(string)
	if ok {
		documents[key] = doc
	}
}

func Get(key string) (*Document, bool) {
	// Потрібно повернути документ по ключу
	// Якщо документ знайдено, повертаємо `true` та поінтер на документ
	// Інакше повертаємо `false` та `nil`
	doc, ok := documents[key]
	return &doc, ok
}

func Delete(key string) bool {
	// Видаляємо документа по ключу.
	// Повертаємо `true` якщо ми знайшли і видалили документі
	// Повертаємо `false` якщо документ не знайдено
	_, ok := documents[key]
	if ok {
		delete(documents, key)
		return true
	}
	return false
}

func List() []Document {
	// Повертаємо список усіх документів
	values := make([]Document, 0, len(documents))

	for _, doc := range documents {
		values = append(values, doc)
	}

	return values
}
