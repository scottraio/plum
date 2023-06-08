package memory

import (
	"fmt"

	"github.com/scottraio/plum/logger"
)

type Memory struct {
	History []ChatHistory
}
type ChatHistory struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func LoadMemory(history []ChatHistory) *Memory {
	m := &Memory{}
	m.History = append(m.History, history...)
	return m
}

func (m *Memory) Add(content string, role string, color string) {
	logger.Log(role, content, color)
	m.History = append(m.History, ChatHistory{content, role})
}

func (c *ChatHistory) Memory() *Memory {
	return LoadMemory([]ChatHistory{*c})
}

func (m *Memory) Format() string {
	var output string
	output = "\n"
	for i := range m.History {
		output += fmt.Sprintf("%d. role => %s content => %s\n", i, m.History[i].Role, m.History[i].Content)
	}
	return output
}
