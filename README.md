# Plum
Plum is a framework for rapidly building LLM-based applications. 
Inspired by [Langchain](https://github.com/hwchase17/langchain).

## Features
1. Side load data into a LLMs (ChatGPT) to have a conversation with your data. 
2. Rest API for interacting with your agents
3. Powerful CLI for querying agents or training models
4. Written in go. 

## Use Cases
1. Customer service chat bots
2. Respond to emails
3. Summarization of team meeting notes
4. Next-gen AI applications

## Getting Started

### Prerequisites
Currently, Plum supports OpenAI (Embeddings and LLM) and Pinecone (Vector Store) only. 

#### OpenAI
1. Obtain an OpenAI API key
2. Set the OPENAI_API_KEY environment variable (.env is supported)

#### Pinecone
1. Signup for Pinecone and create an index.
2. Get the API key and set the PINECONE_API_KEY environment variable (.env is supported)
3. Set the PINECONE_ENV environment variable (.env is supported)
4. Set the PINECONE_PROJECT_ID environment variable (.env is supported)

Find your pod's url. [index]-[product_id].svc.[env].pinecone.io

#### Get Plum
```
go get github.com/scottraio/plum

```

- [ ] TODO: `plum new [insert app name]`

## Anatomy of a Plum Application
* Agents have many tools. 
* Tools use one or more Model or Skill. 
* Models do two things:
  * Take source data, transform text into vectors and store the data in a vector store (training). 
  * Models also retrieve data from the vector store.
* Skills do a wide variety of things, like fetch google results, execute scripts, or send emails.

```
[input] => [Plum AI App] => [output]
                 |
                 |
                 |
                 - Search your own documents with powerful vector search
                 - Call REST API Endpoints
                 - Get latest news articles
                 - Execute Scripts
```


### Plum Application Structure
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
Plum models represent the "database" of a Plum application. 
They take source data, vectorize it, and store it in a vector store.
Once trained, they can be queried to return results.

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
Plum Skills are the opposite of a model. They perform an action in real time.
This can be performing a google search, executing a script, or calling a REST API endpoint.

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

## Roadmap
1. [ ] Add more Embedding Options e.g. [wego](https://github.com/ynqa/wego)
2. [ ] Continuous Mode (allow agents to call other agents)
3. [ ] Add more LLM options
4. [ ] Add more vector store options 
5. [ ] Add more vector store options