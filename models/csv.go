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

func (c *Csv) CSVMatrixPinnedFeaturesToStrings(csvString string) []string {
	csvReader := csv.NewReader(strings.NewReader(csvString))
	csvData, _ := csvReader.ReadAll()

	columns := make([]string, len(csvData[0])-2)

	for i := 1; i < len(csvData[0])-1; i++ {
		columns[i-1] = csvData[0][i]
	}

	featureNames := csvData[1:]

	docs := []string{}

	for _, featureValue := range featureNames {
		for modelIndex, modelName := range columns {
			doc := ""
			featureName := featureValue[0]
			value := featureValue[modelIndex+1]
			doc = doc + fmt.Sprintf("%s %s: %s ", modelName, featureName, value)

			docs = append(docs, doc)
		}
	}

	return docs
}

func (c *Csv) CSVMatrixPinnedHeadersToStrings(content string) []string {
	csvReader := csv.NewReader(strings.NewReader(content))
	csvData, _ := csvReader.ReadAll()

	columnNames := make([]string, len(csvData[0])-2)

	for i := 1; i < len(csvData[0])-1; i++ {
		columnNames[i-1] = csvData[0][i]
	}

	featureNames := csvData[1:]

	docs := []string{}
	for modelIndex, modelName := range columnNames {

		for _, featureValue := range featureNames {
			var doc string
			featureName := featureValue[0]
			value := featureValue[modelIndex+1]
			doc = fmt.Sprintf("%s - %s: %s \n", modelName, featureName, value)
			docs = append(docs, doc)
		}

	}

	return docs
}
