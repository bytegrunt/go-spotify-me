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
			// Return to the menu and blur the table
			switch m.currentView {
			case viewArtists:
				m.artistTable.Blur()
			case viewSongs:
				m.songTable.Blur()
			}
			if m.currentView != viewMenu {
				m.currentView = viewMenu
				return m, nil
			}
			return m, tea.Quit

		case "a", "A":
			// Only switch to the Artists view if in the main menu
			if m.currentView == viewMenu {
				m.artistTable.Focus()
				return m, func() tea.Msg {
					response, err := fetchArtistsPage("https://api.spotify.com/v1/me/top/artists?time_range=medium_term")
					if err != nil {
						return errMsg{err}
					}
					return switchToArtistsMsg{response}
				}
			}

		case "s", "S":
			// Only switch to the Songs view if in the main menu
			if m.currentView == viewMenu {
				m.songTable.Focus()
				return m, func() tea.Msg {
					response, err := fetchSongsPage("https://api.spotify.com/v1/me/top/tracks?time_range=medium_term")
					if err != nil {
						return errMsg{err}
					}
					return switchToSongsMsg{response}
				}
			}

		case "1":
			// fetch short term artists or songs
			switch m.currentView {
			case viewArtists:
				m.artistTable.Focus()
				return m, func() tea.Msg {
					response, err := fetchArtistsPage("https://api.spotify.com/v1/me/top/artists?time_range=short_term")
					if err != nil {
						return errMsg{err}
					}
					return switchToArtistsMsg{response}
				}
			case viewSongs:
				m.songTable.Focus()
				return m, func() tea.Msg {
					response, err := fetchSongsPage("https://api.spotify.com/v1/me/top/tracks?time_range=short_term")
					if err != nil {
						return errMsg{err}
					}
					return switchToSongsMsg{response}
				}
			}
		case "2":
			// medium
			switch m.currentView {
			case viewSongs:
				m.songTable.Focus()
				return m, func() tea.Msg {
					response, err := fetchSongsPage("https://api.spotify.com/v1/me/top/tracks?time_range=medium_term")
					if err != nil {
						return errMsg{err}
					}
					return switchToSongsMsg{response}
				}

			case viewArtists:
				m.artistTable.Focus()
				return m, func() tea.Msg {
					response, err := fetchArtistsPage("https://api.spotify.com/v1/me/top/artists?time_range=medium_term")
					if err != nil {
						return errMsg{err}
					}
					return switchToArtistsMsg{response}
				}
			}
		case "3":
			// long
			switch m.currentView {
			case viewSongs:
				m.songTable.Focus()
				return m, func() tea.Msg {
					response, err := fetchSongsPage("https://api.spotify.com/v1/me/top/tracks?time_range=long_term")
					if err != nil {
						return errMsg{err}
					}
					return switchToSongsMsg{response}
				}
			case viewArtists:
				m.artistTable.Focus()
				return m, func() tea.Msg {
					response, err := fetchArtistsPage("https://api.spotify.com/v1/me/top/artists?time_range=long_term")
					if err != nil {
						return errMsg{err}
					}
					return switchToArtistsMsg{response}
				}
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

				// Fetch the user's information
				me, err := fetchMe()
				if err != nil {
					m.err = fmt.Errorf("failed to fetch user info: %v", err)
					return m, nil
				}
				m.me = me

				// Initialize artist table
				m.artistTable = table.New(
					table.WithColumns([]table.Column{
						{Title: "Name", Width: 40},
						{Title: "Genres", Width: 50},
						{Title: "Popularity", Width: 10},
					}),
					table.WithFocused(false),
				)

				// Initialize song table
				m.songTable = table.New(
					table.WithColumns([]table.Column{
						{Title: "Name", Width: 40},
						{Title: "Artist", Width: 20},
						{Title: "Album", Width: 30},
						{Title: "Popularity", Width: 10},
					}),
					table.WithFocused(false),
				)

				// Switch to the menu view after successful login
				m.currentView = viewMenu
				return m, nil
			}
		}

		// Delegate key events to the focused table
		switch m.currentView {
		case viewArtists:
			m.artistTable, cmd = m.artistTable.Update(msg)
		case viewSongs:
			m.songTable, cmd = m.songTable.Update(msg)
		}

	case tea.WindowSizeMsg:
		m.windowSize = msg

		// Recalculate column widths
		m.artistColWidths = calculateColumnWidths(msg.Width, []float64{0.4, 0.4, 0.2})
		m.songColWidths = calculateColumnWidths(msg.Width, []float64{0.4, 0.2, 0.2, 0.2})

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
