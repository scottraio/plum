package async

import (
	"strconv"
	"strings"

	logger "github.com/scottraio/plum/logger"
	"github.com/scottraio/plum/memory"

	agents "github.com/scottraio/plum/agents"
)

// Agent represents an AI agent with decision-making capabilities.
type Agent struct {
	agents.Engine
	agents.Agent
	Decision Decision
}

// Engine Interface Functions
// -----------------

// Run executes the agent's decision-making process.
func (a *Agent) Answer(input string) string {
	a.Input = input

	decision := a.decide(input)
	outputs := a.runActions(decision.Actions)
	answer := a.summarize(outputs)

	logger.Log("Answer", answer, "green")
	return answer
}

// Remember stores the agent's memory.
func (a *Agent) Remember(memory *memory.Memory) agents.Engine {
	a.Agent.Memory = memory
	return agents.Engine(a)
}

// Private Functions
// -----------------

func (a *Agent) decide(input string) Decision {
	decision := &DecisionPrompt{
		Question: input,
		Context:  a.Agent.Context,
		Memory:   a.Agent.Memory.Format(),
		Tools:    agents.DescribeTools(a.Agent.Tools)}

	return decision.Decide(*a)
}

func (a *Agent) summarize(toolOutputs []string) string {
	s := Summary{
		Question: a.Input,
		Context:  a.Agent.Context,
		Summary:  strings.Join(toolOutputs, "\n"),
		Memory:   a.Agent.Memory.Format()}

	return s.Summarize(*a)
}

// RunActions runs the actions in the agent's decision.
func (a *Agent) runActions(actions []Action) []string {
	summary := []string{}
	no_actions := len(actions)
	logger.Log("Number of actions", strconv.Itoa(no_actions), "gray")

	// Create a channel to receive the summaries from each goroutine
	ch := make(chan string, no_actions)

	for _, action := range actions {
		logger.Log("Tool", action.Tool, "gray")
		logger.Log("Tool Input", action.ToolInput, "gray")

		// Start a new goroutine for each action
		go func(action Action) {
			ch <- a.runAction(action)
		}(action)
	}

	// Collect the summaries from each goroutine
	for i := 0; i < len(actions); i++ {
		summary = append(summary, <-ch)
	}

	return summary
}

// GetNextAction returns the next action to take.
func (a *Agent) runAction(act Action) string {
	var actionResult string

	// TODO: This should be a goroutine
	for _, tool := range a.Agent.Tools {
		if tool.Name == act.Tool {
			actionResult = tool.Func(act.ToolInput)
			logger.Log("Tool Output", actionResult, "white")
			break
		}
	}

	return actionResult
}
