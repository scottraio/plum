package decision

import (
	"encoding/json"
	"fmt"

	"github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
	"github.com/scottraio/plum/memory"
)

const DECISION_PROMPT = `
Follow these instructions to answer the question: 
{{.Instructions}}

Tools: {{.Tools}}

Respond with the following JSON format:
{
	"Question": "{{.Input}}",
	"Thought": "the thought about what action(s) and input(s) are required to answer the question.",
	"Actions": [{
		"Tool": "the tool name to use",
		"Thought": "the thought about what the input to the tool should be",
		"Input": "the input to the tool",
	}]
}

Let's begin!
`

// Decision represents a structured decision made by the agent.
type Decision struct {
	Input        string
	Context      string
	Memory       string
	Tools        string
	Instructions string
	Path         string
	Rules        []string

	DecisionResp DecisionResp
}

type DecisionResp struct {
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
func (d *Decision) Decide(mem memory.Memory, llm llms.LLM) DecisionResp {
	prompt := llms.InjectObjectToPrompt(d, DECISION_PROMPT)

	mem.Add("Background: You are a Plum Agent, a powerful language model that can assist with a wide range of tasks, including answering questions and providing in-depth explanations and discussions on various topics. You can process and understand large amounts of ext, generate human-like responses, and provide valuable insights  and information. You can make decisions and perform complex decision making and reasoning. As a JSON API, a Plum Agent determines the necessary actions to take based on the input received from the user. A Plum Agent understands csv, markdown, json, html and plain text.", "background")
	mem.Add("Context: "+d.Context, "context")
	mem.Add("Instructions: "+prompt, "output_format")

	if len(d.Rules) > 0 {
		rules := "Follow these rules: "
		for _, rule := range d.Rules {
			rules += "\n" + rule
		}

		mem.Add(rules, "system")
	}

	// Log prompt to log file, do not show in stdout
	logger.PersistLog(prompt)

	// Run the LLM
	decision := llm.Run(mem)

	// Parse the JSON response to get the Decision object
	err := json.Unmarshal([]byte(decision), &d.DecisionResp)
	if err != nil {
		logger.Log("Error", "There was an error with the response from the LLM, retrying: "+fmt.Sprintf("%v", err)+" original decision: "+decision, "red")
		d.Decide(mem, llm)
	}

	// set the prompt for future use
	d.DecisionResp._Prompt = decision

	for _, action := range d.DecisionResp.Actions {
		toolSelection := fmt.Sprintf("Selected the %s tool, because %s. Query: %s", action.Tool, action.Thought, action.ToolInput)
		logger.Log(d.Path+"Decision", toolSelection, "yellow")
		mem.Add(toolSelection, "system")
	}

	// Inject the agent's input and memory into the prompt
	return d.DecisionResp
}
