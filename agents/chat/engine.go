package chat

import (
	"strconv"
	"strings"

	"github.com/scottraio/plum/agents"
	llm "github.com/scottraio/plum/llms"
	logger "github.com/scottraio/plum/logger"
	"github.com/scottraio/plum/memory"
)

// Agent represents an AI agent with decision-making capabilities.
type ChatAgent struct {
	agents.Engine
	agents.Agent
	Decision agents.Decision
}

const DECISION_PROMPT = `
Background
----------
A Plum Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. As a JSON API, a Plum Agent determines the necessary actions to take based on the input 
received from the user. A Plum Agent understands csv, markdown, json, html and plain text.

Context
-------
{{.Context}}

Tools
-----
{{.Tools}}

Memory
------
{{.Memory}}

Instructions
------------
To answer the question, you need to create a plan of action by considering which tools to use. 
Then, you will use the selected tools to take the required actions. 

Respond in the following JSON format:
-------------------------------------
{
	"Question": "{{.Input}}",
	"Thought": "Think about what action and input are required to answer the question.",
	"Actions": [{
		"Tool": "the tool name to use",
		"Input": "the input to the tool",
	}]
}

Let's get started!
`

// Engine Interface Functions
// -----------------

// Run executes the agent's decision-making process.
func (a *ChatAgent) Answer(input string) string {
	a.Input = input
	decision := a.Decide(input, DECISION_PROMPT)

	outputs := a.runActions(decision.Actions)

	answer := a.summarize(outputs)
	logger.Log("Answer", answer, "green")
	return answer
}

// Remember stores the agent's memory.
func (a *ChatAgent) Remember(memory *memory.Memory) agents.Engine {
	a.Agent.Memory = memory
	return agents.Engine(a)
}

// RunActions runs the actions in the agent's decision.
func (a *ChatAgent) runActions(actions []agents.Action) []string {
	summary := []string{}
	no_actions := len(actions)
	logger.Log("Number of actions", strconv.Itoa(no_actions), "gray")

	// Create a channel to receive the summaries from each goroutine
	ch := make(chan string, no_actions)

	for _, action := range actions {
		logger.Log("Tool", action.Tool, "gray")
		logger.Log("Tool Input", action.ToolInput, "gray")

		// Start a new goroutine for each action
		go func(action agents.Action) {
			ch <- a.RunAction(action)
		}(action)
	}

	// Collect the summaries from each goroutine
	for i := 0; i < len(actions); i++ {
		summary = append(summary, <-ch)
	}

	return summary
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

As a summarization function of the Plum Agent, you will understand the research and memory, to accurately respond to the question. 
received from the user.

{{.Context}}}

Research
---------
You've done the research this is what you found: 
{{.Summary}}

Memory
---------
You remember the following:
{{.Memory}}

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

Answer this question: {{.Input}}
`

func (a *ChatAgent) summarize(toolOutputs []string) string {
	s := agents.SummaryPrompt{
		Input:   a.Input,
		Context: a.Agent.Context,
		Summary: strings.Join(toolOutputs, "\n"),
		Memory:  a.Agent.Memory.Format()}

	prompt := llm.InjectObjectToPrompt(s, SUMMARY_PROMPT)
	// Log prompt to log file, do not show in stdout
	logger.PersistLog(prompt)

	return a.LLM.Run(prompt)
}
