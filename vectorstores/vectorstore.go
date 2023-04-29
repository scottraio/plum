package vectorstores

type VectorStore interface {
	Client() VectorStore
	Upsert(text string, fields map[string]string) error
	Query(input string, filter map[string]string, options map[string]interface{}) string
	Purge(namespace string) error
	WithNamespace(namespace string) VectorStore
}

// Vector is a struct that represents a vector that can be stored in a vector store.
type Vector struct {
	Text     string
	MetaData map[string]string
}
