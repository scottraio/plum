package plum

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	llm "github.com/scottraio/plum/llms"
	llms "github.com/scottraio/plum/llms"
)

const DECISION_PROMPT = `
Background
----------
A Plum Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. As a JSON API, a Plum Agent determines the necessary actions to take based on the input 
received from the user. A Plum Agent understands csv, markdown, json, html and plain text.

{{.DecisionContext}}

Instructions
------------
To answer the question, you need to create a plan of action by considering which tools to use. 
Then, you will use the selected tools to take the required actions. 

Please choose one or more tools from the following list to take action:
{{.GetToolsTextForPrompt}}

You may use the following information to answer the question:
{{.PromptMemory}}

Respond in the following JSON format:
-------------------------------------
{
	"Question": "{{.Input}}",
	"Thought": "Think about what action and input are required to answer the question.",
	"Actions": [{
		"Tool": "the tool name to use",
		"Input": "the input to the tool",
	}]
}

Let's get started!
`

// Agent represents an AI agent with decision-making capabilities.
type Agent struct {
	Input           string
	VectorInput     []float32
	Prompt          string
	DecisionPrompt  string
	DecisionContext string
	SummaryContext  string
	Decision        Decision
	Memory          *Memory
	PromptMemory    string
	Tools           []Tool
	ToolNames       []string
	App             AppConfig
}

// Decision represents a structured decision made by the agent.
type Decision struct {
	Question string   `json:"Question"`
	Thought  string   `json:"Thought"`
	Actions  []Action `json:"Actions"`
}
type Action struct {
	Tool        string `json:"Tool"`
	ToolInput   string `json:"Input"`
	Reasoning   string `json:"Reasoning"`
	Observation string `json:"Observation"`
}

// NewAgent creates a new agent with the given input, prompt, memory, and tools.
func NewAgent(decision_context string, summary_context string, tools []Tool) *Agent {
	agent := &Agent{
		App:             GetApp(),
		Input:           "",
		DecisionPrompt:  DECISION_PROMPT,
		DecisionContext: decision_context,
		SummaryContext:  summary_context,
		Memory:          &Memory{},
		Tools:           tools}
	return agent
}

// Run executes the agent's decision-making process.
func (a *Agent) Run(input string, memory *Memory) string {
	a.InjectInputsToDecisionPrompt(input, memory)
	a.Decide(a.App.LLM)

	summary := a.RunActions()

	s := Summary{
		Question:       a.Input,
		SummaryContext: a.SummaryContext,
		Summary:        RemoveCommonWords(strings.Join(summary, "\n")),
		PromptMemory:   a.PromptMemory}

	answer := s.Summarize()

	a.App.Log("Answer", answer, "green")
	return answer
}

// Decide makes a decision based on the agent's input and memory.
func (a *Agent) Decide(llm llms.LLM) Decision {
	a.App.Log("Agent", "Thinking...", "gray")

	decision := llm.Run(a.DecisionPrompt)

	// Parse the JSON response to get the Decision object
	err := json.Unmarshal([]byte(decision), &a.Decision)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	// Verbose logging
	a.App.Log("Question", a.Input, "blue")
	a.App.Log("Thought", a.Decision.Thought, "gray")

	// Inject the agent's input and memory into the prompt
	return a.Decision
}

// RunActions runs the actions in the agent's decision.
func (a *Agent) RunActions() []string {
	summary := []string{}
	no_actions := len(a.Decision.Actions)
	a.App.Log("Number of actions", strconv.Itoa(no_actions), "gray")

	// Create a channel to receive the summaries from each goroutine
	ch := make(chan string, no_actions)

	for _, action := range a.Decision.Actions {
		a.App.Log("Tool", action.Tool, "gray")
		a.App.Log("Tool Input", action.ToolInput, "gray")

		// Start a new goroutine for each action
		go func(action Action) {
			ch <- a.RunAction(action)
		}(action)
	}

	// Collect the summaries from each goroutine
	for i := 0; i < len(a.Decision.Actions); i++ {
		summary = append(summary, <-ch)
	}

	return summary
}

// GetNextAction returns the next action to take.
func (a *Agent) RunAction(act Action) string {
	var actionResult string

	for _, tool := range a.Tools {
		if tool.Name == act.Tool {
			actionResult = tool.Func(act.ToolInput)
			a.App.Log("Tool Output", actionResult, "white")
			break
		}
	}

	return actionResult
}

//
// Prompt functions
//

// GetToolNamesAsJSON returns the agent's tool names as a JSON string.
func (a *Agent) GetToolNamesAsJSON() string {
	toolNamesJSON, err := json.Marshal(a.ToolNames)
	if err != nil {
		return ""
	}
	return string(toolNamesJSON)
}

// GetToolNamesAsJSON returns the agent's tool names as a JSON string.
func (a *Agent) GetToolsTextForPrompt() string {
	prompt := ""
	for _, tool := range a.Tools {
		prompt += tool.Prompt()
	}
	return prompt
}

// InjectInputsToDecisionPrompt injects the agent's input and memory into the decision prompt.
func (a *Agent) InjectInputsToDecisionPrompt(input string, memory *Memory) *Agent {
	a.Input = input
	a.Memory = memory
	a.PromptMemory = memory.Format()

	for _, tool := range a.Tools {
		a.ToolNames = append(a.ToolNames, tool.Name+" ("+tool.Description+")")
	}

	a.DecisionPrompt = llm.InjectObjectToPrompt(a, a.DecisionPrompt)
	return a
}
