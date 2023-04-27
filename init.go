package plum

import (
	embeddings "github.com/scottraio/plum/embeddings"
	llms "github.com/scottraio/plum/llms"
	store "github.com/scottraio/plum/vectorstores"
)

type Initialize struct {
	Embedding         func(input string) []float32
	LLM               string
	VectorStoreConfig VectorStoreConfig
}

func InitLLM(init Initialize) llms.LLM {
	var l llms.LLM

	switch init.LLM {
	case "openai":
		apiKey := GetDotEnvVariable("OPENAI_API_KEY")
		l = llms.InitOpenAI(apiKey)
	default:
		init.LLM = "openai"
		InitLLM(init)
	}

	return l
}

func InitEmbeddings(embed string) func(input string) []float32 {
	var e func(input string) []float32

	switch embed {
	case "openai":
		apiKey := GetDotEnvVariable("OPENAI_API_KEY")
		embed := embeddings.InitOpenAI(apiKey)
		e = embed.EmbedText
	default:
		InitEmbeddings("openai")
	}

	return e
}

func InitVectorStore(init Initialize) map[string]store.VectorStore {
	var v map[string]store.VectorStore

	switch init.VectorStoreConfig.Db {
	case "pinecone":
		apiKey := GetDotEnvVariable("PINECONE_API_KEY")
		env := GetDotEnvVariable("PINECONE_ENV")
		projectId := GetDotEnvVariable("PINECONE_PROJECT_ID")
		plumEnv := GetDotEnvVariable("PLUM_ENV")

		v = store.InitPinecone(apiKey, env, projectId, init.VectorStoreConfig.Indexes, init.Embedding, plumEnv)
	default:
		init.VectorStoreConfig.Db = "pinecone"
		InitVectorStore(init)
	}

	return v
}
