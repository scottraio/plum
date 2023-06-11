package agents

import (
	plum "github.com/scottraio/plum"
	agents "github.com/scottraio/plum/agents"
	"github.com/scottraio/plum/example/skills"
)

// CustomerServiceAgent represents a customer service agent.
func AutoAgent() agents.Agent {
	// Create the agent.
	return plum.Agent(agents.Agent{
		DecisionContext: "You are an expert programmer. You know how to build apps and websites.",
		AnswerContext:   "",
		Tools:           AutoTools(),
		Method:          "sequential_selection",
	})
}

func AutoTools() []agents.Tool {
	// Tools are the actions the agent can take.
	return []agents.Tool{
		{
			Name:        "ShellCommand",
			Description: "Useful for executing shell commands to write software or output code",
			HowTo:       "",
			Func: func(input agents.Input) string {
				return skills.ShellCommand(input.Text)
			},
		},
	}
}
