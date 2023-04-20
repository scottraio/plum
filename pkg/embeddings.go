package plum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetEmbeddings(input string) ([]byte, error) {
	app := GetApp()

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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.OpenAIToken))

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

func EmbedText(text string) []float32 {
	// Generate embeddings for the search query text using the OpenAI API client
	embeddingsBytes, err := GetEmbeddings(text)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the data array from the JSON object
	var data map[string]interface{}
	err = json.Unmarshal(embeddingsBytes, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Get the embeddings array from the data map
	embeddingsArr := data["data"].([]interface{})[0].(map[string]interface{})["embedding"].([]interface{})

	// Convert the embeddings array to a float64 slice
	embeddings := make([]float64, len(embeddingsArr))
	for i, v := range embeddingsArr {
		embeddings[i] = v.(float64)
	}

	return convertFloat64ArrayToFloat32Array(embeddings)
}

func convertFloat64ArrayToFloat32Array(float64Array []float64) []float32 {
	float32Array := make([]float32, len(float64Array))
	for i, v := range float64Array {
		float32Array[i] = float32(v)
	}
	return float32Array
}
