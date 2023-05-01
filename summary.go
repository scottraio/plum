package plum

import (
	"strings"

	llm "github.com/scottraio/plum/llms"
)

const SUMMARY_PROMPT = `
Background
----------
A Plum Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. 

As a summarization function of the Plum Agent, you will understand the research and memory, to accurately respond to the question. 
received from the user.

{{.SummaryContext}}}

Research
---------
You've done the research this is what you found: 
{{.Summary}}

Memory
---------
You remember the following:
{{.PromptMemory}}

Instructions
-------------
(legend: + = style/formatting, * = required for answer)

* Use the research and memory (if available) to answer the question.
* DO NOT make up answers, if you're unsure say "Hmmm... I'm not sure".

+ You are helpful and friendly.
+ Answer in markdown format.
+ Answer in a natural conversational, friendly tone.
+ Use lists when possible.
+ Return the answer ONLY.
+ Respond is a confident, concise, and clear manner.
+ Avoid generic phrases like "based on the information I have" or "according to my research."


Begin!

Answer this question: {{.Question}}
`

// Summary represents a summary of multiple actions ran by an agent.
type Summary struct {
	SummaryContext string `json:"summary_context"`
	PromptMemory   string `json:"prompt_memory"`
	Summary        string `json:"summary"`
	Question       string `json:"question"`
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
func (s *Summary) Summarize() string {
	SummarizedPrompt := llm.InjectObjectToPrompt(s, SUMMARY_PROMPT)
	return App.LLM.Run(SummarizedPrompt)
}
