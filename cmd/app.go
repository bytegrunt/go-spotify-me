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
	currentView     viewType
	clientID        string
	me              Me              // User information
	textInput       textinput.Model // Text input for Client ID
	artists         APIResponse
	songs           APIResponse
	artistTable     table.Model // Table for artists
	artistColWidths []int
	songTable       table.Model // Table for songs
	songColWidths   []int
	windowSize      tea.WindowSizeMsg
	err             error
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
			err: fmt.Errorf("failed to log in: %w", err),
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

	// Initialize artist table
	artistColWidths := calculateColumnWidths(100, []float64{0.4, 0.4, 0.2})
	artistTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "Name", Width: 40},
			{Title: "Genres", Width: 50},
			{Title: "Popularity", Width: 10},
		}),
		table.WithFocused(false),
	)

	// Initialize song table
	songColWidths := calculateColumnWidths(100, []float64{0.4, 0.3, 0.2, 0.1})
	songTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "Name", Width: 40},
			{Title: "Artist", Width: 20},
			{Title: "Album", Width: 30},
			{Title: "Popularity", Width: 10},
		}),
		table.WithFocused(false),
	)

	return appModel{
		currentView:     viewMenu,
		clientID:        clientID,
		me:              me,
		artistTable:     artistTable,
		artistColWidths: artistColWidths,
		songTable:       songTable,
		songColWidths:   songColWidths,
	}
}
