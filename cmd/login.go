package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bytegrunt/go-spotify-me/internal/auth"
	"github.com/zalando/go-keyring"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		logger.Fatal("Failed to initialize zap logger", zap.Error(err))
	}
}

func Login() error {
	clientID, err := GetClientID()
	if err != nil {
		return fmt.Errorf("failed to get client ID: %w", err)
	}

	authConfig := auth.AuthConfig{
		RedirectURI: "http://127.0.0.1:9000/callback",
		AuthURL:     "https://accounts.spotify.com/authorize",
		TokenURL:    "https://accounts.spotify.com/api/token",
		ClientID:    clientID,
	}

	_, isValid := auth.GetValidAccessToken()
	if isValid {
		return nil
	}

	// Check for refresh token in the keyring
	refreshToken, err := keyring.Get("go-spotify-me-cli", "refresh_token")
	if err != nil {
		logger.Debug("Refresh token not found in keyring", zap.Error(err))

		// Check for refresh token in the hidden file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Fatal("Failed to get user home directory", zap.Error(err))
		}

		filePath := filepath.Join(homeDir, ".go-spotify-me-cli")

		// Validate that the filePath is within the user's home directory
		if !strings.HasPrefix(filePath, homeDir) {
			logger.Fatal("Invalid file path", zap.String("filePath", filePath))
		}

		// Attempt to read the file
		data, err := os.ReadFile(filePath)
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "refresh_token=") {
					refreshToken = strings.TrimPrefix(line, "refresh_token=")
					break
				}
			}
		}
	}

	// If a refresh token is found, try to refresh the access token
	if refreshToken != "" {
		logger.Debug("Using existing refresh token to get a new access token.")
		err := auth.RefreshAccessToken(authConfig, refreshToken)
		if err == nil {
			return nil // Successfully refreshed the token, exit the command
		}
		logger.Debug("Failed to refresh access token", zap.Error(err))
		logger.Debug("Falling back to regular login flow.")
	}

	// Generate the code verifier and code challenge
	codeVerifier := auth.GenerateCodeVerifier()
	codeChallenge := auth.GenerateCodeChallenge(codeVerifier)

	// Generate the authorization URL
	authURLWithParams := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&scope=user-read-private user-read-email user-top-read&code_challenge=%s&code_challenge_method=S256",
		authConfig.AuthURL, url.QueryEscape(authConfig.ClientID), authConfig.RedirectURI, codeChallenge)

	logger.Debug("Generated authorization URL", zap.String("url", authURLWithParams))

	// Open the URL in the default browser
	err = openBrowser(authURLWithParams)
	if err != nil {
		logger.Error("Failed to open browser", zap.Error(err))
	}

	// Start a local server to handle the callback
	startCallbackServer(authConfig, codeVerifier)
	return nil
}

// GetClientID retrieves the Client ID from the keyring or environment variable.
func GetClientID() (string, error) {
	// Check the keyring for the Client ID
	clientID, err := keyring.Get("go-spotify-me-cli", "client_id")
	if err == nil {
		return clientID, nil
	}

	// Fallback to environment variable
	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	if clientID != "" {
		return clientID, nil
	}

	return "", nil
}

// Start a local server to handle the callback
func startCallbackServer(authConfig auth.AuthConfig, codeVerifier string) {
	var wg sync.WaitGroup
	wg.Add(1) // Add one task to the WaitGroup

	server := &http.Server{
		Addr:              "127.0.0.1:9000",
		ReadHeaderTimeout: 5 * time.Second, // Set a timeout to mitigate Slowloris attacks
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		if code == "" {
			http.Error(w, "Authorization code not found", http.StatusBadRequest)
			return
		}

		if _, err := fmt.Fprintln(w, "Authorization successful! You can close this window."); err != nil {
			logger.Error("Error writing response", zap.Error(err))
		}

		// Exchange the authorization code for an access token
		go func() {
			auth.ExchangeCodeForToken(authConfig, code, codeVerifier)
			if err := server.Close(); err != nil {
				logger.Error("Error closing server", zap.Error(err))
			}
			wg.Done() // Mark the task as done
		}()
	})

	// Start the server
	go func() {
		logger.Info("Waiting for the authorization code...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	go func() {
		time.Sleep(180 * time.Second)
		logger.Info("Timeout reached. Shutting down the server.")
		if err := server.Close(); err != nil {
			logger.Error("Error closing server", zap.Error(err))
		}
		wg.Done() // Mark the task as done if timeout occurs
	}()

	wg.Wait() // Wait for the server to close
}

// Function to open the browser
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "windows": // Windows
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // Linux and other OS
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
