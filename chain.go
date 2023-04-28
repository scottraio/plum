package plum

import (
	"bytes"
	"text/template"
)

type Chain struct {
	Input          string
	LLM            interface{}
	Memory         Memory
	PromptTemplate string
	Prompt         string
	Action         func(input string) string
	Answer         Answer
}

type Answer struct {
	Question     string `json:"question"`
	Context      string `json:"context"`
	Value        string `json:"value"`
	PromptMemory string `json:"prompt_memory"`
}

func (c *Chain) Result() string {
	c.Answer.Question = c.Input
	c.Answer.Context = c.Action(c.Input)
	c.Answer.PromptMemory = c.Memory.Format()

	return c.Run()
}

func (c *Chain) Run() string {
	c.InjectInputsToChainPrompt()
	c.Answer.Value = App.LLM.Run(c.Prompt)

	return c.Answer.Value
}

// InjectInputsToDecisionPrompt injects the agent's input and memory into the decision prompt.
func (c *Chain) InjectInputsToChainPrompt() *Chain {
	c.Prompt = c.injectInputsToPrompt(c.PromptTemplate)
	return c
}

// InjectInputsToPrompt injects the agent's input and memory into the prompt.
func (c *Chain) injectInputsToPrompt(prompt string) string {
	tmpl, err := template.New("").Parse(prompt)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, c.Answer); err != nil {
		return ""
	}

	return buf.String()
}
