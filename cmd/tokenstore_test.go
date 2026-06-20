package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zalando/go-keyring"
)

func TestOSTokenStore_SaveAndGetRefreshToken(t *testing.T) {
	// Mock keyring for success
	keyring.MockInit()

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	store := NewOSTokenStore()
	err := store.SaveTokens("mock-access", "mock-refresh", time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	token, err := store.GetRefreshToken()
	if err != nil {
		t.Fatalf("Expected no error getting refresh token, got %v", err)
	}
	if token != "mock-refresh" {
		t.Fatalf("Expected mock-refresh, got %v", token)
	}
}

func TestOSTokenStore_SaveAndGetValidAccessToken(t *testing.T) {
	keyring.MockInit()
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	store := NewOSTokenStore()
	err := store.SaveTokens("mock-access", "mock-refresh", time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	token, valid := store.GetValidAccessToken()
	if !valid {
		t.Fatalf("Expected valid access token")
	}
	if token != "mock-access" {
		t.Fatalf("Expected mock-access, got %v", token)
	}
}

func TestOSTokenStore_ExpiredAccessToken(t *testing.T) {
	keyring.MockInit()
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	store := NewOSTokenStore()
	// Expiration is in the past
	err := store.SaveTokens("mock-access", "mock-refresh", time.Now().Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, valid := store.GetValidAccessToken()
	if valid {
		t.Fatalf("Expected invalid access token since it is expired")
	}
}

func TestOSTokenStore_FallbackToFile(t *testing.T) {
	// Mock keyring to return an error to simulate fallback
	keyring.MockInitWithError(fmt.Errorf("keyring unavailable"))

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	store := NewOSTokenStore()
	err := store.SaveTokens("mock-access-file", "mock-refresh-file", time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Expected no error on fallback to file, got %v", err)
	}

	// Verify it can get the refresh token from file
	token, err := store.GetRefreshToken()
	if err != nil {
		t.Fatalf("Expected no error getting refresh token from file, got %v", err)
	}
	if token != "mock-refresh-file" {
		t.Fatalf("Expected mock-refresh-file, got %v", token)
	}

	// Verify file content explicitly just to be sure
	filePath := filepath.Join(tempDir, ".go-spotify-me-cli")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Expected token file to exist: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Fatalf("Expected file to not be empty")
	}
}
