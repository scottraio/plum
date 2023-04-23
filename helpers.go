package plum

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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

func (a *AppConfig) Vectorize(text string) []float32 {
	return a.Embedding.EmbedText(text)
}

// use godot package to load/read the .env file and
// return the value of the key
func GetDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Could not load .env file")
	}

	return os.Getenv(key)
}

func FatalIfError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
