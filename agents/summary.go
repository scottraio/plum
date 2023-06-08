package agents

// Summary represents a summary of multiple actions ran by an agent.
type SummaryPrompt struct {
	Context string
	Memory  string
	Summary string
	Input   string
}

// Summary Prompt
// -----------------

const SUMMARY_PROMPT = `
Background
----------
A Plum Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. 

Memory
---------
{{.Memory}}

Knowledge
---------
{{.Summary}}

Instructions
-------------
{{.Context}}

Use the knowledge and memory (if available) to answer the question. DO NOT make up answers.

Answers are helpful, friendly, informative, and confident. Use lists when possible. 

Additionally, evaluate the accuracy of the answer to the question. Determine a score (float64) from 0 to 1 that depicts the quality of the answer. Respond in JSON format.


Respond in the following JSON format ONLY!
------------------------------------------
{
	"question" : "{{.Input}}",
	"answer" : "the answer to the question",
	"score" : "the score of quality of the answer as float64",
	"reason": "the reason for the score with ways to improve"
}

Let's get started!
`

// func (a *Agent) summarize(toolOutputs []string) string {
// 	var scoreResult ScoreResult

// 	s := SummaryPrompt{
// 		Input:   a.Input,
// 		Context: a.Context,
// 		Summary: strings.Join(toolOutputs, "\n"),
// 		Memory:  a.Memory.Format()}

// 	prompt := llm.InjectObjectToPrompt(s, SUMMARY_PROMPT)
// 	// Log prompt to log file, do not show in stdout
// 	logger.PersistLog(prompt)

// 	answer := a.LLM.Run(prompt)

// 	err := json.Unmarshal([]byte(answer), &scoreResult)
// 	if err != nil {
// 		logger.Log("Error", "There was an error with the response from the LLM when trying to Unmarshal the JSON from the Summary Prompt, retrying: "+fmt.Sprintf("%v", err)+" original decision: "+scoreResult.Answer, "red")
// 		return a.summarize(toolOutputs)
// 	}

// 	logger.Log("Score", fmt.Sprintf("%f", scoreResult.Score), "purple")
// 	logger.Log("Reason", scoreResult.Reason, "purple")

// 	if scoreResult.Score > 0.9 {
// 		logger.Log("Answer", scoreResult.Answer, "green")
// 		return scoreResult.Answer
// 	} else {
// 		logger.Log("Retrying", "Score too low", "purple")

// 		return a.Answer(fmt.Sprintf("Question: %s, Previous Answer: %s \n Suggestions: %s", scoreResult.Question, scoreResult.Answer, scoreResult.Reason))
// 	}
// }
