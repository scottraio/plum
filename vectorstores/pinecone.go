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
)

type PineconeConfig struct {
	PineconeKey       string
	PineconeEnv       string
	PineconeProjectId string
}

type Pinecone struct {
	VectorStore
	Config    PineconeConfig
	IndexName string
	Namespace string
	_Client   pinecone_grpc.VectorServiceClient
	_Context  context.Context
}

// InitPinecone initializes a Pinecone client.
func InitPinecone(apiKey string, env string, productId string, index string) VectorStore {
	config := &PineconeConfig{
		PineconeKey:       apiKey,
		PineconeEnv:       env,
		PineconeProjectId: productId,
	}

	pinecone := Pinecone{
		Config:    *config,
		IndexName: index,
	}

	pinecone.NewClient()

	return &pinecone
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
func (p *Pinecone) Query(embeddings []float32) string {
	// client
	queryResult, queryErr := p._Client.Query(p._Context, &pinecone_grpc.QueryRequest{
		Queries: []*pinecone_grpc.QueryVector{
			{Values: embeddings},
		},
		TopK:          1000,
		IncludeValues: true,
		Namespace:     p.Namespace,
	})

	if queryErr != nil {
		log.Fatalf("query error: %v", queryErr)
	} else {
		log.Printf("query result: %v", queryResult)
	}

	// return the first result
	return "implement me"
}

// Upsert upserts a document into the Pinecone index.
func (p *Pinecone) Upsert(namespace string, embeddings []float32) string {
	log.Print("upserting data...")

	upsertResult, upsertErr := p._Client.Upsert(p._Context, &pinecone_grpc.UpsertRequest{
		Vectors: []*pinecone_grpc.Vector{
			{
				Id:     uuid.New().String(),
				Values: embeddings,
			},
		},
		Namespace: namespace,
	})

	if upsertErr != nil {
		log.Fatalf("upsert error: %v", upsertErr)
	} else {
		log.Printf("upsert result: %v", upsertResult)
	}

	// return the first result
	return "success"
}

func (p *Pinecone) Index(index string) VectorStore {
	p.IndexName = index
	return p.NewClient()
}

// capText caps the text at 1000 characters.
func (p *Pinecone) capText(text string) string {
	if len(text) > 1000 {
		return text[:1000]
	}

	return text
}
