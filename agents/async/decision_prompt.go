package async

import (
	"encoding/json"
	"fmt"

	llm "github.com/scottraio/plum/llms"
	logger "github.com/scottraio/plum/logger"
)

const ASYNC_DECISION_PROMPT = `
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

// Decision represents a structured decision made by the agent.
type DecisionPrompt struct {
	Question string
	Context  string
	Memory   string
	Tools    string

	Decision Decision
}

// Decision represents a structured decision made by the agent.
type Decision struct {
	Question string   `json:"Question"`
	Thought  string   `json:"Thought"`
	Actions  []Action `json:"Actions"`
}

type Action struct {
	Tool        string `json:"Tool"`
	ToolInput   string `json:"Input"`
	Reasoning   string `json:"Reasoning"`
	Observation string `json:"Observation"`
}

// Decide makes a decision based on the agent's input and memory.
func (a *DecisionPrompt) Decide(agent Agent) Decision {
	prompt := llm.InjectObjectToPrompt(a, ASYNC_DECISION_PROMPT)

	logger.Log("Agent", "Thinking...", "gray")

	decision := agent.LLM.Run(prompt)

	// Parse the JSON response to get the Decision object
	err := json.Unmarshal([]byte(decision), &a.Decision)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	// Verbose logging
	logger.Log("Question", a.Question, "blue")
	logger.Log("Thought", a.Decision.Thought, "gray")

	// Inject the agent's input and memory into the prompt
	return a.Decision
}
