package chat

import (
	"strconv"

	agents "github.com/scottraio/plum/agents"
	llm "github.com/scottraio/plum/llms"
	logger "github.com/scottraio/plum/logger"
)

func (a *ChatAgent) Decide(input string, prompt string) agents.Decision {
	decision := &agents.DecisionPrompt{
		Input:   input,
		Context: a.Context,
		Memory:  a.Memory.Format(),
		Tools:   agents.DescribeTools(a.Tools)}

	prompt = llm.InjectObjectToPrompt(decision, prompt)
	return decision.Decide(prompt, a.LLM)
}

// GetNextAction returns the next action to take.
func (a *ChatAgent) RunAction(act agents.Action) string {
	var actionResult string

	// TODO: This should be a goroutine
	for _, tool := range a.Tools {
		if tool.Name == act.Tool {
			actionResult = tool.Func(act.ToolInput)
			logger.Log("Tool Output", actionResult, "white")
			break
		}
	}

	return actionResult
}

// RunActions runs the actions in the agent's decision.
func (a *ChatAgent) RunActions(actions []agents.Action) []string {
	summary := []string{}
	no_actions := len(actions)
	logger.Log("Number of actions", strconv.Itoa(no_actions), "gray")

	// Create a channel to receive the summaries from each goroutine
	ch := make(chan string, no_actions)

	for _, action := range actions {
		logger.Log("Tool", action.Tool, "gray")
		logger.Log("Tool Input", action.ToolInput, "gray")

		// Start a new goroutine for each action
		go func(action agents.Action) {
			ch <- a.RunAction(action)
		}(action)
	}

	// Collect the summaries from each goroutine
	for i := 0; i < len(actions); i++ {
		summary = append(summary, <-ch)
	}

	return summary
}
