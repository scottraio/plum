package models

import (
	"context"

	plum "github.com/scottraio/plum"
	models "github.com/scottraio/plum/models"
	retriever "github.com/scottraio/plum/retrievers"
	store "github.com/scottraio/plum/vectorstores"
	"google.golang.org/genproto/googleapis/type/datetime"
)

type KnowledgeModel struct {
	models.Model

	// Meta attributes you want to store about the model
	Title     string            `json:"Title"`
	CreatedAt datetime.DateTime `json:"CreatedAt"`
}

func Knowledge() *models.Model {
	var knowledge *KnowledgeModel

	// create the model
	knowledge = &KnowledgeModel{
		// Model is the base model that you want to use
		Model: models.Model{
			Name: "Comparison",
			// VectorStore is the vector store that you want to use for this model
			VectorStore: plum.App.VectorStore["knowledge"],

			HowTo: "",

			// The Train/Return functions are the functions that you want to use for training and returning results
			// You decide how to train and how to retrieve the results
			// Plum provides a few functions that you can use to help you with this

			// Train is a function that returns the data to be used for training
			Train: func(ctx context.Context) []store.Vector {
				return knowledge.Train(ctx)
			},

			// Return is a function that returns the result that you want to use in your prompt
			// How the result is used is up to you
			Return: func(queryBuilder retriever.QueryBuilder) string {
				return knowledge.Return(queryBuilder)
			},
		},
	}

	return &knowledge.Model
}

// Fetch gets the data from the source
func (c *KnowledgeModel) Train(ctx context.Context) []store.Vector {
	var results []store.Vector
	// Return an array of strings that you want to use for training
	// The []strings will be used to create the vectors
	// The results are stored in the "text" parameter in the vector store
	return results
}

// Return gets the result from the vector store
func (c *KnowledgeModel) Return(queryBuilder retriever.QueryBuilder) string {
	query := c.Model.VectorStore.Query(
		queryBuilder.Query,
		c.Model.SetAttributes(queryBuilder.Filters),
		queryBuilder.Options,
	)
	return query
}
