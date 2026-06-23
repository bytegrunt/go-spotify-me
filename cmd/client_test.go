package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpotifyClient_Success(t *testing.T) {
	// Mock server that returns a successful JSON response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test_token" {
			t.Errorf("Expected Authorization header 'Bearer test_token', got '%s'", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true, "message": "hello"}`))
	}))
	defer server.Close()

	client := &DefaultSpotifyClient{
		httpClient: &http.Client{},
		tokenProvider: func() (string, bool) {
			return "test_token", true
		},
	}

	response, err := client.doRequest(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response["success"] != true {
		t.Errorf("Expected success=true in response, got %v", response["success"])
	}
	if response["message"] != "hello" {
		t.Errorf("Expected message='hello' in response, got %v", response["message"])
	}
}

func TestSpotifyClient_ErrorStatus(t *testing.T) {
	// Mock server that returns an error status code
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	}))
	defer server.Close()

	client := &DefaultSpotifyClient{
		httpClient: &http.Client{},
		tokenProvider: func() (string, bool) {
			return "test_token", true
		},
	}

	_, err := client.doRequest(server.URL)
	if err == nil {
		t.Fatalf("Expected error for non-200 status code, got nil")
	}

	expectedErrMsg := "API request failed with status code 404: Not Found"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestSpotifyClient_InvalidJSON(t *testing.T) {
	// Mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client := &DefaultSpotifyClient{
		httpClient: &http.Client{},
		tokenProvider: func() (string, bool) {
			return "test_token", true
		},
	}

	_, err := client.doRequest(server.URL)
	if err == nil {
		t.Fatalf("Expected error for invalid JSON, got nil")
	}
}
