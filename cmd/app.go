package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type viewType int

const (
	viewMenu viewType = iota
	viewArtists
	viewSongs
	viewEnterClientID
)

type appModel struct {
	currentView viewType
	clientID    string
	me          Me              // User information
	textInput   textinput.Model // Text input for Client ID
	artists     APIResponse
	songs       APIResponse
	artistTable table.Model // Table for artists
    songTable   table.Model // Table for songs
	windowSize  tea.WindowSizeMsg
	err         error
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.ClearScreen,
		tea.WindowSize(),
	)
}

func InitialAppModel(clientID string) appModel {
	if clientID == "" {
		ti := textinput.New()
		ti.Placeholder = "Enter your Spotify Client ID"
		ti.Focus()
		ti.CharLimit = 100
		ti.Width = 50

		return appModel{
			currentView: viewEnterClientID,
			textInput:   ti,
		}
	}

	err := Login()
	if err != nil {
		return appModel{
			err: fmt.Errorf("failed to log in: %v", err),
		}
	}

	me, err := fetchMe()
	if err != nil {
		me = Me{
			DisplayName: "Unknown",
			Email:       "Unknown",
			Product:     "Unknown",
			ProfileURL:  "Unknown",
		}
	}

	// Define columns for the tables
    columns := []table.Column{
        {Title: "Name", Width: 30},
        {Title: "Genres", Width: 40},
        {Title: "Popularity", Width: 10},
    }

    // Initialize artist table
    artistTable := table.New(
        table.WithColumns(columns),
        table.WithFocused(true),
    )

    // Initialize song table
    songTable := table.New(
        table.WithColumns([]table.Column{
            {Title: "Name", Width: 30},
            {Title: "Artist", Width: 30},
            {Title: "Album", Width: 30},
            {Title: "Popularity", Width: 10},
        }),
        table.WithFocused(true),
    )

	return appModel{
		currentView: viewMenu,
        clientID:    clientID,
        me:          me,
        artistTable: artistTable,
        songTable:   songTable,
	}
}
