package plum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Pinecone struct {
	ApiKey      string
	IndexName   string
	ProjectName string
	PineconeEnv string
	Version     string
}

type PineconeResponse struct {
	Results []interface{} `json:"results"`
	Matches []struct {
		ID       string                 `json:"id"`
		Score    float64                `json:"score"`
		Values   []interface{}          `json:"values"`
		Metadata map[string]interface{} `json:"metadata"`
	} `json:"matches"`
}

type QueryPayload struct {
	IncludeValues   bool      `json:"includeValues"`
	IncludeMetadata bool      `json:"includeMetadata"`
	Vector          []float32 `json:"vector"`
	TopK            int       `json:"topK"`
	Namespace       string    `json:"namespace"`
}

// Query queries the Pinecone index.
func (p *Pinecone) Query(version string, input string) string {
	// copy credentials
	p.copyCredentials()

	// make the request
	url := fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io/query", p.IndexName, p.ProjectName, p.PineconeEnv)
	body := p.makeHTTPRequest(url, version, input)

	// parse the response
	var resp PineconeResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return ""
	}

	// return the first result
	return p.capText(resp.Matches[0].Metadata["text"].(string))
}

// Upsert upserts a document into the Pinecone index.
func (p *Pinecone) Upsert(version string, input string) string {
	// copy credentials
	p.copyCredentials()

	// make the request
	url := fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io/vectors/upsert", p.IndexName, p.ProjectName, p.PineconeEnv)
	body := p.makeHTTPRequest(url, version, input)

	// parse the response
	var resp PineconeResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return ""
	}

	// return the first result
	return resp.Matches[0].Metadata["text"].(string)
}

// makeHTTPRequest makes the HTTP request to the Pinecone API.
func (p *Pinecone) makeHTTPRequest(url string, version string, input string) []byte {
	// get the embeddings
	embeddings := EmbedText(input)

	// create the payload
	payload := QueryPayload{
		IncludeValues:   false,
		IncludeMetadata: true,
		Vector:          embeddings,
		TopK:            1000,
		Namespace:       version,
	}

	// convert the payload to JSON
	payloadBytes, _ := json.Marshal(payload)

	// make the request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	// get the response
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// return the response
	return body
}

// copyCredentials copies the credentials from the app config.
func (p *Pinecone) copyCredentials() {
	app := GetApp()

	p.ApiKey = app.PineconeKey
	p.ProjectName = app.PineconeProjectId
	p.PineconeEnv = app.PineconeEnv
}

// capText caps the text at 1000 characters.
func (p *Pinecone) capText(text string) string {
	if len(text) > 1000 {
		return text[:1000]
	}

	return text
}
