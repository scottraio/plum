package zero_shot

import (
	"bytes"
	"html/template"

	"github.com/scottraio/plum/agents"
)

type ZeroShotEngine struct {
	agents.Agent

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

func (z *ZeroShotEngine) Run() string {
	z.InjectInputsToChainPrompt()
	z.Answer.Value = z.LLM.Run(z.Prompt)

	return z.Answer.Value
}

// InjectInputsToDecisionPrompt injects the agent's input and memory into the decision prompt.
func (z *ZeroShotEngine) InjectInputsToChainPrompt() *ZeroShotEngine {
	z.Prompt = z.injectInputsToPrompt(z.PromptTemplate)
	return z
}

// InjectInputsToPrompt injects the agent's input and memory into the prompt.
func (z *ZeroShotEngine) injectInputsToPrompt(prompt string) string {
	tmpl, err := template.New("").Parse(prompt)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, z); err != nil {
		return ""
	}

	return buf.String()
}
