package plum

import (
	agents "github.com/scottraio/plum/agents"
	async "github.com/scottraio/plum/agents/async"
)

func AsyncAgent(context string, tools []agents.Tool) agents.Engine {
	agent := &async.Agent{
		Agent: agents.Agent{
			Context: context,
			Tools:   tools,
			LLM:     App.LLM,
		},
	}

	return agents.Engine(agent)
}
