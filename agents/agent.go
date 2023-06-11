package agents

import (
	decision "github.com/scottraio/plum/agents/decision"
	llm "github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
	memory "github.com/scottraio/plum/memory"
)

// The main Agent object.
// Agents are the way to interact with the LLM.
// Agents are mostly few-shot prompting.
//
// The built in Agents include:
//  1. Chat - A Conversational QA Agent
//  2. Auto - A goal-oriented agent that executes Steps and Actions recursively until a goal is reached.
//  3. ZeroShot - Not a few-shot agent, just a quick way to send a input+prompt=answer.

type Agent struct {
	Name  string
	Input string
	Path  string

	DecisionContext string
	AnswerContext   string

	LLM    llm.LLM
	Memory memory.Memory

	Tools     []Tool // user-defined
	ToolNames []string

	Method   string // user-defined
	Decision decision.Decision

	BelongsTo string // user-defined
	HasMany   []string

	DecisionRules []string
	AnswerRules   []string
}

// Run executes the agent's decision-making process.
func (a *Agent) Answer(input string) string {
	a.Input = input
	a.Memory = memory.Memory{}

	logger.Log(a.Path+"Thinking", input, "cyan")

	decision := a.Decide()

	outputs := a.runActions(decision.Actions)

	for _, output := range outputs {
		a.Memory.Add(output, "assistant")
	}

	a.Memory.Add(a.AnswerContext, "context")

	rules := "Follow these rules:"
	for _, rule := range a.AnswerRules {
		rules += "\n" + rule
	}

	a.Memory.Add(rules, "system")

	// Question
	a.Memory.Add(input, "user")

	answer := a.LLM.Run(a.Memory)

	// Answer
	a.Memory.Add(answer, "answer")
	logger.Log(a.Path+"Answer", answer, "green")

	return answer
}

// Remember stores the agent's memory.
func (a *Agent) Remember(memory memory.Memory) *Agent {
	a.Memory = memory

	return a
}

func (a *Agent) CallingAgent(name string) *Agent {
	a.Path = name + " > " + a.Name + " > "

	return a
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
func (a *Agent) Decide() decision.DecisionResp {
	decisionMethod := decision.GetDecisionMethod(a.Method)

	decide := &decision.Decision{
		Input:        a.Input,
		Context:      a.DecisionContext,
		Instructions: decisionMethod.Instructions(),
		Tools:        DescribeTools(a.Tools),
		Path:         a.Path,
		Rules:        a.DecisionRules,
	}

	return decide.Decide(a.Memory, a.LLM)
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
// func (a *Agent) RunStep(step Step, prompt string) Decision {
// 	input := fmt.Sprintf("Input: %s. Step: %s", a.Input, step.Description)
// 	return a.Decide(input, prompt)
// }

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
func (a *Agent) RunAction(act decision.Action) string {
	var actionResult string

	// TODO: This should be a goroutine
	for _, tool := range a.Tools {
		if tool.Name == act.Tool {
			input := Input{
				CallingAgent: a.Name,
				Text:         act.ToolInput,
				Action:       act,
				Memory:       a.Memory,
				CurrentStep:  act.StepDescription,
				ToolName:     tool.Name,
				ToolHowTo:    tool.HowTo,
				LLM:          a.LLM,
			}

			tool.CallingAgent = a.Name

			actionResult = tool.Func(input)

			if actionResult == "" {
				actionResult = "No output. (Input: " + act.ToolInput + "))"
			}

			// logger.Log("Tool "+act.Tool+" Thought", act.Thought, "gray")
			// logger.Log("Tool "+act.Tool+" Input", act.ToolInput, "gray")
			// logger.Log("Tool "+act.Tool+" Output", actionResult, "gray")
			break
		}
	}

	return actionResult
}

//
// Prompt helper functions
//

// Returns a string of the agent's plans (steps and/or actions).
// func (a *Agent) Plans() string {
// 	mainPlan := ""

// 	if len(a.Decision.Steps) > 0 {
// 		return a.PlansWithSteps(a.Decision.Steps, mainPlan)
// 	} else {
// 		return a.PlansWithActions(a.Decision.Actions, mainPlan)
// 	}
// }

// // Returns a string of the agent's steps.
// func (a *Agent) PlansWithSteps(steps []decision.Step, mainPlan string) string {
// 	for i, step := range steps {
// 		mainPlan += fmt.Sprintf("Step %d. %s", i, step.Description)
// 		a.PlansWithActions(step.Actions, mainPlan)
// 	}

// 	return mainPlan
// }

// // Returns a string of the agent's actions.
// func (a *Agent) PlansWithActions(actions []decision.Action, mainPlan string) string {
// 	for i, action := range actions {
// 		mainPlan += fmt.Sprintf("Action %d. %s", i, action.Tool)
// 	}

// 	return mainPlan
// }

// RunActions runs the actions in the agent's decision.
func (a *Agent) runActions(actions []decision.Action) []string {
	summary := []string{}
	no_actions := len(actions)
	// logger.Log("Number of actions", strconv.Itoa(no_actions), "gray")

	// Create a channel to receive the summaries from each goroutine
	ch := make(chan string, no_actions)

	for _, action := range actions {
		// logger.Log("Tool", action.Tool, "gray")
		//a.Memory.Add(action.ToolInput, "user")

		// Start a new goroutine for each action
		go func(action decision.Action) {
			ch <- a.RunAction(action)
		}(action)
	}

	// Collect the summaries from each goroutine
	for i := 0; i < len(actions); i++ {
		summary = append(summary, <-ch)
	}

	return summary
}
