package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zalando/go-keyring"
)

type APIResponse struct {
	Artists []Artist
	Songs   []Song
	Next    string
	Prev    string
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

		case "a", "A":
			// Switch to the Artists view
			return m, func() tea.Msg {
				response, err := fetchArtistsPage("https://api.spotify.com/v1/me/top/artists?time_range=medium_term")
				if err != nil {
					return errMsg{err}
				}
				return switchToArtistsMsg{response}
			}

		case "s", "S":
			// Switch to the Songs view
			return m, func() tea.Msg {
				response, err := fetchSongsPage("https://api.spotify.com/v1/me/top/tracks?time_range=medium_term")
				if err != nil {
					return errMsg{err}
				}
				return switchToSongsMsg{response}
			}

		case "right": // Handle next page for Artists or Songs
			if m.currentView == viewArtists && m.artists.Next != "" {
				return m, func() tea.Msg {
					response, err := fetchArtistsPage(m.artists.Next)
					if err != nil {
						return errMsg{err}
					}
					return switchToArtistsMsg{response}
				}
			} else if m.currentView == viewSongs && m.songs.Next != "" {
				return m, func() tea.Msg {
					response, err := fetchSongsPage(m.songs.Next)
					if err != nil {
						return errMsg{err}
					}
					return switchToSongsMsg{response}
				}
			}

		case "left": // Handle previous page for Artists or Songs
			if m.currentView == viewArtists && m.artists.Prev != "" {
				return m, func() tea.Msg {
					response, err := fetchArtistsPage(m.artists.Prev)
					if err != nil {
						return errMsg{err}
					}
					return switchToArtistsMsg{response}
				}
			} else if m.currentView == viewSongs && m.songs.Prev != "" {
				return m, func() tea.Msg {
					response, err := fetchSongsPage(m.songs.Prev)
					if err != nil {
						return errMsg{err}
					}
					return switchToSongsMsg{response}
				}
			}

		case "enter":
			// Handle entering the Client ID
			if m.currentView == viewEnterClientID {
				m.clientID = m.textInput.Value()
				err := keyring.Set("go-spotify-me-cli", "client_id", m.clientID)
				if err != nil {
					m.err = fmt.Errorf("failed to store client ID in keyring: %v", err)
					return m, nil
				}

				// Call the Login function after saving the Client ID
				err = Login()
				if err != nil {
					m.err = fmt.Errorf("failed to log in: %v", err)
					return m, nil
				}

				// Switch to the menu view after successful login
				m.currentView = viewMenu
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.windowSize = msg

	case switchToArtistsMsg:
		m.artists = msg.response
		rows := []table.Row{}
		for _, artist := range m.artists.Artists {
			rows = append(rows, table.Row{artist.Name, artist.Genres, fmt.Sprintf("%d", artist.Popularity)})
		}
		m.artistTable.SetRows(rows)
		m.currentView = viewArtists
	
	case switchToSongsMsg:
		m.songs = msg.response
		rows := []table.Row{}
		for _, song := range m.songs.Songs {
			rows = append(rows, table.Row{song.Name, song.Artist, song.Album, fmt.Sprintf("%d", song.Popularity)})
		}
		m.songTable.SetRows(rows)
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
