package agents

import (
	llm "github.com/scottraio/plum/llms"
	memory "github.com/scottraio/plum/memory"
)

type Engine interface {
	Answer(question string) string
	Remember(memory *memory.Memory) Engine
}

type Agent struct {
	Input   string
	Context string

	LLM    llm.LLM
	Memory *memory.Memory

	Tools     []Tool
	ToolNames []string
}
