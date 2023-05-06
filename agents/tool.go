package agents

import (
	"encoding/json"

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
		<Tool>
			<Name>{{.Name}}</Name>
			<Description>{{.Description}}</Description>
			<How to use>{{.HowTo}}</How to use>
		</Tool>
	`

	return llm.InjectObjectToPrompt(t, template)
}

// GetToolNamesAsJSON returns the agent's tool names as a JSON string.
func GetToolNamesAsJSON(toolNames []string) string {
	toolNamesJSON, err := json.Marshal(toolNames)
	if err != nil {
		return ""
	}
	return string(toolNamesJSON)
}

// GetToolNamesAsJSON returns the agent's tool names as a JSON string.
func DescribeTools(tools []Tool) string {
	prompt := ""
	for _, tool := range tools {
		prompt += tool.Prompt()
	}
	return prompt
}
