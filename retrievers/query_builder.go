package retriever

import "encoding/json"

// QueryBuilder is a struct that represents a query that can be sent to a vector store.
type QueryBuilder struct {
	Query   string                 `json:"Query"`
	Filters map[string]string      `json:"Filters"`
	Options map[string]interface{} `json:"Options"`
}

func (qb *QueryBuilder) ToString() string {
	toolInput, _ := json.Marshal(qb)
	return string(toolInput)
}
