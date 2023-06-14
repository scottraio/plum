package agents

import (
	"encoding/json"
	"fmt"

	decision "github.com/scottraio/plum/agents/decision"
	llm "github.com/scottraio/plum/llms"
	memory "github.com/scottraio/plum/memory"
)

type Tool struct {
	Name         string
	Description  string
	InputType    string
	CallingAgent string
	Func         func(input Input) string
}

func (t *Tool) Prompt() string {
	template := "%s Description: %s Input: %s \n"
	return fmt.Sprintf(template, t.Name, t.Description, t.InputTypeInstructions())
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
	for i, tool := range tools {
		prompt += fmt.Sprintf("\n %d. Name %s Description: %s Input: %s", i+1, tool.Name, tool.Description, tool.InputTypeInstructions())
	}
	return prompt
}

// ToolInput represents the input to a tool.
type Input struct {
	Text          string
	CallingAgent  string
	Action        decision.Action
	Memory        memory.Memory
	Plans         string
	CurrentStep   string
	ToolName      string
	ToolInputType string
	LLM           llm.LLM
}

func (t *Tool) InputTypeInstructions() string {
	switch t.InputType {
	case "text":
		return "A plain-text input is required."
	default:
		return "A plain-text input is required."
	}
}
