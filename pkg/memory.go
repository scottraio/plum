package framework

import "fmt"

type Memory struct {
	History []ChatHistory
}

type ChatHistory struct {
	Query  string `json:"query"`
	Answer string `json:"answer"`
}

func LoadMemory(memory ChatHistory) Memory {
	m := Memory{}
	m.History = append(m.History, memory)
	return m
}

func (m *Memory) Format() string {
	var output string
	output = "\n"
	for i := range m.History {
		output += fmt.Sprintf("%d. [in] %s [out] %s\n", i, m.History[i].Query, m.History[i].Answer)
	}
	return output
}
