package agents

import (
	"fmt"

	llm "github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
	memory "github.com/scottraio/plum/memory"
)

// The Engine interface is the main interface for all agents.
// It is used on the project-side of the application.
// New Agents can be created by the project by implementing the Engine interface+Agent Struct.
type Engine interface {
	Answer(question string) string
	Remember(memory *memory.Memory) Engine
}

// The main Agent object.
// Agents are the way to interact with the LLM.
// Agents are mostly few-shot prompting.
//
// The built in Agents include:
//  1. Chat - A Conversational QA Agent
//  2. Auto - A goal-oriented agent that executes Steps and Actions recursively until a goal is reached.
//  3. ZeroShot - Not a few-shot agent, just a quick way to send a input+prompt=answer.

type Agent struct {
	Engine

	Input   string
	Context string // user-defined

	LLM    llm.LLM
	Memory *memory.Memory

	Tools     []Tool // user-defined
	ToolNames []string

	Decision Decision
}

// A decision always has an input and a thought.
// It may also either have one of the following: actions or steps.
// the Decide function will return a decision object.
//
//	type Decision struct {
//		Input   string   `json:"Question"`
//		Thought string   `json:"Thought"`
//		Actions []Action `json:"Actions"`
//		Steps   []Step   `json:"Steps"`
//	}
func (a *Agent) Decide(input string, prompt string) Decision {
	decide := &DecisionPrompt{
		Input:   input,
		Context: a.Context,
		Memory:  a.Memory.Format(),
		Tools:   DescribeTools(a.Tools)}

	prompt = llm.InjectObjectToPrompt(decide, prompt)
	return decide.Decide(prompt, a.LLM)
}

// Run step is basically another decision tree.
// It takes the Description and Validates field, forms a sentence, and passes it to the Input field for the Decision.
//
// Objects reference:
//
//	 type Step struct {
//		 Description string `json:"Description"`
//		 Validate    string `json:"Validate"`
//	 }
//
// The prompt arg takes any agent prompt that outputs Actions.
func (a *Agent) RunStep(step Step, prompt string) Decision {
	input := fmt.Sprintf("Input: %s. Step: %s", a.Input, step.Description)
	return a.Decide(input, prompt)
}

// Runs the Action, which in turn runs the user-defined func on the Tool.
//
// Objects reference:
//
//	 type Action struct {
//		 Tool      string `json:"Tool"`
//		 ToolInput string `json:"Input"`
//		 Validate  string `json:"Validate"`
//	 }
//
//	 type Tool struct {
//		 Name        string
//		 Description string
//		 HowTo       string
//		 Func        func(query string) string
//	 }
//
// The output is the output of the Tool Func.
func (a *Agent) RunAction(act Action) string {
	var actionResult string

	// TODO: This should be a goroutine
	for _, tool := range a.Tools {
		if tool.Name == act.Tool {
			input := Input{
				Text:        act.ToolInput,
				Agent:       a,
				Action:      act,
				CurrentStep: act.StepDescription,
			}

			actionResult = tool.Func(input)

			if actionResult == "" {
				actionResult = "No output. (Input: " + act.ToolInput + "))"
			}

			logger.Log("Tool "+act.Tool+" Input", act.ToolInput, "gray")
			logger.Log("Tool "+act.Tool+" Output", actionResult, "gray")
			break
		}
	}

	return actionResult
}
