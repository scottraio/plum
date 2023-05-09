package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
)

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

func FatalIfEmpty(label string, msg string) {
	if msg == "" {
		logger.LogAndFail(label+"is required", msg, "red")
	}
}

func InjectObjectToPrompt(obj interface{}, prompt string) string {
	return llms.InjectObjectToPrompt(obj, prompt)
}
