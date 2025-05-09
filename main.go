package main

import (
	"fmt"
	"os"

	"github.com/bytegrunt/go-spotify-me/cmd"
	tea "github.com/charmbracelet/bubbletea"
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

func main() {
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
