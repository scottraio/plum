package auto

import (
	agents "github.com/scottraio/plum/agents"
	llm "github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
)

func (a *AutoAgent) Decide(input string, prompt string) agents.Decision {
	decision := &agents.DecisionPrompt{
		Input:   input,
		Context: a.Context,
		Memory:  a.Memory.Format(),
		Tools:   agents.DescribeTools(a.Tools)}

	prompt = llm.InjectObjectToPrompt(decision, prompt)
	return decision.Decide(prompt, a.LLM)
}

// RunActions runs the actions in the agent's decision.
func (a *AutoAgent) RunSteps(steps []agents.Step) {
	for _, step := range steps {
		logger.Log("Step Description", step.Description, "gray")
		logger.Log("Step Validation", step.Validate, "gray")

		// Start a new goroutine for each action
		a.RunStep(step)
	}
}

// GetNextAction returns the next action to take.
func (a *AutoAgent) RunStep(step agents.Step) string {
	var actionResult string

	// TODO: This should be a goroutine
	for _, action := range step.Actions {
		a.RunAction(action)
	}

	return actionResult
}

// GetNextAction returns the next action to take.
func (a *AutoAgent) RunAction(act agents.Action) string {
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
