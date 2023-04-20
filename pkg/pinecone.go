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

func (p *Pinecone) Query(input string) string {
	body := p.makeHTTPRequest(input)

	var resp PineconeResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return ""
	}

	return p.capText(resp.Matches[0].Metadata["text"].(string))
}

func (p *Pinecone) makeHTTPRequest(input string) []byte {
	app := GetApp()

	p.ApiKey = app.PineconeKey
	p.ProjectName = "9eaeaa7"
	p.PineconeEnv = app.PineconeEnv

	url := fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io/query", p.IndexName, p.ProjectName, p.PineconeEnv)

	embeddings := EmbedText(input)
	payload := QueryPayload{
		IncludeValues:   false,
		IncludeMetadata: true,
		Vector:          embeddings,
		TopK:            1000,
		Namespace:       "40",
	}

	// convert the payload to JSON
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func (p *Pinecone) capText(text string) string {
	if len(text) > 1000 {
		return text[:1000]
	}

	return text
}
