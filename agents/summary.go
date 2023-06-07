package agents

import (
	"strings"

	llm "github.com/scottraio/plum/llms"
	"github.com/scottraio/plum/logger"
)

// Summary represents a summary of multiple actions ran by an agent.
type SummaryPrompt struct {
	Context string
	Memory  string
	Summary string
	Input   string
}

// Summary Prompt
// -----------------

const SUMMARY_PROMPT = `
Background
----------
A Plum Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. 

Memory
---------
{{.Memory}}

Knowledge
---------
{{.Summary}}

Instructions
-------------
{{.Context}}

Use the knowledge and memory (if available) to answer the question. DO NOT make up answers.

Answers are helpful, friendly, informative, and confident. Use lists when possible. 

Answer this question in Markdown format: {{.Input}}

Begin!
`

func (a *Agent) summarize(toolOutputs []string) string {
	s := SummaryPrompt{
		Input:   a.Input,
		Context: a.Context,
		Summary: strings.Join(toolOutputs, "\n"),
		Memory:  a.Memory.Format()}

	prompt := llm.InjectObjectToPrompt(s, SUMMARY_PROMPT)
	// Log prompt to log file, do not show in stdout
	logger.PersistLog(prompt)

	return a.LLM.Run(prompt)
}
