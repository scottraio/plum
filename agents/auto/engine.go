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
	
	1.Let's think step by step. You will be provided with a goal and will need to create high-level steps to accomplish it. Your knowledge of Tools will help define what's possible.
	2. For each step, you will need to plan the actions required to achieve the goal of the step. Use your memory for local context.
	3. Each action will require a Tool and Input. The Tool is the name of the tool to be used and the Input is the input required by the tool.
	4. If an action uses a tool that will exceed the token limit (3000), branch it to a new prompt by setting Branch to true.
	5. Respond back to the software function with a an efficient plan of actions in JSON.

-------------------------------------------------------------------------------

Let's begin!

{
	"Input": "{{.Input}}",
	"Thought": "Consider the high-level actions required to achieve the goal.",
	"Steps": [{
		"Description": "a step to take",
		"Validate": "how would you validate or test the step?"
		"Actions": [{
			"Tool": "name of the tool to be used",
			"Input": "input required by the tool",
			"Branch": "false"
		}]
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

// RunActions runs the actions in the agent's decision.
func (a *AutoAgent) runSteps(steps []agents.Step) string {
	output := ""
	for _, step := range steps {
		logger.Log("Step Planned", step.Description, "purple")
		for _, action := range step.Actions {
			logger.Log(action.Tool, action.ToolInput, "purple")
		}
	}

	logger.Log("Running", "......", "yellow")

	for _, step := range steps {
		logger.Log("Running Step", step.Description, "yellow")
		logger.Log("Step Validation", step.Validate, "yellow")

		// Start a new goroutine for each action
		for _, action := range step.Actions {
			action.StepDescription = step.Description
			result := a.RunAction(action)
			output = output + result
		}
	}
	return output
}
