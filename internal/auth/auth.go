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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytegrunt/go-spotify-me/internal/logging"
	"github.com/zalando/go-keyring"
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
func ExchangeCodeForToken(authConfig AuthConfig, code, codeVerifier string) {
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

	// Save the access token, refresh token, and expiration time to a hidden file
	SaveAccessTokenToFile(accessToken, refreshToken, expirationTime)

	fmt.Println("Refresh Token stored successfully.")
}

// Save the access token, refresh token, and expiration time to a hidden file
func SaveAccessTokenToFile(accessToken, refreshToken string, expirationTime time.Time) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal("Failed to get user home directory", zap.Error(err))
	}

	filePath := filepath.Join(homeDir, ".go-spotify-me-cli")

	// Validate that the filePath is within the user's home directory
	if !strings.HasPrefix(filePath, homeDir) {
		logger.Fatal("Invalid file path", zap.String("filePath", filePath))
	}

	// Attempt to open the file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		logger.Fatal("Failed to open file for writing", zap.Error(err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Error closing file", zap.Error(err))
		}
	}()

	var data string

	// Attempt to store the refresh token in the keyring
	err = keyring.Set("go-spotify-me-cli", "refresh_token", refreshToken)
	if err != nil {
		logger.Error("Failed to store refresh token in keyring", zap.Error(err))
		logger.Info("Falling back to saving the refresh token in the hidden file.")
		data = fmt.Sprintf("access_token=%s\nrefresh_token=%s\nexpires_at=%s\n", accessToken, refreshToken, expirationTime.Format(time.RFC3339))
	} else {
		data = fmt.Sprintf("access_token=%s\nexpires_at=%s\n", accessToken, expirationTime.Format(time.RFC3339))
	}

	_, err = file.WriteString(data)
	if err != nil {
		logger.Fatal("Failed to write to file", zap.Error(err))
	}

	logging.DebugLog("Access token saved to %s", filePath)
}

// Refresh the access token using the refresh token
func RefreshAccessToken(authConfig AuthConfig, refreshToken string) error {
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

	// Save the new access token to a hidden file
	SaveAccessTokenToFile(accessToken, refreshToken, time.Now().Add(3600*time.Second)) // Assuming 1 hour expiration

	logging.DebugLog("Access Token refreshed successfully.")
	return nil
}

// Check if the token is still valid and return the token if valid
func GetValidAccessToken() (string, bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal("Failed to get user home directory", zap.Error(err))
	}

	filePath := filepath.Join(homeDir, ".go-spotify-me-cli")

	// Validate that the filePath is within the user's home directory
	if !strings.HasPrefix(filePath, homeDir) {
		logger.Fatal("Invalid file path", zap.String("filePath", filePath))
	}

	// Attempt to open the file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		logger.Fatal("Failed to open file", zap.Error(err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Error closing file", zap.Error(err))
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read token file", zap.Error(err))
		return "", false
	}

	lines := strings.Split(string(data), "\n")
	var accessToken, expiresAtStr string
	for _, line := range lines {
		if strings.HasPrefix(line, "access_token=") {
			accessToken = strings.TrimPrefix(line, "access_token=")
		} else if strings.HasPrefix(line, "expires_at=") {
			expiresAtStr = strings.TrimPrefix(line, "expires_at=")
		}
	}

	if accessToken == "" || expiresAtStr == "" {
		return "", false
	}

	expirationTime, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		logger.Error("Failed to parse expiration time", zap.Error(err))
		return "", false
	}

	if time.Now().After(expirationTime) {
		return "", false
	}

	return accessToken, true
}
