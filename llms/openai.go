package llms

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
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
func (ai *OpenAI) Run(prompt string) string {
	resp, err := ai._Client.CreateChatCompletion(
		ai._Context,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "There was an error with the AI Engine. Please try again."
	}

	return resp.Choices[0].Message.Content
}
