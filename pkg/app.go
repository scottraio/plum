package framework

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var App *AppConfig

// AppConfig represents the app config.
type AppConfig struct {
	PineconeKey string
	PineconeEnv string
	Pinecone    Pinecone
	OpenAIToken string
	OpenAI      OpenAI
	Port        string
	Verbose     bool
}

// Init initializes the app config.
func Init() AppConfig {
	App = &AppConfig{
		PineconeKey: getDotEnvVariable("PINECONE_API_KEY"),
		PineconeEnv: getDotEnvVariable("PINECONE_ENV"),
		OpenAIToken: getDotEnvVariable("OPENAI_API_KEY"),
		Port:        getDotEnvVariable("PORT"),
		Verbose:     getDotEnvVariable("VERBOSE") == "true",
	}

	App.OpenAI = *NewOpenAI(App.OpenAIToken)

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

// use godot package to load/read the .env file and
// return the value of the key
func getDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Could not load .env file")
	}

	return os.Getenv(key)
}
