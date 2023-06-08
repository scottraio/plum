package decision

const SEQUENTIAL_INSTRUCTIONS = `
To answer the question, you need to create a plan of action by considering which tools to use. Think step by step.
Then, you will use the selected tools to take the required actions. 
`

// Decision represents a structured decision made by the agent.
type SequentialDecision struct {
	DecisionMethod DecisionMethod
}

func (d *SequentialDecision) Instructions() string {
	return SEQUENTIAL_INSTRUCTIONS
}

// // Decide makes a decision based on the agent's input and memory.
// func (a *Decision) StepsToString() string {
// 	steps := ""
// 	for i, step := range a.Steps {
// 		steps += fmt.Sprintf("Step %d: %s", i, step.Description)
// 	}
// 	return steps
// }

// func (action *Action) ActionToString() string {
// 	json, _ := json.Marshal(action)
// 	return string(json)
// }

// func (action *Action) Branch(input Input, context string) string {
// 	prompt := `
// 	Steps:
// 	` + "```" + input.Plans + "```" + `

// 	Current Step:
// 	` + "```" + input.CurrentStep + "```" + `

// 	Current Tool:
// 	` + "```" + input.ToolName + "```" + `

// 	Original Input:
// 	` + "```" + input.Text + "```" + `

// 	You have been given steps, this is the main plan. You are on the Current Step,
// 	using the Current Tool.

// 	` + context

// 	return input.LLM.Run(prompt)
// }
