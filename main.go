package main

import (
	"log"

	"github.com/bytegrunt/go-spotify-me/cmd"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
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
