package decision

const SINGLE_INSTRUCTIONS = `
To answer the question, you need to think about the question and the tools you have available. Select 1 tool only.
`

// Decision represents a structured decision made by the agent.
type SingleDecision struct {
	DecisionMethod DecisionMethod
}

func (d *SingleDecision) Instructions() string {
	return SINGLE_INSTRUCTIONS
}
