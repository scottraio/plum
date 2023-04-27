package plum

import (
	"bytes"
	"strings"
	"text/template"
)

// Summary represents a summary of multiple actions ran by an agent.
type Summary struct {
	PromptMemory string `json:"prompt_memory"`
	Summary      string `json:"summary"`
	Question     string `json:"question"`
}

func RemoveCommonWords(summary string) string {
	var simple string

	listOfCommonWords := []string{
		"the ",
		"a ",
		"an ",
		"is ",
		"are ",
		"was ",
		"were ",
		"has ",
		"have ",
		"had ",
		"been ",
		"to ",
		"of ",
		"for ",
	}

	// replace common words from summary like "the" and "a"
	for _, word := range listOfCommonWords {
		simple = strings.ReplaceAll(summary, word, "")
	}

	return simple
}

// Summarize summarizes the prompt.
func (s *Summary) Summarize(prompt string) string {
	SummarizedPrompt := s.injectSummaryToPrompt(prompt)

	return App.LLM.Run(SummarizedPrompt)
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
