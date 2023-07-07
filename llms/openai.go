package llms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/scottraio/plum/logger"
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
func (ai *OpenAI) Run(memory memory.Memory) string {
	var messages []openai.ChatCompletionMessage

	for _, h := range memory.History {

		logger.Log("Prompt", h.Content, "white")
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	resp, err := ai._Client.CreateChatCompletion(
		ai._Context,
		openai.ChatCompletionRequest{
			Model:    "gpt-3.5-turbo-16k",
			Messages: messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "There was an error with the AI Engine. Please try again."
	}

	return resp.Choices[0].Message.Content
}

func (ai *OpenAI) Decide(memory memory.Memory) string {
	var messages []openai.ChatCompletionMessage

	for _, h := range memory.History {
		logger.Log("Prompt", h.Content, "white")
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	functions := ai.DecisionFunction()

	funcJSON, _ := json.Marshal(functions)
	logger.Log("Functions", string(funcJSON), "white")

	req := openai.ChatCompletionRequest{
		Model:     "gpt-3.5-turbo-16k",
		Messages:  messages,
		Functions: functions,
	}

	resp, err := ai._Client.CreateChatCompletion(ai._Context, req)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		fmt.Printf("ChatCompletion error resp: %v\n", resp)
		return "There was an error with the AI Engine. Please try again."
	}

	logger.Log("Decision", fmt.Sprintf("%v", resp.Choices[0].Message), "white")
	return resp.Choices[0].Message.FunctionCall.Arguments
}

func (ai *OpenAI) Answer(memory memory.Memory) string {
	var messages []openai.ChatCompletionMessage

	for _, h := range memory.History {
		logger.Log("Prompt", h.Content, "white")
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	functions := ai.AnswerFunction()

	funcJSON, _ := json.Marshal(functions)
	logger.Log("Functions", string(funcJSON), "white")

	req := openai.ChatCompletionRequest{
		Model:     "gpt-3.5-turbo-16k",
		Messages:  messages,
		Functions: functions,
	}

	resp, err := ai._Client.CreateChatCompletion(ai._Context, req)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		fmt.Printf("ChatCompletion error resp: %v\n", resp)
		return "There was an error with the AI Engine. Please try again."
	}

	if len(resp.Choices) > 0 {
		message := resp.Choices[0].Message
		if message.FunctionCall != nil {
			arguments := message.FunctionCall.Arguments
			// Use the 'arguments' variable for further processing
			return arguments
		} else {
			// Handle the case when 'resp.Choices[0].Message.FunctionCall' is nil
			return "There was an error with the AI Engine. Couldn't find Message."
		}
	} else {
		// Handle the case when 'resp.Choices' is empty
		return "There was an error with the AI Engine. Choices are empty."
	}

}

// {
// 	"Question": "{{.Input}}",
// 	"Thought": "the thought about what action(s) and input(s) are required to answer the question.",
// 	"Actions": [{
// 		"Tool": "the tool name to use",
// 		"Thought": "the thought about what the input to the tool should be",
// 		"Notes": "Notes on improvements for future prompts",
// 		"Input": "the input to the tool"
// 	}]
// }`

func (ai *OpenAI) DecisionFunction() []*openai.FunctionDefine {
	paramsSchema := json.RawMessage(`{
		"type": "object",
		"properties": {
			"Question": {
				"type": "string"
			},
			"Thought": {
				"type": "string"
			},
			"Actions": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"Tool": {
							"type": "string"
						},
						"Thought": {
							"type": "string"
						},
						"Notes": {
							"type": "string"
						},
						"Input": {
							"type": "string"
						}
					},
					"required": ["Tool", "Thought", "Notes", "Input"]
				}
			}
		},
		"required": ["Question", "Thought", "Actions"]
	}
`)

	function := &openai.FunctionDefine{
		Name:        "decision_response_function",
		Description: "The actions to take to answer the question",
		Parameters:  paramsSchema,
	}

	return []*openai.FunctionDefine{function}
}

func (ai *OpenAI) AnswerFunction() []*openai.FunctionDefine {
	paramsSchema := json.RawMessage(`{
			"type": "object",
			"properties": {
				"Question": {
					"type": "string"
				},
				"Answer": {
					"type": "string"
				},
				"Reason": {
					"type": "string"
				},
				"Score": {
					"type": "number",
					"minimum": 0,
					"maximum": 1
				},
				"Notes": {
					"type": "string"
				}
			},
			"required": ["Question", "Answer", "Reason", "Score", "Notes"]
		}
		
`)

	function := &openai.FunctionDefine{
		Name:        "answer_response_function",
		Description: "The answer to the question",
		Parameters:  paramsSchema,
	}

	return []*openai.FunctionDefine{function}
}
