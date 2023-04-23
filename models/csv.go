package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
)

type Csv struct {
	EmbedFunc func(string) []float32
}

func (c *Csv) CSVToVector(content string) [][]float32 {
	vectors := [][]float32{}
	csvReader := csv.NewReader(strings.NewReader(content))
	csvMap := c.CSVToMap(csvReader)

	// process each row
	for _, row := range csvMap {
		var doc string
		for key, value := range row {
			doc += fmt.Sprintf("%s: %s\n", key, value)
		}

		vector := c.EmbedFunc(doc)
		vectors = append(vectors, vector)
	}

	return vectors
}

func (c *Csv) CSVToStrings(content string) []string {
	csvReader := csv.NewReader(strings.NewReader(content))
	csvMap := c.CSVToMap(csvReader)

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

func (c *Csv) CSVToMap(reader *csv.Reader) []map[string]string {
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
