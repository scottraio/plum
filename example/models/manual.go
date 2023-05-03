package models

import (
	"context"

	plum "github.com/scottraio/plum"
	models "github.com/scottraio/plum/models"
	store "github.com/scottraio/plum/vectorstores"
)

type ManualModel struct {
	models.Model
	models.Markdown
}

func Manual() *models.Model {
	var manual *ManualModel

	// create the model
	manual = &ManualModel{
		// Model is the base model that you want to use
		Model: models.Model{
			Name: "Manual",

			// VectorStore is the vector store that you want to use for this model
			VectorStore: plum.App.VectorStore["knowledge"],

			// How to understand the data
			HowTo: `
				You are given excerpts from many markdown files. 
				Summarize the data and return a synopsis of the query. 
			`,

			// Return the data to be used for training, the vectors will be stored in the vector store
			Train: func(ctx context.Context) []store.Vector {
				// fetch source documents
				// split documents into chunks
				// build list of Vectors from chunks
				return []store.Vector{}
			},

			// Return result that you want to use in your prompt
			Return: func(query string) string {
				// Query the vector store
				result := manual.Find(query, map[string]string{}, map[string]interface{}{
					"TopK": float64(3),
				})

				// return the results, included with Model name and HowTo read the results
				return manual.Describe(query, result)
			},
		},
	}

	return &manual.Model
}
