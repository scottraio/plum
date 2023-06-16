package answer

import (
	"fmt"

	"github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/memory"
)

// Decision represents a structured decision made by the agent.
type Answer struct {
	Input      string
	Context    string
	Memory     string
	Tools      string
	Method     AnswerMethod
	Path       string
	Rules      []string
	Truths     []string
	Outputs    string
	ScratchPad []string
}

type AnswerMethod interface {
	Output(string) AnswerMethod
	Format() string
	Instructions() string
	SetQuestion(string)
	Log()
	Validate() bool
	FinalAnswer() string
	GetNotes() string
}

func GetAnswerMethod(method string) AnswerMethod {
	switch method {
	case "scored":
		return &ScoredAnswer{}
	case "summarized":
		return &SummarizedAnswer{}
	default:
		return &SummarizedAnswer{}
	}
}

// Decide makes a decision based on the agent's input and memory.
func (a *Answer) Answer(mem memory.Memory, llm llms.LLM) AnswerMethod {

	mem.Add(a.PromptBackground(), "system")
	mem.Add(a.PromptContext(), "system")
	mem.Add(a.PromptRules(), "system")
	mem.Add(a.Method.Instructions(), "system")
	mem.Add(a.Outputs, "system")
	mem.Add(a.GetScratchPad(), "system")
	mem.Add(a.Input, "user")
	mem.Add("JSON Response:", "system")

	// Run the LLM
	answer := llm.Answer(mem)
	output := a.Method.Output(answer)

	a.Method.Log()
	return output
}

func (a *Answer) PromptContext() string {
	return "Context: " + a.Context
}

func (a *Answer) PromptBackground() string {
	background := "Background: You are a Plum Answer Agent, a powerful language model that can assist with a wide range of questions and provide in-depth explanations and discussions on various topics in context. You can process and understand large amounts of text, generate human-like responses, and provide valuable insights and information. You can make decisions and perform complex decision making and reasoning. As a JSON API, a Plum Answer Agent will carefully craft an answer to the question based on the input received from the user."

	return background
}

func (a *Answer) PromptRules() string {
	rulesPrompt := "Follow these rules: \n"

	rules := []string{"Always respond with valid JSON, do not respond with anything other than JSON."}
	rules = append(rules, a.Truths...)
	rules = append(rules, a.Rules...)

	for i, rule := range rules {
		rulesPrompt += fmt.Sprintf("\n %d. %s", i+1, rule)
	}

	return rulesPrompt
}

func (a *Answer) GetScratchPad() string {
	scratchPad := "Scratch Pad: \n"

	for i, note := range a.ScratchPad {
		scratchPad += fmt.Sprintf("\n %d. %s", i+1, note)
	}

	return scratchPad
}
