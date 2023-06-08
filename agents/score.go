package agents

const SCORE_PROMPT = `
Background
----------
A Plum Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. As a JSON API, a Plum Agent determines the necessary actions to take based on the input 
received from the user. A Plum Agent understands csv, markdown, json, html and plain text.

Follow these instructions:
-------------------------------------------------
1. Evaluate the accuracy of the answer to the question.
2. Determine a score from 0 to 1 that depicts the accuracy of the answer.
3. Respond in JSON with the score.


Respond in the following JSON format ONLY!:
-------------------------------------
{
	"question" : "{{.Input}}",
	"answer" : "{{.Answer}}",
	"score" : "the score of accuracy of the answer as float64",
}

Let's get started!
`

type ScorePrompt struct {
	Input  string
	Answer string
}

type ScoreResult struct {
	Score    float64 `json:"score"`
	Answer   string  `json:"answer"`
	Question string  `json:"question"`
	Reason   string  `json:"reason"`
}

// func (agent *Agent) ScoreAnswer(question string, answer string) float64 {
// 	var scoreResult ScoreResult
// 	logger.Log("Agent", "Scoring...", "cyan")

// 	scorePrompt := &ScorePrompt{
// 		Input:  question,
// 		Answer: answer}

// 	prompt := llm.InjectObjectToPrompt(scorePrompt, SCORE_PROMPT)
// 	decision := agent.LLM.Run(prompt)

// 	err := json.Unmarshal([]byte(decision), &scoreResult)
// 	if err != nil {
// 		logger.Log("Error", "There was an error with the response from the LLM, retrying: "+fmt.Sprintf("%v", err)+" original decision: "+decision, "red")
// 		agent.ScoreAnswer(question, answer)
// 	}

// 	// Verbose logging
// 	logger.Log("Score", fmt.Sprintf("%f", scoreResult.Score), "cyan")

// 	return scoreResult.Score
// }
