package memory

import "fmt"

type Memory struct {
	History []ChatHistory
}

type ChatHistory struct {
	Query  string `json:"query"`
	Answer string `json:"answer"`
}

func LoadMemory(history []ChatHistory) *Memory {
	m := &Memory{}
	for _, chatHistory := range history {
		m.History = append(m.History, chatHistory)
	}
	return m
}

func (c *ChatHistory) Memory() *Memory {
	return LoadMemory([]ChatHistory{*c})
}

func (m *Memory) Format() string {
	var output string
	output = "\n"
	for i := range m.History {
		output += fmt.Sprintf("%d. in => %s out => %s\n", i, m.History[i].Query, m.History[i].Answer)
	}
	return output
}
