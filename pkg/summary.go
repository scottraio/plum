package plum

import (
	"bytes"
	"text/template"
)

// Summary represents a summary of multiple actions ran by an agent.
type Summary struct {
	PromptMemory string `json:"prompt_memory"`
	Summary      string `json:"summary"`
	Question     string `json:"question"`
}

// Summarize summarizes the prompt.
func (s *Summary) Summarize(prompt string) string {
	app := GetApp()
	SummarizedPrompt := s.injectSummaryToPrompt(prompt)

	return app.OpenAI.Run(SummarizedPrompt)
}

// injectSummaryToPrompt injects the summary to the prompt.
func (s *Summary) injectSummaryToPrompt(prompt string) string {
	tmpl, err := template.New("").Parse(prompt)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s); err != nil {
		return ""
	}

	return buf.String()
}
