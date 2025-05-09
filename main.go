package main

import (
	"fmt"
	"os"

	"github.com/bytegrunt/go-spotify-me/cmd"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

// Removed the init function and replaced it with an explicit InitializeLogger function.
var logger *zap.Logger

func InitializeLogger() error {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize zap logger: %w", err)
	}
	return nil
}

func main() {
	// Initialize the logger
	if err := InitializeLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Check for the --clear-config flag
	for _, arg := range os.Args {
		if arg == "--clear-config" {
			if err := cmd.ClearConfig(); err != nil {
				logger.Fatal("Failed to clear configuration", zap.Error(err))
			}
			fmt.Println("Configuration cleared successfully.")
			return
		}
	}

	clientID, err := cmd.GetClientID()
	if err != nil {
		logger.Fatal("Failed to retrieve client ID", zap.Error(err))
	}

	// Initialize the app model with the client ID
	p := tea.NewProgram(cmd.InitialAppModel(clientID), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Fatal("Error starting TUI", zap.Error(err))
	}
}
