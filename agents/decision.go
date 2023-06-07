package agents

import (
	"encoding/json"
	"fmt"

	"github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
)

const DECISION_PROMPT = `
Background
----------
A Plum Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. As a JSON API, a Plum Agent determines the necessary actions to take based on the input 
received from the user. A Plum Agent understands csv, markdown, json, html and plain text.

Tools
-----
{{.Tools}}

Memory
------
{{.Memory}}


Follow these instructions to answer the question:
-------------------------------------------------
{{.Context}}

{{.Instructions}}


Respond in the following JSON format:
-------------------------------------
{
	"Question": "{{.Input}}",
	"Thought": "the thought about what action(s) and input(s) are required to answer the question.",
	"Actions": [{
		"Tool": "the tool name to use",
		"Thought": "the thought about what the input to the tool should be",
		"Input": "the input to the tool",
	}]
}

Let's get started!
`

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
	Thought   string `json:"Thought"`

	StepDescription string
}

// Decision represents a structured decision made by the agent.
type DecisionPrompt struct {
	Input        string
	Context      string
	Memory       string
	Tools        string
	Instructions string

	Decision Decision
}

type DecisionMethod interface {
	Instructions() string
}

func GetDecisionMethod(method string) DecisionMethod {
	switch method {
	case "parallel":
		return &ParallelDecision{}
	case "single":
		return &SingleDecision{}
	case "iteration":
		return &IterationDecision{}
	case "sequential":
		return &SequentialDecision{}
	case "multiple_selection":
		return &MultipleSelectionDecision{}
	default:
		return &ParallelDecision{}
	}
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
