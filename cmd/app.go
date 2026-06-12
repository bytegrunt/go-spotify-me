package cmd

import (
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
	viewLoading
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
	dataProvider    DataProvider
}

func (m appModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, tea.EnterAltScreen, tea.ClearScreen, tea.WindowSize())

	if m.clientID != "" && m.currentView == viewLoading {
		cmds = append(cmds, performLogin(m.dataProvider))
	}

	return tea.Batch(cmds...)
}

func InitialAppModel(clientID string) appModel {
	dp := NewDataProvider()

	if clientID == "" {
		ti := textinput.New()
		ti.Placeholder = "Enter your Spotify Client ID"
		ti.Focus()
		ti.CharLimit = 100
		ti.Width = 50

		return appModel{
			currentView:  viewEnterClientID,
			textInput:    ti,
			dataProvider: dp,
		}
	}

	return appModel{
		currentView:  viewLoading,
		clientID:     clientID,
		dataProvider: dp,
	}
}
