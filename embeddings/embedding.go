package embeddings

type Embedding interface {
	EmbedText(text string) []float32
}
