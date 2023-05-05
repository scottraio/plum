package plum

import (
	agents "github.com/scottraio/plum/agents"
	auto "github.com/scottraio/plum/agents/auto"
	chat "github.com/scottraio/plum/agents/chat"
)

func ChatAgent(context string, tools []agents.Tool) agents.Engine {
	agent := &chat.ChatAgent{
		Agent: agents.Agent{
			Context: context,
			Tools:   tools,
			LLM:     App.LLM,
		},
	}

	return agents.Engine(agent)
}

func AutoAgent(context string, tools []agents.Tool) agents.Engine {
	agent := &auto.AutoAgent{
		Agent: agents.Agent{
			Context: context,
			Tools:   tools,
			LLM:     App.LLM,
		},
	}

	return agents.Engine(agent)
}
