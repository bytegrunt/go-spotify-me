/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
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
	"github.com/bytegrunt/go-spotify-me/internal/logging"
	"github.com/zalando/go-keyring"
)

var clientId string
var forceLogin bool // Flag to force login

// loginCmd represents the login command
func Login() {
	// Check if clientId is set
	if clientId == "" {
		log.Fatal("client-id is required. Set it via the environment variable SPOTIFY_CLIENT_ID or the --client-id flag.")
	}

	authConfig := auth.AuthConfig{
		RedirectURI: "http://127.0.0.1:6969/callback",
		AuthURL:     "https://accounts.spotify.com/authorize",
		TokenURL:    "https://accounts.spotify.com/api/token",
		ClientID:    clientId,
	}

	// Skip refresh token checks if force flag is set
	if !forceLogin {
		_, isValid := auth.GetValidAccessToken()
		if isValid {
			return
		}

		// Check for refresh token in the keyring
		refreshToken, err := keyring.Get("go-spotify-cli", "refresh_token")
		if err != nil {
			logging.DebugLog("Refresh token not found in keyring: %v", err)

			// Check for refresh token in the hidden file
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Failed to get user home directory: %v", err)
			}

			filePath := filepath.Join(homeDir, ".go-spotify-cli")
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
			logging.DebugLog("Using existing refresh token to get a new access token.")
			err := auth.RefreshAccessToken(authConfig, refreshToken)
			if err == nil {
				return // Successfully refreshed the token, exit the command
			}
			logging.DebugLog("Failed to refresh access token: %v", err)
			logging.DebugLog("Falling back to regular login flow.")
		}
	}

	// Generate the code verifier and code challenge
	codeVerifier := auth.GenerateCodeVerifier()
	codeChallenge := auth.GenerateCodeChallenge(codeVerifier)

	// Generate the authorization URL
	authURLWithParams := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&scope=user-read-private user-read-email user-top-read&code_challenge=%s&code_challenge_method=S256",
		authConfig.AuthURL, url.QueryEscape(authConfig.ClientID), authConfig.RedirectURI, codeChallenge)

	logging.DebugLog("Generated authorization URL: %s", authURLWithParams)

	// Open the URL in the default browser
	err := openBrowser(authURLWithParams)
	if err != nil {
		log.Fatalf("Failed to open browser: %v", err)
	}

	// Start a local server to handle the callback
	startCallbackServer(authConfig, codeVerifier)
}

func init() {
	// Set clientId from environment variable
	envClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if envClientId != "" {
		clientId = envClientId
	}
}

// Start a local server to handle the callback
func startCallbackServer(authConfig auth.AuthConfig, codeVerifier string) {
	var wg sync.WaitGroup
	wg.Add(1) // Add one task to the WaitGroup

	server := &http.Server{
		Addr: "127.0.0.1:6969",
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		if code == "" {
			http.Error(w, "Authorization code not found", http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, "Authorization successful! You can close this window.")

		// Exchange the authorization code for an access token
		go func() {
			auth.ExchangeCodeForToken(authConfig, code, codeVerifier)
			server.Close() // Close the server after processing the token
			wg.Done()      // Mark the task as done
		}()
	})

	// Start the server
	go func() {
		log.Println("Waiting for the authorization code...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for the server to finish or timeout after 120 seconds
	go func() {
		time.Sleep(120 * time.Second)
		log.Println("Timeout reached. Shutting down the server.")
		server.Close()
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
