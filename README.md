# Plum Framework
!!!Plum is under active development!!!

Plum is a framework designed for quickly building LLM-based (large language model) applications. 

## Features

1. Side load data into LLMs (ChatGPT) to enable conversations with your data.
2. A Rest API for interacting with your agents.
3. A CLI for querying agents or training models.

## Use Cases

1. Customer service chatbots.
2. Responding to emails.
3. Summarizing team meeting notes.
4. Next-gen AI applications.

## Getting Started

### Prerequisites

Currently, Plum uses OpenAI and Pinecone. There are plans to add more LLMs and vector databases in the future.

#### OpenAI

1. Obtain an OpenAI API key.
2. Set the `OPENAI_API_KEY` environment variable (`.env` is supported).

#### Pinecone

1. Signup for Pinecone and create an index.
2. Get the API key and set the `PINECONE_API_KEY` environment variable .
3. Set the `PINECONE_ENV` environment variable.
4. Set the `PINECONE_PROJECT_ID` environment variable.
Find your pod's URL: `[index]-[product_id].svc.[env].pinecone.io`.

To get Plum, use the following command: `go get github.com/scottraio/plum`.

## Anatomy of a Plum Application

The Plum framework includes agents, models, skills, and services that make up the building blocks of a Plum application.

### Plum Agents

Plum agents are the entry point for Plum applications, responsible for creating a plan of action and executing tools to achieve a goal.

```go
agent := plum.AsyncAgent(`
					You are the official [Company Name] customer service assistant. 
					Help and assist me with troubleshooting, guiding, and answer questions on [Company Name] products only.
				 `, CustomerServiceTools())

agent.Answer("Where's my order?")
```

### Plum Tools

Plum tools give agents the ability to perform specific tasks, typically returning results from models and skills.

```go
func CustomerServiceTools() []agents.Tool {	
	// Tools are the actions the agent can take.
	return []agents.Tool{
		{
			Name:        "OrderNumberLookup",
			Description: "Useful for finding tracking information, order status, and more",
			HowTo:       "Use the order status and tracking information to find the answer.",
			Func: func(query string) string {
				return plum.App.Skills["order_lookup"].Query(query)
			},
		},
		{
			Name:        "ProductManuals",
			Description: "Useful for finding general information about our products",
			HowTo:       "Use the information returned to find the answer.",
			Func: func(query string) string {
				lookup := plum.App.Models["manual"].Return(query)
				return lookup
			},
		},
	}
}
```

### Plum Models

Plum models are the "database" of Plum applications, taking source data, vectorizing it, and storing it in a vector store. Once trained, they can be queried to return results.

```go
&models.Model{
	Name: "Manual",

	// VectorStore is the vector store that you want to use for this model
	VectorStore: plum.App.VectorStore["knowledge"],

	// How to understand the data
	HowTo: `
		You are given excerpts from many markdown files. 
		Summarize the data and return a synopsis of the query. 
	`,

	// Return the data to be used for training, the vectors will be stored in the vector store
	Train: func(ctx context.Context) []store.Vector {
		// fetch source documents
		// split documents into chunks
		// build list of Vectors from chunks
		return []store.Vector{}
	},

	// Return result that you want to use in your prompt
	Return: func(query string) string {
		// Query the vector store
		result := manual.Find(query, map[string]string{}, map[string]interface{}{
			"TopK": float64(3),
		})

		// return the results, included with Model name and HowTo read the results
		return manual.Describe(query, result)
	},
},
```

### Plum Skills

Plum skills perform an action in real-time, such as performing a Google search, executing a script, or calling a REST API endpoint.

## Roadmap

1. Add more embedding options such as wego.
2. Implement continuous mode (allow agents to call other agents).
3. Add more LLM options.
4. Add more vector store options.
5. Add more vector store options.
