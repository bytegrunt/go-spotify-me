package auth

import (
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/CyberGrit/go-spotify-me/internal/logging"
	"go.uber.org/zap"
)

var logger *zap.Logger

type AuthConfig struct {
	RedirectURI string
	AuthURL     string
	TokenURL    string
	ClientID    string
}

// Generate a random code verifier
func GenerateCodeVerifier() string {
	verifier := make([]byte, 64)
	_, err := cryptoRand.Read(verifier) // Use crypto/rand for secure random bytes
	if err != nil {
		logger.Fatal("Failed to generate secure random bytes", zap.Error(err))
	}

	// Convert bytes to a-z characters
	for i := range verifier {
		verifier[i] = (verifier[i] % 26) + 97 // a-z
	}
	return string(verifier)
}

// Generate a code challenge from the code verifier
func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// Exchange the authorization code for an access token
func ExchangeCodeForToken(authConfig AuthConfig, store TokenStore, code, codeVerifier string) {
	data := url.Values{}
	data.Set("client_id", authConfig.ClientID)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", authConfig.RedirectURI)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", authConfig.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Fatal("Failed to create token request", zap.Error(err))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal("Failed to exchange code for token", zap.Error(err))
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error("Error closing response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Fatal("Failed to get token", zap.String("body", string(body)))
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		logger.Fatal("Failed to parse token response", zap.Error(err))
	}

	accessToken := tokenResponse["access_token"].(string)
	refreshToken := tokenResponse["refresh_token"].(string)
	expiresIn := int(tokenResponse["expires_in"].(float64)) // Convert to int

	// Calculate expiration time
	expirationTime := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Save the access token, refresh token, and expiration time using TokenStore
	err = store.SaveTokens(accessToken, refreshToken, expirationTime)
	if err != nil && logger != nil {
		logger.Error("Failed to save tokens", zap.Error(err))
	}

	fmt.Println("Refresh Token stored successfully.")
}

// Refresh the access token using the refresh token
func RefreshAccessToken(authConfig AuthConfig, store TokenStore, refreshToken string) error {
	data := url.Values{}
	data.Set("client_id", authConfig.ClientID)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", authConfig.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error("Error closing response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to refresh token: %s", body)
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	accessToken := tokenResponse["access_token"].(string)

	// Save the new access token using TokenStore
	err = store.SaveTokens(accessToken, refreshToken, time.Now().Add(3600*time.Second)) // Assuming 1 hour expiration
	if err != nil && logger != nil {
		logger.Error("Failed to save tokens", zap.Error(err))
	}

	logging.DebugLog("Access Token refreshed successfully.")
	return nil
}

