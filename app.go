package plum

import (
	embeddings "github.com/scottraio/plum/embeddings"
	llms "github.com/scottraio/plum/llms"
	store "github.com/scottraio/plum/vectorstores"
)

var App *AppConfig

// AppConfig represents the app config.
type AppConfig struct {
	Port              string
	Verbose           bool
	VectorStore       map[string]store.VectorStore
	VectorStoreConfig VectorStoreConfig
	LLM               llms.LLM
	Embedding         embeddings.Embedding
}

type VectorStoreConfig struct {
	Db      string
	Indexes []string
}

// Init initializes the app config.
func Boot(init Initialize) AppConfig {
	App = &AppConfig{
		Port:        GetDotEnvVariable("PORT"),
		Verbose:     GetDotEnvVariable("VERBOSE") == "true",
		Embedding:   InitEmbeddings(init),
		LLM:         InitLLM(init),
		VectorStore: InitVectorStore(init),
	}

	return *App
}

// GetApp returns the app config.
func GetApp() AppConfig {
	return *App
}
