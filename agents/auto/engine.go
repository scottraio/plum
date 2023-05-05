package auto

import (
	agents "github.com/scottraio/plum/agents"
	"github.com/scottraio/plum/logger"
	memory "github.com/scottraio/plum/memory"
)

// Agent represents an AI agent with decision-making capabilities.
type AutoAgent struct {
	agents.Engine
	agents.Agent
}

const DECISION_PROMPT = `
Background
----------
A Plum Auto Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. As a JSON API, a Plum Auto Agent determines the necessary actions to take based on the input 
received from the user. 

{{.Context}}

Instructions
------------
You will be given a goal and will abstract the high-level steps to accomplish the goal. 


Please choose one or more tools from the following list to take action:
{{.Tools}}

You may use the following information to answer the question:
{{.Memory}}

Respond in the following JSON format:
-------------------------------------
{
	"Input": "{{.Input}}",
	"Thought": "Think about the high-level action to take",
	"Steps": [{
		"Description": "the high-level step to take",
		"Validate": "how would you validate or test the step?"

		"Actions": [{
			"Tool": "the tool name to use",
			"Input": "the input to the tool",
			"Validate": "the validation input to the tool",
		}]
	}]
}

Let's get started!
`

// Engine Interface Functions
// -----------------

// Run executes the agent's decision-making process.
func (a *AutoAgent) Answer(input string) string {
	a.Input = input

	decision := a.Decide(input, DECISION_PROMPT)

	logger.Log("Answer", decision, "green")

	return answer
}

// Remember stores the agent's memory.
func (a *AutoAgent) Remember(memory *memory.Memory) agents.Engine {
	a.Agent.Memory = memory
	return agents.Engine(a)
}
