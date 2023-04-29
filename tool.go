package plum

import (
	llm "github.com/scottraio/plum/llms"
	retriever "github.com/scottraio/plum/retrievers"
)

type Tool struct {
	Name        string
	Description string
	HowTo       string
	Func        func(query retriever.QueryBuilder) string
}

func (t *Tool) Prompt() string {
	template := `
		
		Name: {{.Name}} 
		Reasoning: {{.Description}}
			
		How to use: 
		{{.HowTo}}
		
		----------------------------------------------

		`

	return llm.InjectObjectToPrompt(t, template)
}
