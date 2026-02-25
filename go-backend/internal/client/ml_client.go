package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MLClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewMLClient(baseURL string) *MLClient {
	return &MLClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Embed request/response
type EmbedRequest struct {
	Texts []string `json:"texts"`
}

type EmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
	Model      string      `json:"model"`
	Count      int         `json:"count"`
}

func (c *MLClient) Embed(texts []string) (*EmbedResponse, error) {
	body, err := json.Marshal(EmbedRequest{Texts: texts})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/embed",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("ml service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml service returned status %d", resp.StatusCode)
	}

	var result EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Chunk request/response
type ChunkRequest struct {
	Text      string `json:"text"`
	ChunkSize int    `json:"chunk_size"`
	Overlap   int    `json:"overlap"`
}

type ChunkResponse struct {
	Chunks []string       `json:"chunks"`
	Stats  map[string]any `json:"stats"`
}

func (c *MLClient) Chunk(text string, chunkSize, overlap int) (*ChunkResponse, error) {
	body, err := json.Marshal(ChunkRequest{
		Text:      text,
		ChunkSize: chunkSize,
		Overlap:   overlap,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/chunk",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("ml service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml service returned status %d", resp.StatusCode)
	}

	var result ChunkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Ping ML service
func (c *MLClient) Ping() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("ml service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ml service unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
