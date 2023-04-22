package vectorstores

type VectorStore interface {
	Client() VectorStore
	Upsert(namespace string, text string) error
	Query(input string) string
	WithNamespace(namespace string) VectorStore
}
