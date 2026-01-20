package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

// -----------------------------------------------------------------------------
// TOGETHER AI CLIENT — SAFE, MINIMAL, NON-AUTONOMOUS
//
// This client:
// - makes a single bounded request
// - enforces timeouts
// - does not retry aggressively
// - exposes no internal state
//
// Failure here MUST NOT affect VANTAGE core execution.
// -----------------------------------------------------------------------------

type togetherRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type togetherResponse struct {
	Output struct {
		Text string `json:"text"`
	} `json:"output"`
}

// callTogether performs a single advisory inference request.
func callTogether(prompt string) (string, error) {

	apiKey := os.Getenv("TOGETHER_API_KEY")
	if apiKey == "" {
		return "", errors.New("TOGETHER_API_KEY not set")
	}

	model := os.Getenv("VANTAGE_AI_MODEL")
	if model == "" {
		return "", errors.New("VANTAGE_AI_MODEL not set")
	}

	reqBody := togetherRequest{
		Model:       model,
		Prompt:      prompt,
		Temperature: 0.0, // Determinism enforced
		MaxTokens:   800, // Hard cap
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.together.xyz/v1/completions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var parsed togetherResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	return parsed.Output.Text, nil
}
