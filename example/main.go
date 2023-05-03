package example

import (
	exampleAgents "github.com/scottraio/plum/example/agents"
	"github.com/scottraio/plum/memory"

	plum "github.com/scottraio/plum"

	models "github.com/scottraio/plum/example/models"
)

func main() {
	// Initialize the app config.
	boot := plum.Boot(plum.Initialize{
		Embedding: plum.InitEmbeddings("openai"),
		LLM:       "openai",
		VectorStoreConfig: plum.VectorStoreConfig{
			Db:      "pinecone",
			Indexes: []string{"knowledge", "structured"}},
	})

	// Register the models.
	// TODO: Automatically register models.
	boot.RegisterModel("knowledge", models.Manual())

	// Register the agents.
	// TODO: Automatically register agents.
	boot.RegisterAgent("customer_service", exampleAgents.CustomerServiceAgent())

	history := []memory.ChatHistory{
		{
			Query:  "I need help with my [Product Name]",
			Answer: "I'm sorry to hear that. What seems to be the problem?",
		},
	}

	memory := memory.LoadMemory(history)

	// Call the agent with the input and chat history.
	agent := plum.App.Agents["customer_service"]
	agent.Remember(memory).Answer("I need help with my [Product Name]")
}
