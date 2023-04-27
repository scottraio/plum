package plum

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Log events with color
func (a *AppConfig) Log(label string, message string, color string) {
	msg := strings.ReplaceAll(message, "\n", "")
	var colorStart string
	var colorEnd string

	defaultColorStart := "\033[1m\033[37m"
	defaultColorEnd := "\033[0m"

	// Get current date and time
	now := time.Now()

	// Format the date and time as desired
	dateTimeStr := now.Format("2006-01-02 15:04:05")

	if a.Verbose {
		switch color {
		case "purple":
			colorStart = "\033[1m\033[35m"
			colorEnd = "\033[0m"
		case "green":
			colorStart = "\033[1m\033[32m"
			colorEnd = "\033[0m"
		case "lightblue":
			colorStart = "\033[1m\033[36m"
			colorEnd = "\033[0m"
		case "gray":
			colorStart = "\033[1m\033[30m"
			colorEnd = "\033[0m"
		case "orange":
			colorStart = "\033[1m\033[33m"
			colorEnd = "\033[0m"
		default:
			colorStart = defaultColorStart
			colorEnd = defaultColorEnd
		}

		fmt.Println(colorStart, "["+dateTimeStr+"]", "["+label+"]", msg, colorEnd)
	}
}

func (a *AppConfig) Vectorize(text string) []float32 {
	return a.Embedding(text)
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
