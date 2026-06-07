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
	dataProvider    DataProvider
}

func (m appModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, tea.EnterAltScreen, tea.ClearScreen, tea.WindowSize())

	if m.currentView == viewMenu {
		cmds = append(cmds, func() tea.Msg {
			err := m.dataProvider.Login()
			if err != nil {
				return initLoginMsg{err: fmt.Errorf("failed to log in: %w", err)}
			}
			me, err := m.dataProvider.FetchMe()
			if err != nil {
				return initLoginMsg{err: fmt.Errorf("failed to fetch user info: %w", err)}
			}
			return initLoginMsg{me: me}
		})
	}

	return tea.Batch(cmds...)
}

func InitialAppModel(clientID string) appModel {
	dp := &DefaultDataProvider{}

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
		artistTable:     artistTable,
		artistColWidths: artistColWidths,
		songTable:       songTable,
		songColWidths:   songColWidths,
		dataProvider:    dp,
		me: Me{
			DisplayName: "Loading...",
			Email:       "Loading...",
			Product:     "Loading...",
			ProfileURL:  "Loading...",
		},
	}
}
