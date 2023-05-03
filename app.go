package plum

import (
	"github.com/scottraio/plum/agents"
	llms "github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
	models "github.com/scottraio/plum/models"
	skills "github.com/scottraio/plum/skills"
	util "github.com/scottraio/plum/util"
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
	Embedding         func(input string) []float32
	Models            map[string]*models.Model
	Skills            map[string]*skills.Skill
	Agents            map[string]agents.Engine
	Env               string
}

type VectorStoreConfig struct {
	Db      string
	Indexes []string
}

// Init initializes the app config.
func Boot(init Initialize) AppConfig {
	App = &AppConfig{
		Env:         util.GetDotEnvVariable("PLUM_ENV"),
		Port:        util.GetDotEnvVariable("PORT"),
		Verbose:     util.GetDotEnvVariable("VERBOSE") == "true",
		Embedding:   init.Embedding,
		LLM:         InitLLM(init),
		VectorStore: InitVectorStore(init),
		Models:      make(map[string]*models.Model),
		Skills:      make(map[string]*skills.Skill),
		Agents:      make(map[string]agents.Engine),
	}

	return App.boot()
}

// GetApp returns the app config.
func GetApp() AppConfig {
	return *App
}

// Register Models
func (a *AppConfig) RegisterModel(name string, m *models.Model) {
	logger.Log("Model", "Model "+name+" Registered ", "purple")
	a.Models[name] = m
}

// Register Skills
func (a *AppConfig) RegisterSkill(name string, skill *skills.Skill) {
	logger.Log("Skill", "Skill "+name+" Registered ", "purple")
	a.Skills[name] = skill
}

// Register Agents
func (a *AppConfig) RegisterAgent(name string, ag agents.Engine) {
	logger.Log("Agent", "Agent "+name+" Registered ", "purple")
	a.Agents[name] = ag
}

// bootLog logs the app config.
func (a *AppConfig) boot() AppConfig {
	logger.Log("App", "Plum "+Version, "purple")

	for key := range a.VectorStore {
		logger.Log("Vector Store", "Index "+key+" Registered ", "purple")
	}

	return *a
}
