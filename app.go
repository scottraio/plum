package plum

import (
	"fmt"

	embeddings "github.com/scottraio/plum/embeddings"
	llms "github.com/scottraio/plum/llms"
	store "github.com/scottraio/plum/vectorstores"
)

var App *AppConfig

// AppConfig represents the app config.
type AppConfig struct {
	Port               string
	Verbose            bool
	VectorStore        store.VectorStore
	VectorStoreIndexes []string
	LLM                llms.LLM
	Embedding          embeddings.Embedding
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

// Log events with color
func (a *AppConfig) Log(label string, message string, color string) {
	if a.Verbose {
		switch color {
		case "green":
			fmt.Println("\033[1m\033[32m["+label+"]", message, "\033[0m")
			return
		default:
			fmt.Println("\033[1m\033[37m["+label+"]", message, "\033[0m")
			return
		}
	}
}
