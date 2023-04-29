package plum

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	llm "github.com/scottraio/plum/llms"
	llms "github.com/scottraio/plum/llms"
	retriever "github.com/scottraio/plum/retrievers"
)

const DECISION_PROMPT = `
You are a JSON api that sdetermines what action, actions or no action, to take based on the question and tools.

Instructions:
-------------
1. You will create a plan of action by thinking about what actions do i need to take to answer the question. 

2. Each action queries a model at a time, you can use the same tool twice.

3. Each action can only be one of the following tools: 
{{.GetToolsTextForPrompt}} 

4. You remember the following: 
{{.PromptMemory}}


Respond with the following JSON format:
---------------------------------------
{
	"Question": "{{.Input}}",
	"Thought": "think about what actions and inputs do i need to take to answer the question?",
	"Actions": [
		{
			"Tool": "the tool to use",
			"Reasoning": "the reasoning for using the tool",
			"Input": {
				"Query": "keywords that describe the question", 
				// omit if empty
				"Filters": { "meta_key": "value for meta key" }, 
				// omit if empty
				"Options": { "topK": "value as int" }
			}
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
	Tool      string                 `json:"Tool"`
	ToolInput retriever.QueryBuilder `json:"Input"`
	Reasoning string                 `json:"Reasoning"`
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
		Summary:      RemoveCommonWords(strings.Join(summary, "\n")),
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

		a.App.Log("Tool Input", action.ToolInput.ToString(), "gray")

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
