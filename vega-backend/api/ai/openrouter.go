package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type openRouterResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type OpenRouterAPI struct {
	Model  string
	apiKey string
}

func NewOpenRouterAPI(model string) (*OpenRouterAPI, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY is not set")
	}
	return &OpenRouterAPI{Model: model, apiKey: apiKey}, nil
}

func (o *OpenRouterAPI) Complete(messages []Message) (*Message, error) {
	body, err := json.Marshal(CompletionRequest{Model: o.Model, Messages: messages})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenRouter returned status %d", resp.StatusCode)
	}

	var result openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result.Choices[0].Message, nil
}
