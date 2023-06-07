package agents

const PARALLEL_INSTRUCTIONS = `
To answer the question, you will need to form an input for each tool. 
Then, you will use the selected tools to take the required actions. 
Must use every tool listed.
`

// Decision represents a structured decision made by the agent.
type ParallelDecision struct {
	DecisionMethod DecisionMethod
}

func (d *ParallelDecision) Instructions() string {
	return PARALLEL_INSTRUCTIONS
}
