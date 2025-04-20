/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bytegrunt/go-spotify-me/cmd"
	"github.com/bytegrunt/go-spotify-me/internal/auth"
	"github.com/bytegrunt/go-spotify-me/internal/logging"
)

func main() {

	_, isValid := auth.GetValidAccessToken()
		if !isValid {
			logging.DebugLog("Access token is not valid. Running login command...")
			cmd.Login()
			_, isValid = auth.GetValidAccessToken()
			if !isValid {
				fmt.Println("Failed to obtain a valid access token after login.")
				return
			}
		}

    // Initialize the app model
    p := tea.NewProgram(cmd.InitialAppModel(), tea.WithAltScreen())
    if _,err := p.Run(); err != nil {
        log.Fatalf("Error starting TUI: %v", err)
    }
}
