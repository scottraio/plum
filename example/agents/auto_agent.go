package agents

import (
	plum "github.com/scottraio/plum"
	agents "github.com/scottraio/plum/agents"
)

// CustomerServiceAgent represents a customer service agent.
func AutoAgent() agents.Engine {
	// Create the agent.
	return plum.AutoAgent(`
		You are an expert programmer. You know how to build apps and websites.
	`, CustomerServiceTools())
}

func AutoTools() []agents.Tool {
	// Tools are the actions the agent can take.
	return []agents.Tool{
		{
			Name:        "ShellCommand",
			Description: "Useful for executing shell commands to write software or output code",
			HowTo:       "",
			Func: func(input agents.Input) string {
				return plum.App.Skills["shell"].Return(input.Text)
			},
		},
	}
}
