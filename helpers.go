package plum

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

func CSVToMap(reader *csv.Reader) []map[string]string {
	rows := []map[string]string{}
	var header []string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}

func (a *AppConfig) Vectorize(text string) []float32 {
	return a.Embedding.EmbedText(text)
}

func (a *AppConfig) CSVToVector(content string) [][]float32 {
	vectors := [][]float32{}
	csvReader := csv.NewReader(strings.NewReader(content))
	csvMap := CSVToMap(csvReader)

	// process each row
	for _, row := range csvMap {
		var doc string
		for key, value := range row {
			doc += fmt.Sprintf("%s: %s\n", key, value)
		}

		vector := a.Embedding.EmbedText(doc)
		vectors = append(vectors, vector)
	}

	return vectors
}

func (a *AppConfig) CSVToStrings(content string) []string {
	csvReader := csv.NewReader(strings.NewReader(content))
	csvMap := CSVToMap(csvReader)

	docs := []string{}

	// process each row
	for _, row := range csvMap {
		var doc string
		for key, value := range row {
			doc += fmt.Sprintf("%s: %s\n", key, value)
		}
		docs = append(docs, doc)
	}

	return docs
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
