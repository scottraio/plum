package llms

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/scottraio/plum/memory"
)

type OpenAIConfig struct {
	OpenAIToken string
}

type OpenAI struct {
	LLM
	Config   OpenAIConfig
	_Client  *openai.Client
	_Context context.Context
}

func InitOpenAI(apiKey string) LLM {
	config := &OpenAIConfig{
		OpenAIToken: apiKey,
	}

	ai := OpenAI{
		Config:   *config,
		_Client:  openai.NewClient(config.OpenAIToken),
		_Context: context.Background(),
	}

	return &ai
}

// Client returns a OpenAI client.
func (ai *OpenAI) Client() LLM {
	return ai
}

// Run returns a response from OpenAI.
func (ai *OpenAI) Run(memory *memory.Memory) string {
	var messages []openai.ChatCompletionMessage

	for _, h := range memory.History {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	resp, err := ai._Client.CreateChatCompletion(
		ai._Context,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "There was an error with the AI Engine. Please try again."
	}

	return resp.Choices[0].Message.Content
}
