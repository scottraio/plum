package decision

const MULTIPLE_SELECTION_INSTRUCTIONS = `
Think step by step and create a plan of action by considering which tools to use. 

Once you have a plan, you will use the selected tools to take the required actions. Select as many tools as needed.

`

// Decision represents a structured decision made by the agent.
type MultipleSelectionDecision struct {
	DecisionMethod DecisionMethod
}

func (d *MultipleSelectionDecision) Instructions() string {
	return MULTIPLE_SELECTION_INSTRUCTIONS
}
