package models

import (
	"context"
	"strings"

	"github.com/schollz/progressbar/v3"
	store "github.com/scottraio/plum/vectorstores"
)

type Model struct {
	VectorStore store.VectorStore

	Name       string            `json:"Name"`
	Attributes map[string]string `json:"Attributes"`

	HowTo  string
	Train  func(ctx context.Context) []store.Vector
	Return func(query string) string
}

func (m *Model) SetAttributes(filters map[string]string) map[string]string {
	attrs := m.GetAttributes()

	for key, value := range filters {
		attrs[key] = value
	}

	m.Attributes = attrs

	return attrs
}

// GetAttributes returns the attributes for the model
func (m *Model) GetAttributes() map[string]string {
	attrs := make(map[string]string, len(m.Attributes)+1)

	for key, value := range m.Attributes {
		attrs[key] = value
	}

	attrs["type"] = m.Name
	return attrs
}

// SetAttribute sets an attribute for the model
func (m *Model) SetAttribute(key string, value string) error {
	m.Attributes[key] = value
	return nil
}

func (m *Model) TrainModel(attrs map[string]string) error {
	var err error
	ctx := context.Background()
	vectors := m.Train(ctx)

	// Create a new progress bar with the total number of documents
	bar := progressbar.Default(int64(len(vectors)))

	// Iterate over the documents and insert them into the vector store
	for _, vector := range vectors {
		for k, v := range vector.MetaData {
			attrs[k] = v
		}

		err = m.VectorStore.Upsert(vector.Text, attrs, map[string]interface{}{
			"Namespace": m.Namespace(),
		})

		// Increment the progress bar after each iteration
		bar.Add(1)
	}

	return err
}

func (m *Model) Purge() error {
	return m.VectorStore.Purge(m.Namespace())
}

func (m *Model) Find(query string, filters map[string]string, opts map[string]interface{}) string {
	opts["Namespace"] = m.Namespace()

	return m.VectorStore.Query(query, filters, opts)
}

func (m *Model) Describe(query string, result string) string {
	return `
		Model Name: ` + m.Name + `
		Model Context: ` + m.HowTo + `

		Model Query (generated from Question): ` + query + `
		Model Summary: ` + result + `
		------------------------------------------------------------
	`
}

func (m *Model) Namespace() string {
	return strings.ToLower(m.Name)
}
