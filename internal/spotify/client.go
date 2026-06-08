package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Client is an interface for making requests to the Spotify API
type Client interface {
	Get(token string, url string) (map[string]interface{}, error)
}

// DefaultClient is the concrete implementation of the Client interface
type DefaultClient struct {
	HTTPClient *http.Client
}

// NewClient creates a new DefaultClient
func NewClient() *DefaultClient {
	return &DefaultClient{
		HTTPClient: &http.Client{},
	}
}

// Get makes a GET request to the Spotify API and returns the response or an error
func (c *DefaultClient) Get(token string, url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return response, nil
}
