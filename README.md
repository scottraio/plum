# Plum

Plum is a framework for rapidly building LLM-based applications. 

## Features
1. Side load data into a LLMs (ChatGPT) to have a conversation with your data. 
2. Rest API for interacting with your agents
3. Powerful CLI for querying agents or training models
4. Written in go. 

## Get Plum

```
go get github.com/scottraio/plum
```

## Anatomy of a Plum Application

```
/
agents/
	- my_agent.go
models/
	- recipe.go
skills/
	- epicurious.go
services/
	- mysql.go
boot.go
main.go
server.go
```

### Plum Agents
Plum agents are the main entry point for Plum applications. They are responsible for doing near-shot prompting to 
return a plan of action that executes tools to achieve a goal, mostly QA. 

```
A Plum Agent is a powerful language model that can assist with a wide range of tasks, including 
answering questions and providing in-depth explanations and discussions on various topics. It can 
process and understand large amounts of text, generate human-like responses, and provide valuable insights 
and information. As a JSON API, a Plum Agent determines the necessary actions to take based on the input 
received from the user. A Plum Agent understands csv, markdown, json, html and plain text.

Respond in the following JSON format:
-------------------------------------
{
	"Question": "{{.Input}}",
	"Thought": "Think about what action and input are required to answer the question.",
	"Actions": [{
		"Tool": "the tool name to use",
		"Input": "the input to the tool",
	}]
}
```

Plum agents are designed to be conversational. You give them a question, they create plan of action and execute tools, then return the result. 

#### Plum Tools
A plum tool gives the agent the ability to perform a specific task. Plum tools are the building blocks of a Plum agent.
They typically return results from models and skills, but can accomplish anything you write. 

```go
Tool{
	Name:        "Manuals",
	Description: "Useful for finding high level information, troubleshooting, and operating procedures.",
	HowTo:       `Example Input: "dp1234 description, troubleshooting, operating procedures, cleaning, maintenance"`,
	Func: func(query string) string {
		lookup := manual.Return(query)
		return lookup
	},
},
```

### Plum Models
A plum model represents the "database" of a Plum application. 
It is a collection of data, that's been vectorized with an embedding func, that can be queried by plum tools. 

```go
Model{
	Name: "Recipe",

	// VectorStore is the vector store that you want to use for this model
	VectorStore: plum.App.VectorStore["knowledge"],

	// How to understand the data
	HowTo: `
		You are given excerpts from many markdown files. 
		Summarize the data and return a synopsis of the query. 
	`,

	// Train is a function that returns the data to be used for training
	Train: func(ctx context.Context) []store.Vector {
		// ETL the data
		// The results of the []store.Vector will be stored in the VectorStore defined above.
		// 1 store.Vector is 1 record in the vectorstore
	},

	// Return is a function that returns the result that you want to use in your prompt
	Return: func(query string) string {
		// Query the vectorstore
		// Return the result
	},
	},
},
```

### Plum Skills

A Plum Skill represents any action or task not found in the vector store. This could be a google search, 
REST API call, executing a script, or sending an email. 

```go
Skill{
	HowTo: `
		You are given a JSON response from a REST API.
		Always just return a sentence from the JSON attributes and values.
	`,

	Return: func(query string) string {
		return http.Lookup(query)
	},
}
```