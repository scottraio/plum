package chat

import (
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

{{.Context}}

Instructions
------------
To answer the question, you need to create a plan of action by considering which tools to use. 
Then, you will use the selected tools to take the required actions. 

Please choose one or more tools from the following list to take action:
{{.Tools}}

You may use the following information to answer the question:
{{.Memory}}

Respond in the following JSON format:
-------------------------------------
{
	"Question": "{{.Question}}",
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
	outputs := a.RunActions(decision.Actions)
	answer := a.Summarize(outputs)
	logger.Log("Answer", answer, "green")
	return answer
}

// Remember stores the agent's memory.
func (a *ChatAgent) Remember(memory *memory.Memory) agents.Engine {
	a.Agent.Memory = memory
	return agents.Engine(a)
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

func (a *ChatAgent) Summarize(toolOutputs []string) string {
	s := agents.SummaryPrompt{
		Input:   a.Input,
		Context: a.Agent.Context,
		Summary: strings.Join(toolOutputs, "\n"),
		Memory:  a.Agent.Memory.Format()}

	prompt := llm.InjectObjectToPrompt(s, SUMMARY_PROMPT)

	return a.LLM.Run(prompt)
}
