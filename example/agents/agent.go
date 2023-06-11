package agents

import (
	plum "github.com/scottraio/plum"
	agents "github.com/scottraio/plum/agents"
)

// CustomerServiceAgent represents a customer service agent.
func CustomerServiceAgent() agents.Agent {
	// Create the agent.
	return plum.Agent(agents.Agent{
		DecisionContext: "You are the official [Company Name] customer service assistant. Help and assist me with troubleshooting, guiding, and answer questions on [Company Name] products only.",
		AnswerContext:   "",
		Tools:           CustomerServiceTools(),
		Method:          "single_selection",
	})
}

func CustomerServiceTools() []agents.Tool {
	// Tools are the actions the agent can take.
	return []agents.Tool{
		{
			Name:        "OrderNumberLookup",
			Description: "Useful for finding tracking information, order status, and more",
			HowTo:       "Use the order status and tracking information to find the answer.",
			Func: func(input agents.Input) string {
				return plum.App.VectorStore["structured"].Query(input.Text, nil, nil)
			},
		},
		{
			Name:        "ProductManuals",
			Description: "Useful for finding general information about our products",
			HowTo:       "Use the information returned to find the answer.",
			Func: func(input agents.Input) string {
				lookup := plum.App.Models["manual"].Return(input.Text)
				return lookup
			},
		},
	}
}
