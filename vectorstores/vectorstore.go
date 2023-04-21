package vectorstores

type VectorStore interface {
	Client() VectorStore
	Query(input []float32) string
	Upsert(namespace string, vectors [][]float32) error
	Index(index string) VectorStore
}
