package answer

import (
	"fmt"

	"encoding/json"

	"github.com/scottraio/plum/logger"
)

type ScoredAnswer struct {
	Score        float64 `json:"score"`
	Answer       string  `json:"answer"`
	Question     string  `json:"question"`
	Reason       string  `json:"reason"`
	Notes        string  `json:"notes"`
	AnswerMethod AnswerMethod
}

func (a *ScoredAnswer) SetQuestion(input string) {
	a.Question = input
}

func (a *ScoredAnswer) Format() string {
	format := fmt.Sprintf(`Respond with this JSON format: {
		"Question": "%s",
		"Answer": "the answer to the question",
		"Reason": "the reason to the answer",
		"Score": "the score of accuracy of the answer as float64. scores are between 0 and 1",
		"Notes": "Notes on improvements for future prompts",
	}`, a.Question)

	return format
}

func (a *ScoredAnswer) Instructions() string {
	instructions := `Instructions: 
	1. Carefully consider the outputs, rules and context. 
	
	2. Answer the question. 
	
	3. Score your answer from 0 to 1. The closer to 1, the better.

	4. Follow the rules as they only apply to you, not the user's question.
	`

	return instructions
}

func (a *ScoredAnswer) Output(jsonInput string) AnswerMethod {
	// Parse the JSON response to get the Decision object
	err := json.Unmarshal([]byte(jsonInput), a)
	if err != nil {
		logger.Log("Error", "There was an error with the response from the LLM, retrying: "+fmt.Sprintf("%v", err)+" original decision: "+a.Answer, "red")
		//a.Output(jsonInput)
	}

	return a
}

func (a *ScoredAnswer) Validate() bool {
	if a.Score >= 0.95 {
		return true
	} else {
		return false
	}
}

func (a *ScoredAnswer) Log() {
	logger.Log("Answer", a.Answer, "green")
	logger.Log("Score", fmt.Sprintf("%f", a.Score), "green")
	logger.Log("Reason", a.Reason, "green")
	logger.Log("Notes", a.Notes, "green")
}

func (a *ScoredAnswer) FinalAnswer() string {
	return a.Answer
}

func (a *ScoredAnswer) GetNotes() string {
	return a.Notes
}
