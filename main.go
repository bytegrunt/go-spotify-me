package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bytegrunt/go-spotify-me/cmd"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check for the --clear-config flag
	for _, arg := range os.Args {
		if arg == "--clear-config" {
			if err := cmd.ClearConfig(); err != nil {
				log.Fatalf("Failed to clear configuration: %v", err)
			}
			fmt.Println("Configuration cleared successfully.")
			return
		}
	}

	clientID, err := cmd.GetClientID()
	if err != nil {
		log.Fatalf("Failed to retrieve client ID: %v", err)
	}

	// Initialize the app model with the client ID
	p := tea.NewProgram(cmd.InitialAppModel(clientID), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting TUI: %v", err)
	}
}
