package models

import (
	"context"

	store "github.com/scottraio/plum/vectorstores"
)

type Model struct {
	Name        string
	Attributes  map[string]string
	VectorStore store.VectorStore

	Train  func(ctx context.Context) []string
	Return func(input string, filters map[string]string) string
}

func (m *Model) Filters(filters map[string]string) map[string]string {
	return m.GetAttributes(filters)
}

func (m *Model) GetAttributes(attrs map[string]string) map[string]string {
	if attrs == nil {
		attrs = make(map[string]string)
	}

	for key, value := range m.Attributes {
		attrs[key] = value
	}

	attrs["type"] = m.Name
	return attrs
}

func (m *Model) SetAttribute(key string, value string) error {
	m.Attributes[key] = value
	return nil
}

func (m *Model) TrainModel(attrs map[string]string) error {
	var err error
	ctx := context.Background()
	docs := m.Train(ctx)

	for _, doc := range docs {
		err = m.VectorStore.Upsert(doc, m.GetAttributes(attrs))
	}

	return err
}
