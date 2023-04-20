# The Plum Framework



```go

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

func CustomerServiceAgent(input string, memory base.Memory) string {
	tools := []base.Tool{
		base.UseTool(
			"Part Number Lookup",
			"Useful for finding information about parts.",
			func(input string) string {
				pinecone := base.Pinecone{
					IndexName: "parts-db",
				}

				return pinecone.Query(input)
			}),
		base.UseTool(
			"General Info",
			"Useful for finding general information",
			func(input string) string {
				pinecone := base.Pinecone{
					IndexName: "knowledge-db",
				}

				return pinecone.Query(input)
			}),
	}

	agent := base.NewAgent(input, AGENT_PROMPT, memory, tools)

	return agent.Run()
}

```