package plum

import (
	embeddings "github.com/scottraio/plum/embeddings"
	llms "github.com/scottraio/plum/llms"
	util "github.com/scottraio/plum/util"
	store "github.com/scottraio/plum/vectorstores"
)

type Initialize struct {
	Embedding         func(input string) []float32
	LLM               string
	VectorStoreConfig VectorStoreConfig
	Truths            []string
}

func InitLLM(init Initialize) llms.LLM {
	var l llms.LLM

	switch init.LLM {
	case "openai":
		apiKey := util.GetDotEnvVariable("OPENAI_API_KEY")
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
		apiKey := util.GetDotEnvVariable("OPENAI_API_KEY")
		util.FatalIfEmpty("OPENAI_API_KEY", apiKey)

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
		apiKey := util.GetDotEnvVariable("PINECONE_API_KEY")
		env := util.GetDotEnvVariable("PINECONE_ENV")
		projectId := util.GetDotEnvVariable("PINECONE_PROJECT_ID")
		plumEnv := util.GetDotEnvVariable("PLUM_ENV")

		v = store.InitPinecone(apiKey, env, projectId, init.VectorStoreConfig.Indexes, init.Embedding, plumEnv)
	default:
		init.VectorStoreConfig.Db = "pinecone"
		InitVectorStore(init)
	}

	return v
}

func Load() {

}
