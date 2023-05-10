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
	_Prompt string
}

type Step struct {
	Description string `json:"Description"`
	Validate    string `json:"Validate"`
	Actions     []Action
}

type Action struct {
	Tool      string `json:"Tool"`
	ToolInput string `json:"Input"`
	Validate  string `json:"Validate"`

	StepDescription string
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
	logger.Log("Agent", "Thinking...", "cyan")

	// Log prompt to log file, do not show in stdout
	logger.PersistLog(prompt)

	// Run the LLM
	decision := llm.Run(prompt)

	// Parse the JSON response to get the Decision object
	err := json.Unmarshal([]byte(decision), &a.Decision)
	if err != nil {
		logger.Log("Error", "There was an error with the response from the LLM, retrying: "+fmt.Sprintf("%v", err)+" original decision: "+decision, "red")
		a.Decide(prompt, llm)
	}

	// Verbose logging
	logger.Log("Question", a.Input, "cyan")
	logger.Log("Thought", a.Decision.Thought, "cyan")

	// set the prompt for future use
	a.Decision._Prompt = decision

	// Inject the agent's input and memory into the prompt
	return a.Decision
}

// Decide makes a decision based on the agent's input and memory.
func (a *Decision) StepsToString() string {
	steps := ""
	for i, step := range a.Steps {
		steps += fmt.Sprintf("Step %d: %s", i, step.Description)
	}
	return steps
}

func (action *Action) ActionToString() string {
	json, _ := json.Marshal(action)
	return string(json)
}
