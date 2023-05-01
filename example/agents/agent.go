package agents

import (
	plum "github.com/scottraio/plum"
)

const DECISION_CONTEXT = `
You are the official [Company Name] customer service assistant. 
Help and assist me with troubleshooting, guiding, and answer questions on [Company Name] products only.
`

const SUMMARY_CONTEXT = `
You are the official Proluxe customer service AI assistant. 
- You are an expert on commercial food-service equipment, specifically Proluxe equipment.
- You understand the needs of QSR, pizzerias, mexican restaurants, and other food-service establishments.
- You understand the difference between parts, models, and serial numbers.
`

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
			Func: func(query string) string {
				return app.VectorStore["structured"].Query(input, nil, nil)
			},
		},
		{
			Name:        "General Info",
			Description: "Useful for finding general information",
			HowTo:       "Use the knowledge base to find the answer",
			Func: func(query string) string {
				lookup := plum.App.Models["knowledge"].Return(query)
				return lookup
			},
		},
	}

	// Create the agent.
	agent := plum.NewAgent(DECISION_CONTEXT, SUMMARY_CONTEXT, tools)

	// Run the agent.
	return agent.Run(input, &memory)
}
