package plum

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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
