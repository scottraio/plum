package plum

import (
	llm "github.com/scottraio/plum/llms"
)

type Tool struct {
	Name        string
	Description string
	HowTo       string
	Func        func(query string) string
}

func (t *Tool) Prompt() string {
	template := `
		
		Name: {{.Name}} 
		Description: {{.Description}}
			
		How to use: 
		{{.HowTo}}
		
		----------------------------------------------

		`

	return llm.InjectObjectToPrompt(t, template)
}
