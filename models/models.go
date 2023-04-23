package models

import (
	"context"
	"strings"

	store "github.com/scottraio/plum/vectorstores"
)

type Model struct {
	Name        string
	Attributes  []string
	VectorStore store.VectorStore

	Train  func(ctx context.Context) []string
	Return func(input string, namespace string) string
}

func (m *Model) GetAttributes() map[string]string {
	flat := strings.Join(m.Attributes, ", ")

	return map[string]string{
		"attributes": flat,
	}
}

func (m *Model) TrainModel(namespace string) error {
	var err error
	ctx := context.Background()
	docs := m.Train(ctx)

	for _, doc := range docs {
		err = m.VectorStore.Upsert(namespace, doc, m.GetAttributes())
	}

	return err
}
