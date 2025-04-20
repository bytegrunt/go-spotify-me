package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zalando/go-keyring"
)

type viewType int

const (
	viewMenu viewType = iota
	viewLogin
	viewArtists
	viewSongs
	viewEnterClientID
)

type appModel struct {
	currentView viewType
	clientID    string
	textInput   textinput.Model // Text input for Client ID
	artists     APIResponse
	songs       APIResponse
	windowSize  tea.WindowSizeMsg
	err         error
}

type switchToArtistsMsg struct {
	response APIResponse
}

type switchToSongsMsg struct {
	response APIResponse
}

type errMsg struct {
	err error
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
	return appModel{
		currentView: viewMenu,
		clientID:    clientID,
	}
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			if m.currentView != viewMenu {
				m.currentView = viewMenu
				return m, nil
			}
			return m, tea.Quit

		case "l", "L":
			m.currentView = viewLogin
			return m, nil

		case "a", "A":
			// Switch to the Artists view
			return m, func() tea.Msg {
				response, err := fetchArtistsPage("https://api.spotify.com/v1/me/top/artists")
				if err != nil {
					return errMsg{err}
				}
				return switchToArtistsMsg{response}
			}

		case "s", "S":
			// Switch to the Songs view
			return m, func() tea.Msg {
				response, err := fetchSongsPage("https://api.spotify.com/v1/me/top/tracks")
				if err != nil {
					return errMsg{err}
				}
				return switchToSongsMsg{response}
			}

		case "enter":
			// Handle entering the Client ID
			if m.currentView == viewEnterClientID {
				m.clientID = m.textInput.Value()
				err := keyring.Set("go-spotify-cli", "client_id", m.clientID)
				if err != nil {
					m.err = fmt.Errorf("failed to store client ID in keyring: %v", err)
					return m, nil
				}
				m.currentView = viewLogin
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.windowSize = msg

	case switchToArtistsMsg:
		m.artists = msg.response
		m.currentView = viewArtists

	case switchToSongsMsg:
		m.songs = msg.response
		m.currentView = viewSongs

	case errMsg:
		m.err = msg.err
	}

	// Update the text input model
	if m.currentView == viewEnterClientID {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m appModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
	}

	switch m.currentView {
	case viewMenu:
		return "\nWelcome to go-spotify\n\nPress L to Login\nPress A for Top Artists\nPress S for Top Songs\nPress Q to quit"
	case viewLogin:
		return m.renderLogin()
	case viewArtists:
		return m.renderArtists()
	case viewSongs:
		return m.renderSongs()
	case viewEnterClientID:
		return m.renderEnterClientID()
	default:
		return "Unknown view"
	}
}

func (m appModel) renderEnterClientID() string {
	return fmt.Sprintf(
		"Enter your Spotify Client ID:\n\n%s\n\nPress Enter to confirm, or Esc to quit.",
		m.textInput.View(),
	)
}

func (m appModel) renderLogin() string {
	return "Login functionality is not yet implemented.\nPress Q to go back."
}

func (m appModel) renderArtists() string {
	var s strings.Builder

	// Dynamic column widths
	totalWidth := m.windowSize.Width
	if totalWidth == 0 {
		totalWidth = 100 // fallback
	}
	colName := 25
	colGenre := totalWidth - colName - 12 - 5 // leave space for padding and popularity
	colPop := 10

	s.WriteString(fmt.Sprintf("%s %s %s\n",
		truncateOrPad("Name", colName),
		truncateOrPad("Genres", colGenre),
		truncateOrPad("Popularity", colPop),
	))
	s.WriteString(strings.Repeat("-", totalWidth) + "\n")

	for _, artist := range m.artists.Artists {
		s.WriteString(fmt.Sprintf("%s %s %d\n",
			truncateOrPad(artist.Name, colName),
			truncateOrPad(artist.Genres, colGenre),
			artist.Popularity,
		))
	}
	s.WriteString("\n[n] next, [p] previous, [q] back to menu\n")
	return s.String()
}

func (m appModel) renderSongs() string {
	var s strings.Builder

	totalWidth := m.windowSize.Width
	if totalWidth == 0 {
		totalWidth = 100
	}
	colName := 20
	colArtist := 20
	colAlbum := totalWidth - colName - colArtist - 12 - 5
	colPop := 10

	s.WriteString(fmt.Sprintf("%s %s %s %s\n",
		truncateOrPad("Name", colName),
		truncateOrPad("Artist", colArtist),
		truncateOrPad("Album", colAlbum),
		truncateOrPad("Popularity", colPop),
	))
	s.WriteString(strings.Repeat("-", totalWidth) + "\n")

	for _, song := range m.songs.Songs {
		s.WriteString(fmt.Sprintf("%s %s %s %d\n",
			truncateOrPad(song.Name, colName),
			truncateOrPad(song.Artist, colArtist),
			truncateOrPad(song.Album, colAlbum),
			song.Popularity,
		))
	}
	s.WriteString("\n[n] next, [p] previous, [q] back to menu\n")
	return s.String()
}

func truncateOrPad(s string, width int) string {
	if len(s) > width {
		if width > 3 {
			return s[:width-3] + "..."
		}
		return s[:width]
	}
	return fmt.Sprintf("%-*s", width, s)
}
