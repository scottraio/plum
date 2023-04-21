package plum

import (
	embeddings "github.com/scottraio/plum/embeddings"
	llms "github.com/scottraio/plum/llms"
	store "github.com/scottraio/plum/vectorstores"
)

type Initialize struct {
	Embedding        string
	LLM              string
	VectorStore      string
	VectorStoreIndex string
}

func InitLLM(init Initialize) llms.LLM {
	var l llms.LLM

	switch init.LLM {
	case "openai":
		apiKey := GetDotEnvVariable("OPENAI_API_KEY")
		l = llms.InitOpenAI(apiKey)
		break
	default:
		apiKey := GetDotEnvVariable("OPENAI_API_KEY")
		l = llms.InitOpenAI(apiKey)
		break
	}

	return l
}

func InitEmbeddings(init Initialize) embeddings.Embedding {
	var e embeddings.Embedding

	switch init.Embedding {
	case "openai":
		apiKey := GetDotEnvVariable("OPENAI_API_KEY")
		e = embeddings.InitOpenAI(apiKey)
		break
	default:
		apiKey := GetDotEnvVariable("OPENAI_API_KEY")
		e = embeddings.InitOpenAI(apiKey)
		break
	}

	return e
}

func InitVectorStore(init Initialize) store.VectorStore {
	var v store.VectorStore

	switch init.VectorStore {
	case "pinecone":
		apiKey := GetDotEnvVariable("PINECONE_API_KEY")
		env := GetDotEnvVariable("PINECONE_ENV")
		projectId := GetDotEnvVariable("PINECONE_PROJECT_ID")
		v = store.InitPinecone(apiKey, env, projectId, init.VectorStoreIndex)
		break
	default:
		apiKey := GetDotEnvVariable("PINECONE_API_KEY")
		env := GetDotEnvVariable("PINECONE_ENV")
		projectId := GetDotEnvVariable("PINECONE_PROJECT_ID")
		v = store.InitPinecone(apiKey, env, projectId, init.VectorStoreIndex)
		break
	}

	return v
}
