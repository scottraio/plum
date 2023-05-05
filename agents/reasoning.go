package agents

import (
	"encoding/json"
	"fmt"

	"github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
)

type Decision struct {
	Input   string   `json:"Question"`
	Thought string   `json:"Thought"`
	Actions []Action `json:"Actions"`
	Steps   []Step   `json:"Steps"`
}

type Step struct {
	Description string   `json:"Description"`
	Actions     []Action `json:"Actions"`
	Validate    string   `json:"Validate"`
}

type Action struct {
	Tool      string `json:"Tool"`
	ToolInput string `json:"Input"`
	Validate  string `json:"Validate"`
}

// Decision represents a structured decision made by the agent.
type DecisionPrompt struct {
	Input   string
	Context string
	Memory  string
	Tools   string

	Decision Decision
}

// Summary represents a summary of multiple actions ran by an agent.
type SummaryPrompt struct {
	Context string
	Memory  string
	Summary string
	Input   string
}

// Decide makes a decision based on the agent's input and memory.
func (a *DecisionPrompt) Decide(prompt string, llm llms.LLM) Decision {
	logger.Log("Agent", "Thinking...", "gray")

	decision := llm.Run(prompt)

	// Parse the JSON response to get the Decision object
	err := json.Unmarshal([]byte(decision), &a.Decision)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	// Verbose logging
	logger.Log("Question", a.Input, "blue")
	logger.Log("Thought", a.Decision.Thought, "gray")

	// Inject the agent's input and memory into the prompt
	return a.Decision
}
