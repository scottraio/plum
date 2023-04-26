package plum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	llms "github.com/scottraio/plum/llms"
)

const DECISION_PROMPT = `
Instructions:
1. You are a JSON api that determines what action or actions, to take based on the question and tools.

2. Each action can only be one of the following tools: 
{{.GetToolsAsText}} 

3. You remember the following: 
{{.PromptMemory}}


Respond with the following JSON format:

{
	"Question": "{{.Input}}",
	"Thought": "think about what actions do i need to take to answer the question?",
	"Actions": [
		{
			"Tool": "the tool to use",
			"Reasoning": "the reasoning for using the tool",
			"Input": "input to the tool"
		}
	]	
}


Begin!
`

// Agent represents an AI agent with decision-making capabilities.
type Agent struct {
	Input          string
	VectorInput    []float32
	Prompt         string
	DecisionPrompt string
	Decision       Decision
	Memory         *Memory
	PromptMemory   string
	Tools          []Tool
	ToolNames      []string
	App            AppConfig
}

// Decision represents a structured decision made by the agent.
type Decision struct {
	Question string   `json:"Question"`
	Thought  string   `json:"Thought"`
	Actions  []Action `json:"Actions"`
}
type Action struct {
	Tool      string `json:"Tool"`
	ToolInput string `json:"Input"`
	Reasoning string `json:"Reasoning"`
}

// NewAgent creates a new agent with the given input, prompt, memory, and tools.
func NewAgent(prompt string, tools []Tool) *Agent {
	agent := &Agent{
		App:            GetApp(),
		Input:          "",
		DecisionPrompt: DECISION_PROMPT,
		Prompt:         prompt,
		Memory:         &Memory{},
		Tools:          tools}
	return agent
}

// Run executes the agent's decision-making process.
func (a *Agent) Run(input string, memory *Memory) string {
	a.InjectInputsToDecisionPrompt(input, memory)

	a.Decide(a.App.LLM)

	summary := a.RunActions()

	s := Summary{
		Question:     a.Input,
		Summary:      strings.Join(summary, "\n"),
		PromptMemory: a.PromptMemory}

	answer := s.Summarize(a.Prompt)

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
		a.App.Log("Tool Reasoning", action.Reasoning, "gray")
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
func (a *Agent) GetToolsAsText() string {
	prompt := ""
	for _, tool := range a.Tools {
		prompt += tool.Name + " (" + tool.Description + ")\n"
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

	a.DecisionPrompt = a.injectInputsToPrompt(a.DecisionPrompt)
	return a
}

// InjectInputsToPrompt injects the agent's input and memory into the prompt.
func (a *Agent) injectInputsToPrompt(prompt string) string {
	tmpl, err := template.New("").Parse(prompt)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, a); err != nil {
		return ""
	}

	return buf.String()
}
