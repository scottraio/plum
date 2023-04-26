package vectorstores

type VectorStore interface {
	Client() VectorStore
	Upsert(text string, fields map[string]string) error
	Query(input string, filter map[string]string) string
	Purge(namespace string) error
	WithNamespace(namespace string) VectorStore
}
