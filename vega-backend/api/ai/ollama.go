package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const ollamaURL = "http://localhost:11434/api/chat"

type OllamaAPI struct {
	Model string
}

func NewOllamaAPI(model string) *OllamaAPI {
	return &OllamaAPI{Model: model}
}

func (o *OllamaAPI) Complete(messages []Message) (*Message, error) {
	body, err := json.Marshal(map[string]any{
		"model":    o.Model,
		"messages": messages,
		"stream":   false,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, ollamaURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var raw struct {
		Message Message `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &raw.Message, nil
}
