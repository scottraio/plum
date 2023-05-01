package vectorstores

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

type PineconeConfig struct {
	PineconeKey       string
	PineconeEnv       string
	PineconeProjectId string
}

type Pinecone struct {
	VectorStore
	Config    PineconeConfig
	EmbedFunc func(text string) []float32
	IndexName string
	Namespace string
	_Client   pinecone_grpc.VectorServiceClient
	_Context  context.Context
}

// InitPinecone initializes a Pinecone client.
func InitPinecone(apiKey string,
	env string,
	projectId string,
	indexes []string,
	embedFunc func(text string) []float32,
	namespace string) map[string]VectorStore {
	indexMap := make(map[string]VectorStore)

	for _, index := range indexes {
		config := &PineconeConfig{
			PineconeKey:       apiKey,
			PineconeEnv:       env,
			PineconeProjectId: projectId,
		}

		pinecone := Pinecone{
			Config:    *config,
			IndexName: index,
			EmbedFunc: embedFunc,
			Namespace: namespace,
		}

		indexMap[index] = pinecone.NewClient()
	}

	return indexMap
}

// Client returns a Pinecone client.
func (p *Pinecone) NewClient() VectorStore {
	config := &tls.Config{}

	ctx := context.Background()

	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", p.Config.PineconeKey)
	target := fmt.Sprintf("%s-%s.svc.%s.pinecone.io:443", p.IndexName, p.Config.PineconeProjectId, p.Config.PineconeEnv)

	log.Printf("connecting to %v", target)
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithAuthority(target),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	//defer conn.Close()

	client := pinecone_grpc.NewVectorServiceClient(conn)

	p._Client = client
	p._Context = ctx

	return p
}

// Query queries the Pinecone index.
func (p *Pinecone) Query(input string, fields map[string]string, options map[string]interface{}) string {
	var filtered structpb.Struct
	// get the embeddings
	embeddings := p.EmbedFunc(input)

	if fields != nil {
		filtered = structpb.Struct{
			Fields: p.WithFields(fields),
		}
	}

	opts := p.queryOptions(options)

	// client'
	queryResult, queryErr := p._Client.Query(p._Context, &pinecone_grpc.QueryRequest{
		Queries: []*pinecone_grpc.QueryVector{
			{Values: embeddings},
		},
		TopK:            opts["TopK"].(uint32),
		IncludeValues:   false,
		IncludeMetadata: true,
		Namespace:       opts["Namespace"].(string),
		Filter:          &filtered,
	})

	if queryErr != nil {
		log.Fatalf("query error: %v", queryErr)
	}

	var resultString string
	for _, match := range queryResult.Results[0].Matches {
		resultString += p.capText(match.Metadata.Fields["text"].GetStringValue())
	}

	return resultString

}

func (p *Pinecone) queryOptions(options map[string]interface{}) map[string]interface{} {
	defaults := map[string]interface{}{
		"TopK":      uint32(1),
		"Namespace": p.Namespace,
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

// Upsert upserts a document into the Pinecone index.
func (p *Pinecone) Upsert(text string, fields map[string]string, options map[string]interface{}) error {
	var vects []*pinecone_grpc.Vector
	var meta structpb.Struct

	// Set Fields and Options
	fields["text"] = text

	meta = structpb.Struct{
		Fields: p.WithFields(fields),
	}

	opts := p.queryOptions(options)

	// get the embeddings
	embeddings := p.EmbedFunc(text)

	vects = append(vects, &pinecone_grpc.Vector{
		Id:       uuid.New().String(),
		Values:   embeddings,
		Metadata: &meta,
	})

	_, upsertErr := p._Client.Upsert(p._Context, &pinecone_grpc.UpsertRequest{
		Vectors:   vects,
		Namespace: opts["Namespace"].(string),
	})

	// return the first result
	return upsertErr
}

// purge
func (p *Pinecone) Purge(namespace string) error {
	deleteResult, deleteErr := p._Client.Delete(p._Context, &pinecone_grpc.DeleteRequest{
		DeleteAll: true,
		Namespace: namespace,
	})

	if deleteErr != nil {
		log.Fatalf("delete error: %v", deleteErr)
	} else {
		log.Printf("delete result: %v", deleteResult)
	}

	return deleteErr
}

// WithNamespace sets the namespace for the Pinecone.
func (p *Pinecone) WithNamespace(namespace string) VectorStore {
	p.Namespace = namespace
	return p.VectorStore
}

// capText caps the text at 1000 characters.
func (p *Pinecone) capText(text string) string {
	if len(text) > 3000 {
		return text[:3000]
	}

	return text
}

func (p *Pinecone) WithFields(fields map[string]string) map[string]*structpb.Value {
	filtered := make(map[string]*structpb.Value)

	for key, val := range fields {
		filtered[key] = &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: val},
		}
	}

	return filtered
}
