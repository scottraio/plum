package plum

import (
	agents "github.com/scottraio/plum/agents"
)

func Agent(agent agents.Agent) agents.Agent {
	agent.LLM = App.LLM
	agent.Truths = App.Truths
	return agent
}
