package decision

import (
	"encoding/json"
	"fmt"

	"github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
	"github.com/scottraio/plum/memory"
)

// Decision represents a structured decision made by the agent.
type Decision struct {
	Input        string
	Context      string
	Memory       string
	Tools        string
	Instructions string
	Path         string
	Rules        []string
	Truths       []string
	ScratchPad   []string

	DecisionResp DecisionResp
}

type DecisionResp struct {
	Input   string   `json:"Question"`
	Thought string   `json:"Thought"`
	Actions []Action `json:"Actions"`
	_Prompt string
}

type Action struct {
	Tool      string `json:"Tool"`
	ToolInput string `json:"Input"`
	Thought   string `json:"Thought"`
	Notes     string `json:"Notes"`

	StepDescription string
}

type DecisionStrategy interface {
	Instructions() string
}

func GetDecisionStrategy(method string) DecisionStrategy {
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

	mem.Add(d.PromptBackground(), "system")
	mem.Add(d.PromptContext(), "system")
	mem.Add(d.PromptRules(), "system")
	mem.Add(d.PromptInstructions(), "system")
	mem.Add(d.GetScratchPad(), "system")
	mem.Add(d.PromptFormat(), "system")
	mem.Add(d.Input, "user")
	mem.Add("JSON Response:", "system")

	// Run the LLM
	decision := llm.Run(mem)

	// Parse the JSON response to get the Decision object
	err := json.Unmarshal([]byte(decision), &d.DecisionResp)
	if err != nil {
		logger.Log("Error", "There was an error with the response from the LLM, retrying: "+fmt.Sprintf("%v", err)+" original decision: "+decision, "red")
		//d.Decide(mem, llm)
	}

	// set the prompt for future use
	d.DecisionResp._Prompt = decision

	for _, action := range d.DecisionResp.Actions {
		logger.Log("Tool", action.Tool, "yellow")
		logger.Log("Thought", action.Thought, "yellow")
		logger.Log("Input", action.ToolInput, "yellow")
		logger.Log("Notes", action.Notes, "yellow")
	}

	// Inject the agent's input and memory into the prompt
	return d.DecisionResp
}

func (d *Decision) PromptBackground() string {
	background := "Background: You are a Plum Agent, a powerful language model that can assist with a wide range of questions and provide in-depth explanations and discussions on various topics in context. You can process and understand large amounts of text, generate human-like responses, and provide valuable insights and information. You can make decisions and perform complex decision making and reasoning. As a JSON API, a Plum Agent determines the necessary actions to take based on the input received from the user."

	return background
}

func (d *Decision) PromptContext() string {
	return "Context: " + d.Context
}

func (d *Decision) PromptRules() string {
	rulesPrompt := "Follow these rules: \n"
	rules := []string{"Always respond with valid JSON, do not respond with anything other than JSON."}
	rules = append(rules, d.Truths...)
	rules = append(rules, d.Rules...)

	for i, truth := range rules {
		rulesPrompt += fmt.Sprintf("\n %d. %s", i+1, truth)
	}

	return rulesPrompt
}

func (d *Decision) PromptInstructions() string {
	instructions := `
	Follow these instructions to answer the question:
	{{.Instructions}}
	
	Use these tools: 
	{{.Tools}}
	`
	return llms.InjectObjectToPrompt(d, instructions)
}

func (d *Decision) PromptFormat() string {
	format := `Respond with this JSON format: {
		"Question": "{{.Input}}",
		"Thought": "the thought about what action(s) and input(s) are required to answer the question.",
		"Actions": [{
			"Tool": "the tool name to use",
			"Thought": "the thought about what the input to the tool should be",
			"Notes": "Notes on improvements for future prompts",
			"Input": "the input to the tool"
		}]
	}`

	return llms.InjectObjectToPrompt(d, format)
}

func (a *Decision) GetScratchPad() string {
	scratchPad := "Scratch Pad: \n"

	for i, note := range a.ScratchPad {
		scratchPad += fmt.Sprintf("\n %d. %s", i+1, note)
	}

	return scratchPad
}
