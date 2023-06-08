package memory

import (
	"fmt"
)

type Memory struct {
	History []ChatHistory
}
type ChatHistory struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func LoadMemory(history []ChatHistory) Memory {
	m := Memory{}
	m.History = append(m.History, history...)
	return m
}

func (m *Memory) Add(content string, role string) {
	//logger.Log(role, content, m.LogColor(role))
	m.History = append(m.History, ChatHistory{content, m.GetRole(role)})
}

func (m *Memory) GetRole(role string) string {
	switch role {
	case "tool":
		return "system"
	case "output_format":
		return "system"
	case "background":
		return "system"
	case "context":
		return "system"
	case "answer":
		return "assistant"
	case "decision":
		return "assistant"
	case "prompt":
		return "assistant"
	default:
		return role
	}
}

func (m *Memory) LogColor(role string) string {
	switch role {
	case "background":
		return "purple"
	case "answer":
		return "green"
	case "context":
		return "purple"
	case "prompt":
		return "purple"
	case "tool":
		return "yellow"
	case "assistant":
		return "purple"
	case "decision":
		return "purple"
	case "system":
		return "cyan"
	case "user":
		return "yellow"
	default:
		return "white"
	}
}

func (c *ChatHistory) Memory() Memory {
	return LoadMemory([]ChatHistory{*c})
}

func (m Memory) Format() string {
	var output string
	output = "\n"
	for i := range m.History {
		output += fmt.Sprintf("%d. role => %s content => %s\n", i, m.History[i].Role, m.History[i].Content)
	}
	return output
}
