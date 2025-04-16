/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bytegrunt/go-spotify-me/internal/logging"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var clientId string

const redirectURI = "http://127.0.0.1:6969/callback"
const authURL = "https://accounts.spotify.com/authorize"
const tokenURL = "https://accounts.spotify.com/api/token"

var codeVerifier string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Spotify using PKCE",
	Long: `Authenticate with Spotify using the PKCE flow. This command will
initiate the login process and open a browser for user authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if clientId is set
		if clientId == "" {
			log.Fatal("client-id is required. Set it via the environment variable SPOTIFY_CLIENT_ID or the --client-id flag.")
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
			err := refreshAccessToken(refreshToken)
			if err == nil {
				return // Successfully refreshed the token, exit the command
			}
			logging.DebugLog("Failed to refresh access token: %v", err)
			logging.DebugLog("Falling back to regular login flow.")
		}

		// Generate the code verifier and code challenge
		codeVerifier = generateCodeVerifier()
		codeChallenge := generateCodeChallenge(codeVerifier)

		// Generate the authorization URL
		authURLWithParams := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&scope=user-read-private user-read-email&code_challenge=%s&code_challenge_method=S256",
			authURL, url.QueryEscape(clientId), redirectURI, codeChallenge)

		logging.DebugLog("Generated authorization URL: %s", authURLWithParams)

		// Open the URL in the default browser
		err = openBrowser(authURLWithParams)
		if err != nil {
			log.Fatalf("Failed to open browser: %v", err)
		}

		// Start a local server to handle the callback
		startCallbackServer()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Set clientId from environment variable
	envClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if envClientId != "" {
		clientId = envClientId
	}

	// Add a CLI flag for client-id
	loginCmd.Flags().StringVar(&clientId, "client-id", clientId, "Spotify Client ID (overrides SPOTIFY_CLIENT_ID)")
}

// Generate a random code verifier
func generateCodeVerifier() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random generator
	verifier := make([]byte, 64)
	for i := range verifier {
		verifier[i] = byte(rng.Intn(26) + 97) // a-z
	}
	return string(verifier)
}

// Generate a code challenge from the code verifier
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// Start a local server to handle the callback
func startCallbackServer() {
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
			exchangeCodeForToken(code)
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

// Exchange the authorization code for an access token
func exchangeCodeForToken(code string) {
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("Failed to create token request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to exchange code for token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to get token: %s", body)
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Fatalf("Failed to parse token response: %v", err)
	}

	accessToken := tokenResponse["access_token"].(string)
	refreshToken := tokenResponse["refresh_token"].(string)

	// Save the access token to a hidden file
	saveAccessTokenToFile(accessToken, refreshToken)

	fmt.Println("Refresh Token stored successfully.")
}

// Save the access token and refresh token to a hidden file
func saveAccessTokenToFile(accessToken, refreshToken string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	filePath := filepath.Join(homeDir, ".go-spotify-cli")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %v", err)
	}
	defer file.Close()
	var data string

	// Attempt to store the refresh token in the keyring
	err = keyring.Set("go-spotify-cli", "refresh_token", refreshToken)
	if err != nil {
		log.Printf("Failed to store refresh token in keyring: %v", err)
		log.Println("Falling back to saving the refresh token in the hidden file.")
		data = fmt.Sprintf("access_token=%s\nrefresh_token=%s\n", accessToken, refreshToken)
	} else {
		data = fmt.Sprintf("access_token=%s", accessToken)
	}

	_, err = file.WriteString(data)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	log.Printf("Access token and refresh token saved to %s", filePath)
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

// Refresh the access token using the refresh token
func refreshAccessToken(refreshToken string) error {
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

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
	saveAccessTokenToFile(accessToken, refreshToken)

	logging.DebugLog("Access Token refreshed successfully.")
	return nil
}
