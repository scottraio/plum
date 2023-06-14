package decision

const ITERATION_INSTRUCTIONS = `
To answer the question, you need to create a plan of action by considering which tools to use. Think step by step.
Then, you will use the selected tools to take the required actions.
`

// Decision represents a structured decision made by the agent.
type IterationDecision struct {
	DecisionStrategy DecisionStrategy
}

func (d *IterationDecision) Instructions() string {
	return ITERATION_INSTRUCTIONS
}
