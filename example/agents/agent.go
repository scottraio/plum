package agents

import (
	plum "github.com/scottraio/plum"
	retriever "github.com/scottraio/plum/retrievers"
)

const AGENT_PROMPT = `
You are the official [Company Name] customer service assistant. Help and assist me with troubleshooting, guiding, and answer questions on [Company Name] products only.
You are given the following extracted snippets of a many documents and a question. 
If you are unsure of the answer, say "Hmm... I'm not sure".
Answers should be conversational and helpful. Use lists as much as possible. Respond in markdown.

Memory: {{.PromptMemory}}
--------------------
Question: {{.Question}}
--------------------
Summary: {{.Summary}}
--------------------
Helpful Answer:`

// CustomerServiceAgent represents a customer service agent.
func CustomerServiceAgent(input string, memory plum.Memory) string {
	// Get the app config.
	app := plum.GetApp()

	// Tools are the actions the agent can take.
	tools := []plum.Tool{
		{
			Name:        "Part Number Lookup",
			Description: "Useful for finding information about parts.",
			HowTo:       "Use the part number to find the answer",
			Func: func(query retriever.QueryBuilder) string {
				return app.VectorStore["structured"].Query(input, nil, nil)
			},
		},
		{
			Name:        "General Info",
			Description: "Useful for finding general information",
			HowTo:       "Use the knowledge base to find the answer",
			Func: func(query retriever.QueryBuilder) string {
				lookup := plum.App.Models["knowledge"].Return(query)
				return lookup
			},
		},
	}

	// Create the agent.
	agent := plum.NewAgent(AGENT_PROMPT, tools, ``)

	// Run the agent.
	return agent.Run(input, &memory)
}
