package auto

import (
	agents "github.com/scottraio/plum/agents"
	"github.com/scottraio/plum/logger"
	"github.com/scottraio/plum/memory"
)

// Agent represents an AI agent with decision-making capabilities.
type AutoAgent struct {
	agents.Agent
	agents.Engine
}

const AUTO_PROMPT = `
Background:
You are a Plum Agent. A Plum Agent is a highly capable language model that excels at answering questions and providing detailed explanations on various topics. It leverages its ability to process and comprehend vast amounts of text to generate human-like responses and offer valuable insights and information. A Plum Agent is designed to return an input back to a software function, all output are valid JSON format.

Context:
{{.Context}}

Tools:
{{.Tools}}

Instructions:
	
	1.You will be provided with a goal and will need to create high-level steps to accomplish it. Your knowledge of <Tools> will aid you in mapping the steps accordingly.

-------------------------------------------------------------------------------

Let's begin!

{
	"Input": "{{.Input}}",
	"Thought": "Consider the high-level actions required to achieve the goal.",
	"Steps": [{
		"Description": "a high-level step to take",
		"Validate": "how would you validate or test the step?"
	}]
}
`

// Engine Interface Functions
// -----------------

// Run executes the agent's decision-making process.
func (a *AutoAgent) Answer(input string) string {
	a.Input = input

	decision := a.Decide(input, AUTO_PROMPT)
	a.runSteps(decision.Steps)

	logger.Log("Answer", "Done", "green")

	return "Done"
}

// Remember stores the agent's memory.
func (a *AutoAgent) Remember(memory *memory.Memory) agents.Engine {
	a.Agent.Memory = memory
	return agents.Engine(a)
}

const STEP_PROMPT = `
Background:
You are a Plum Agent. A Plum Agent is a highly capable language model that excels at answering questions and providing detailed explanations on various topics. It leverages its ability to process and comprehend vast amounts of text to generate human-like responses and offer valuable insights and information. A Plum Agent is designed to return an input back to a software function, all output are valid JSON format.

Context:
{{.Context}}

Tools:
{{.Tools}}

Memory:
{{.Memory}}

Instructions:

	1. Create a plan of actions given an input and step. Use your memory for local context.  
	2. Plan each action using the selected tools from the list. 
	3. Respond back to the software function with a an efficient plan of actions in JSON. 

-------------------------------------------------------------------------------

Let's begin!

{
	"Input": "{{.Input}}",
	"Thought": "Consider which tools are needed and what inputs are required to answer the question.",
	"Actions": [{
		"Tool": "name of the tool to be used",
		"Input": "input required by the tool",
	}]
}
`

// RunActions runs the actions in the agent's decision.
func (a *AutoAgent) runSteps(steps []agents.Step) string {
	output := ""
	for _, step := range steps {
		logger.Log("Step Created", step.Description, "yellow")
	}
	for _, step := range steps {
		logger.Log("Running Step", step.Description, "yellow")
		logger.Log("Step Validation", step.Validate, "yellow")

		// Start a new goroutine for each action
		decision := a.RunStep(step, STEP_PROMPT)
		for _, action := range decision.Actions {
			result := a.RunAction(action)
			a.Agent.Memory.Add(action.ToolInput, result)
			output = output + result
		}
	}
	return output
}
