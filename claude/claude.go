package claude

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/kr/pretty"
)

type Response struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Role         string    `json:"role"`
	Content      []Content `json:"content"`
	Model        string    `json:"model"`
	StopReason   string    `json:"stop_reason"`
	StopSequence *string   `json:"stop_sequence"`
	Usage        Usage     `json:"usage"`
}

// Content struct to hold the content of the response
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage struct to hold the usage details of the response
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func GetMessage(history []map[string]string, input string) (string, error) {
	// Anthropic API endpoint for sending messages
	url := "https://api.anthropic.com/v1/messages"

	// Use an environment variable for the API key, or replace os.Getenv("ANTHROPIC_API_KEY") with your actual key in quotes
	apiKey := ""

	// Prepare the data for the request
	data := map[string]interface{}{
		"model":       "claude-3-opus-20240229", // Ensure you use the correct model identifier
		"system":      "You are an AI assistant with a deep understanding of dream interpretation and symbolism. Your task is to provide users with insightful and meaningful analyses of the symbols, emotions, and narratives present in their dreams. Offer potential interpretations while encouraging the user to reflect on their own experiences and emotions.",
		"max_tokens":  1024,
		"temperature": 1,
		"messages":    history,
	}

	// Marshal the data into a JSON []byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Add the required headers
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("anthropic-version", "2023-06-01") // Adjust the version date as needed
	req.Header.Add("content-type", "application/json")

	// Initialize a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and print the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}
	if len(response.Content) == 0 {
		pretty.Println("empty response", response)
		return "", errors.New("empty response")
	}

	return response.Content[0].Text, nil
}
