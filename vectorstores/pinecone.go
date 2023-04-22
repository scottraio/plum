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
	embedFunc func(text string) []float32) map[string]VectorStore {
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
func (p *Pinecone) Query(input string, fields map[string]string) string {
	var filtered structpb.Struct
	// get the embeddings
	embeddings := p.EmbedFunc(input)

	if fields != nil {
		filtered = structpb.Struct{
			Fields: p.WithFields(fields),
		}
	}

	// client'
	queryResult, queryErr := p._Client.Query(p._Context, &pinecone_grpc.QueryRequest{
		Queries: []*pinecone_grpc.QueryVector{
			{Values: embeddings},
		},
		TopK:            3,
		IncludeValues:   false,
		IncludeMetadata: true,
		Namespace:       p.Namespace,
		Filter:          &filtered,
	})

	if queryErr != nil {
		log.Fatalf("query error: %v", queryErr)
	}

	// return the first result
	// var resultString string&{"": ""}

	// for _, match := range queryResult.Results[0].Matches {
	// 	resultString += match.Metadata.Fields["text"].Value.StringValue
	// }

	var resultString string
	for _, match := range queryResult.Results[0].Matches {
		resultString += p.capText(match.Metadata.Fields["text"].GetStringValue())
	}

	return resultString

}

// Upsert upserts a document into the Pinecone index.
func (p *Pinecone) Upsert(namespace string, text string, fields map[string]string) error {
	var vects []*pinecone_grpc.Vector
	var meta structpb.Struct

	// get the embeddings
	embeddings := p.EmbedFunc(text)

	fields["text"] = text

	meta = structpb.Struct{
		Fields: p.WithFields(fields),
	}

	vects = append(vects, &pinecone_grpc.Vector{
		Id:       uuid.New().String(),
		Values:   embeddings,
		Metadata: &meta,
	})

	_, upsertErr := p._Client.Upsert(p._Context, &pinecone_grpc.UpsertRequest{
		Vectors:   vects,
		Namespace: namespace,
	})

	// return the first result
	return upsertErr
}

// WithNamespace sets the namespace for the Pinecone.
func (p *Pinecone) WithNamespace(namespace string) VectorStore {
	p.Namespace = namespace
	return p.VectorStore
}

// capText caps the text at 1000 characters.
func (p *Pinecone) capText(text string) string {
	if len(text) > 1000 {
		return text[:1000]
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
