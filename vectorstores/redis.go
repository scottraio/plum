package vectorstores

import (
	redisearch "github.com/RediSearch/redisearch-go/redisearch"

	"google.golang.org/protobuf/types/known/structpb"
)

type RedisConfig struct {
	RedisAddress  string
	RedisPassword string
}

type RedisVectorStore struct {
	VectorStore
	Config    RedisConfig
	EmbedFunc func(text string) []float32
	IndexName string
	Namespace string
	_Client   *redisearch.Client
}

func InitRedisVectorStore(address string,
	password string,
	indexes []string,
	embedFunc func(text string) []float32,
	namespace string) map[string]VectorStore {
	indexMap := make(map[string]VectorStore)

	for _, index := range indexes {
		config := &RedisConfig{
			RedisAddress:  address,
			RedisPassword: password,
		}

		store := RedisVectorStore{
			Config:    *config,
			IndexName: index,
			EmbedFunc: embedFunc,
			Namespace: namespace,
		}

		indexMap[index] = store.NewClient()
	}

	return indexMap
}

func (r *RedisVectorStore) NewClient() VectorStore {
	client := redisearch.NewClient(r.Config.RedisAddress, r.IndexName)

	// Authenticate if a password is provided
	// if r.Config.RedisPassword != "" {
	// 	err := client.Auth(r.Config.RedisPassword)
	// 	if err != nil {
	// 		log.Fatalf("failed to authenticate: %v", err)
	// 	}
	// }

	r._Client = client

	return r
}

func (r *RedisVectorStore) Query(input string, fields map[string]string, options map[string]interface{}) string {
	// Use the RedisAI client to perform the query
	return ""
}

func (r *RedisVectorStore) queryOptions(options map[string]interface{}) map[string]interface{} {
	defaults := map[string]interface{}{
		"TopK":      uint32(1),
		"Namespace": r.Namespace,
	}

	for key, value := range options {
		if key == "TopK" {
			defaults["TopK"] = uint32(value.(float64))
		} else {
			defaults[key] = value
		}
	}

	return defaults
}

func (r *RedisVectorStore) Upsert(text string, fields map[string]string, options map[string]interface{}) error {
	// Use the RedisAI client to upsert the document
	return nil
}

func (r *RedisVectorStore) Purge(namespace string) error {
	// Use the RedisAI client to purge the data

	return nil
}

func (r *RedisVectorStore) WithNamespace(namespace string) VectorStore {
	r.Namespace = namespace
	return r.VectorStore
}

func (r *RedisVectorStore) capText(text string) string {
	if len(text) > 3000 {
		return text[:3000]
	}

	return text
}

func (r *RedisVectorStore) WithFields(fields map[string]string) map[string]*structpb.Value {
	filtered := make(map[string]*structpb.Value)

	for key, val := range fields {
		filtered[key] = &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: val},
		}
	}

	return filtered
}
