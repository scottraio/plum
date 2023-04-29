package llms

import (
	"bytes"
	"text/template"
)

type LLM interface {
	Client() LLM
	Run(prompt string) string
}

// InjectInputsToPrompt injects the agent's input and memory into the prompt.
func InjectObjectToPrompt(obj interface{}, prompt string) string {
	tmpl, err := template.New("").Parse(prompt)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, obj); err != nil {
		return ""
	}

	return buf.String()
}
