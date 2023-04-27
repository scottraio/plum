package example

import (
	plum "github.com/scottraio/plum"
)

func main() {
	// Initialize the app config.
	plum.Boot(plum.Initialize{
		Embedding: plum.InitEmbeddings("openai"),
		LLM:       "openai",
		VectorStoreConfig: plum.VectorStoreConfig{
			Db:      "pinecone",
			Indexes: []string{"knowledge", "structured"}},
	})

	// chat history is a slice of ChatHistory structs. Handled by the client.
	chatHistory := plum.LoadMemory(plum.ChatHistory{
		Query:  "I need help with my [Product Name]",
		Answer: "I'm sorry to hear that. What seems to be the problem?",
	})

	// Call the agent with the input and chat history.
	CustomerServiceAgent("I need help with my [Product Name]", chatHistory)
}
