package agents

import (
	answer "github.com/scottraio/plum/agents/answer"
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

	Context string

	LLM        llm.LLM
	Memory     memory.Memory
	ScratchPad []string

	Tools     []Tool // user-defined
	ToolNames []string

	Strategy string // user-defined
	Decision decision.Decision

	Method string // user-defined
	Output answer.Answer

	BelongsTo string // user-defined
	HasMany   []string

	Truths []string

	DecisionRules []string
	AnswerRules   []string
}

// Run executes the agent's decision-making process.
func (a *Agent) Answer(input string) string {
	a.Input = input
	// TODO: an agent should carry the conversational memory forward.
	a.Memory = memory.Memory{}

	// Make a decision (runs LLM)
	decisionResp := a.Decide()

	// answer the question
	outputs := a.PromptOutputs(decisionResp.Actions)
	return a.Return(outputs)
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
	logger.Log(a.Path+"Thinking", a.Input, "cyan")

	decisionStrat := decision.GetDecisionStrategy(a.Strategy)

	decide := &decision.Decision{
		Input:        a.Input,
		Context:      a.Context,
		Instructions: decisionStrat.Instructions(),
		Tools:        DescribeTools(a.Tools),
		Path:         a.Path,
		Rules:        a.DecisionRules,
		Truths:       a.Truths,
		ScratchPad:   a.ScratchPad,
	}

	decisionResp := decide.Decide(a.Memory, a.LLM)

	return decisionResp
}

func (a *Agent) PromptOutputs(actions []decision.Action) string {
	outputs := a.runActions(actions)

	outputsPrompt := "Outputs: "
	for _, output := range outputs {
		outputsPrompt += "\n" + output
	}

	return outputsPrompt
}

func (a *Agent) Return(outputs string) string {
	logger.Log(a.Path+"Answering", a.Input, "cyan")

	answerMethod := answer.GetAnswerMethod(a.Method)
	answerMethod.SetQuestion(a.Input)

	answer := &answer.Answer{
		Input:      a.Input,
		Context:    a.Context,
		Method:     answerMethod,
		Path:       a.Path,
		Rules:      a.AnswerRules,
		Truths:     a.Truths,
		Outputs:    outputs,
		ScratchPad: a.ScratchPad,
	}

	answerResp := answer.Answer(a.Memory, a.LLM)
	a.ScratchPad = append(a.ScratchPad, "Previous answer example: "+answerResp.FinalAnswer())

	if !answerResp.Validate() {
		return a.Answer(a.Input)
	}

	return answerResp.FinalAnswer()
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
func (a *Agent) RunAction(act decision.Action) string {
	var actionResult string

	// TODO: This should be a goroutine
	for _, tool := range a.Tools {
		if tool.Name == act.Tool {
			input := Input{
				CallingAgent:  a.Name,
				Text:          act.ToolInput,
				Action:        act,
				Memory:        a.Memory,
				CurrentStep:   act.StepDescription,
				ToolName:      tool.Name,
				ToolInputType: tool.InputType,
				LLM:           a.LLM,
			}

			tool.CallingAgent = a.Name

			actionResult = tool.Func(input)
			a.ScratchPad = append(a.ScratchPad, act.Notes)

			if actionResult == "" {
				actionResult = "No output. (Input: " + act.ToolInput + "))"
			}
			break
		}
	}

	return actionResult
}

// RunActions runs the actions in the agent's decision.
func (a *Agent) runActions(actions []decision.Action) []string {
	summary := []string{}

	for _, action := range actions {
		summary = append(summary, a.RunAction(action))
	}

	return summary
}
