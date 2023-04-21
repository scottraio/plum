package example

import (
	plum "github.com/scottraio/plum"
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
		plum.UseTool(
			"Part Number Lookup",
			"Useful for finding information about parts.",
			func(input string) string {
				return app.VectorStore.Index("structured").Query(app.Embedding.EmbedText(input))
			}),
		plum.UseTool(
			"General Info",
			"Useful for finding general information",
			func(input string) string {
				return app.VectorStore.Index("knowledge").Query(app.Embedding.EmbedText(input))
			}),
	}

	// Create the agent.
	agent := plum.NewAgent(input, AGENT_PROMPT, memory, tools)

	// Run the agent.
	return agent.Run()
}
