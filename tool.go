package plum

import (
	models "github.com/scottraio/plum/models"
)

type Tool struct {
	Name        string
	Description string
	Func        func(query *models.QueryBuilder) string
}

func UseTool(name string, desc string, run func(query *models.QueryBuilder) string) Tool {
	return Tool{
		Name:        name,
		Description: desc,
		Func:        run,
	}
}
