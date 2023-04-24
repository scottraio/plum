package vectorstores

type VectorStore interface {
	Client() VectorStore
	Upsert(text string, fields map[string]string) error
	Query(input string, filter map[string]string) string
	WithNamespace(namespace string) VectorStore
}
