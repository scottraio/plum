package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type OpenAIConfig struct {
	OpenAIToken string
}

type OpenAI struct {
	Embedding
	Config OpenAIConfig
}

type EmbedResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func InitOpenAI(apiKey string) Embedding {
	config := &OpenAIConfig{
		OpenAIToken: apiKey,
	}

	openai := OpenAI{
		Config: *config,
	}

	return &openai
}

// EmbedText returns the embeddings for the given text using the OpenAI API.
func (ai *OpenAI) EmbedText(text string) []float32 {
	// Generate embeddings for the search query text using the OpenAI API client
	embeddingsBytes, err := ai.GetEmbeddings(text)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the data array from the JSON object
	var embedResp EmbedResponse
	err = json.Unmarshal(embeddingsBytes, &embedResp)
	if err != nil {
		log.Fatal(err)
	}

	// Get the embeddings array from the data map
	embeddingsArr := embedResp.Data[0].Embedding

	// Convert the embeddings array to a float64 slice
	embeddings := make([]float64, len(embeddingsArr))
	for i, v := range embeddingsArr {
		embeddings[i] = v
	}

	converted := convertFloat64ArrayToFloat32Array(embeddings)

	return converted
}

// GetEmbeddings returns the embeddings for the given text using the OpenAI API.
func (ai *OpenAI) GetEmbeddings(input string) ([]byte, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"input": input,
		"model": "text-embedding-ada-002",
	})
	if err != nil {
		return nil, fmt.Errorf("error serializing payload to JSON: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ai.Config.OpenAIToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP response: %w", err)
	}

	return respData, nil
}

func convertFloat64ArrayToFloat32Array(float64Array []float64) []float32 {
	float32Array := make([]float32, len(float64Array))
	for i, v := range float64Array {
		float32Array[i] = float32(v)
	}
	return float32Array
}
