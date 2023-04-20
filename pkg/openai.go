package framework

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	LLM
	C   *openai.Client
	Ctx context.Context
}

func NewOpenAI(token string) *OpenAI {
	c := openai.NewClient(token)
	ctx := context.Background()
	return &OpenAI{C: c, Ctx: ctx}
}

func (ai *OpenAI) Run(prompt string) string {
	resp, err := ai.C.CreateChatCompletion(
		ai.Ctx,
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
		return "Error"
	}

	return resp.Choices[0].Message.Content
}
