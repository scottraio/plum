package answer

import (
	"encoding/json"
	"fmt"

	"github.com/scottraio/plum/logger"
)

type SummarizedAnswer struct {
	AnswerMethod AnswerMethod

	Answer   string `json:"answer"`
	Question string `json:"question"`
	Reason   string `json:"reason"`
	Notes    string `json:"notes"`
}

func (a *SummarizedAnswer) SetQuestion(input string) {
	a.Question = input
}

func (a *SummarizedAnswer) Format() string {
	format := fmt.Sprintf(`Respond with this JSON format: {
		"Question": "%s",
		"Answer": "the answer to the question",
		"Reason": "the reason to the answer",
		"Notes": "Notes on improvements for future prompts",
	}`, a.Question)

	return format
}

func (a *SummarizedAnswer) Output(jsonInput string) AnswerMethod {
	// Parse the JSON response to get the Decision object
	err := json.Unmarshal([]byte(jsonInput), a)
	if err != nil {
		logger.Log("Error", "There was an error with the response from the LLM, retrying: "+fmt.Sprintf("%v", err)+" original decision: "+a.Answer, "red")
		//a.Output(jsonInput)
	}

	return a
}

func (a *SummarizedAnswer) Instructions() string {
	instructions := `Instructions:
	
	1. Carefully consider the outputs, rules and context. 
	
	2. Answer the question by giving a detailed summary of the outputs. 

	4. Follow the rules as they only apply to you, not the user's question.`
	return instructions
}

func (a *SummarizedAnswer) Validate() bool {
	return true
}

func (a *SummarizedAnswer) Log() {
	logger.Log("Answer", a.Answer, "green")
	logger.Log("Reason", a.Reason, "cyan")
}

func (a *SummarizedAnswer) FinalAnswer() string {
	return a.Answer
}

func (a *SummarizedAnswer) GetNotes() string {
	return a.Notes
}
