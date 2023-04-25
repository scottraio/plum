package plum

import (
	embeddings "github.com/scottraio/plum/embeddings"
	llms "github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/models"
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
	Models            map[string]models.Model
	Jobs              map[string]Job
	Env               string
}

type VectorStoreConfig struct {
	Db      string
	Indexes []string
}

// Init initializes the app config.
func Boot(init Initialize) AppConfig {
	App = &AppConfig{
		Env:         GetDotEnvVariable("PLUM_ENV"),
		Port:        GetDotEnvVariable("PORT"),
		Verbose:     GetDotEnvVariable("VERBOSE") == "true",
		Embedding:   InitEmbeddings(init),
		LLM:         InitLLM(init),
		VectorStore: InitVectorStore(init),
		Models:      make(map[string]models.Model),
		Jobs:        make(map[string]Job),
	}

	return App.boot()
}

// GetApp returns the app config.
func GetApp() AppConfig {
	return *App
}

// Register Models
func (a *AppConfig) RegisterModel(name string, m models.Model) {
	a.Log("Model", "Model "+name+" Registered ", "purple")
	a.Models[name] = m
}

// Register Models
func (a *AppConfig) RegisterJob(name string, j Job) {
	a.Log("Job", "Job "+name+" Registered ", "purple")
	a.Jobs[name] = j
}

// bootLog logs the app config.
func (a *AppConfig) boot() AppConfig {
	a.Log("App", "Plum "+Version, "purple")

	for key, _ := range a.VectorStore {
		a.Log("Vector Store", "Index "+key+" Registered ", "purple")
	}

	return *a
}
