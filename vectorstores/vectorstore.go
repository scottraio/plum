package vectorstores

type VectorStore interface {
	Client() VectorStore
	Query(input []float32) string
	Upsert(namespace string, input []float32) string
	Index(index string) VectorStore
}
