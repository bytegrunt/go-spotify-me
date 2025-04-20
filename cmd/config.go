package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

// ClearConfig removes the client_id and refresh_token from the keyring
// and clears the .go-spotify-me-cli file in the user's home directory.
func ClearConfig() error {
	// Remove client_id from the keyring
	if err := keyring.Delete("go-spotify-me-cli", "client_id"); err != nil {
		fmt.Printf("Failed to delete client_id from keyring: %v\n", err)
	}

	// Remove refresh_token from the keyring
	if err := keyring.Delete("go-spotify-me-cli", "refresh_token"); err != nil {
		fmt.Printf("Failed to delete refresh_token from keyring: %v\n", err)
	}

	// Clear the .go-spotify-me-cli file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, ".go-spotify-me-cli")
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete configuration file: %w", err)
	}

	return nil
}
