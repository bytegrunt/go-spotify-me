package cmd

import (
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

	me, err := fetchMe()
	if err != nil {
		me = Me{
			DisplayName: "Unknown",
			Email:       "Unknown",
			Product:     "Unknown",
			ProfileURL:  "Unknown",
		}
	}

	return appModel{
		currentView: viewMenu,
		clientID:    clientID,
		me:          me,
	}
}
