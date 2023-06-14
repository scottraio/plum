package decision

const MULTIPLE_SELECTION_INSTRUCTIONS = `
1. Think step by step and create a plan of action by selecting which tools to use. 

2. Once you have a plan, you will use the selected tools to take the required actions. Select as many tools as needed.

3. Follow the rules as they only apply to you, not the user's question.

4. Carefully consider the input to the tools.
`

// Decision represents a structured decision made by the agent.
type MultipleSelectionDecision struct {
	DecisionStrategy DecisionStrategy
}

func (d *MultipleSelectionDecision) Instructions() string {
	return MULTIPLE_SELECTION_INSTRUCTIONS
}
