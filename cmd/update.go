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

type initLoginMsg struct {
	me  Me
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
				return m, m.dataProvider.FetchTopArtists("medium_term")
			}

		case "s", "S":
			// Only switch to the Songs view if in the main menu
			if m.currentView == viewMenu {
				m.songTable.Focus()
				return m, m.dataProvider.FetchTopSongs("medium_term")
			}

		case "1":
			// fetch short term artists or songs
			switch m.currentView {
			case viewArtists:
				m.artistTable.Focus()
				return m, m.dataProvider.FetchTopArtists("short_term")
			case viewSongs:
				m.songTable.Focus()
				return m, m.dataProvider.FetchTopSongs("short_term")
			}
		case "2":
			// medium
			switch m.currentView {
			case viewSongs:
				m.songTable.Focus()
				return m, m.dataProvider.FetchTopSongs("medium_term")

			case viewArtists:
				m.artistTable.Focus()
				return m, m.dataProvider.FetchTopArtists("medium_term")
			}
		case "3":
			// long
			switch m.currentView {
			case viewSongs:
				m.songTable.Focus()
				return m, m.dataProvider.FetchTopSongs("long_term")
			case viewArtists:
				m.artistTable.Focus()
				return m, m.dataProvider.FetchTopArtists("long_term")
			}

		case "right": // Handle next page for Artists or Songs
			if m.currentView == viewArtists && m.artists.Next != "" {
				return m, m.dataProvider.FetchArtistsPage(m.artists.Next)
			} else if m.currentView == viewSongs && m.songs.Next != "" {
				return m, m.dataProvider.FetchSongsPage(m.songs.Next)
			}

		case "left": // Handle previous page for Artists or Songs
			if m.currentView == viewArtists && m.artists.Prev != "" {
				return m, m.dataProvider.FetchArtistsPage(m.artists.Prev)
			} else if m.currentView == viewSongs && m.songs.Prev != "" {
				return m, m.dataProvider.FetchSongsPage(m.songs.Prev)
			}

		case "enter":
			// Handle entering the Client ID
			if m.currentView == viewEnterClientID {
				m.clientID = m.textInput.Value()
				err := keyring.Set("go-spotify-me-cli", "client_id", m.clientID)
				if err != nil {
					m.err = fmt.Errorf("failed to store client ID in keyring: %w", err)
					return m, nil
				}

				// Switch to the menu view immediately and let it populate data via commands
				m.currentView = viewMenu

				// Initialize artist table
				m.artistColWidths = calculateColumnWidths(m.windowSize.Width, []float64{0.4, 0.4, 0.2})
				m.artistTable = table.New(
					table.WithColumns([]table.Column{
						{Title: "Name", Width: 40},
						{Title: "Genres", Width: 50},
						{Title: "Popularity", Width: 10},
					}),
					table.WithFocused(false),
				)

				// Initialize song table
				m.songColWidths = calculateColumnWidths(m.windowSize.Width, []float64{0.4, 0.3, 0.2, 0.1})
				m.songTable = table.New(
					table.WithColumns([]table.Column{
						{Title: "Name", Width: 40},
						{Title: "Artist", Width: 20},
						{Title: "Album", Width: 30},
						{Title: "Popularity", Width: 10},
					}),
					table.WithFocused(false),
				)

				return m, func() tea.Msg {
					err := m.dataProvider.Login()
					if err != nil {
						return initLoginMsg{err: fmt.Errorf("failed to log in: %w", err)}
					}
					me, err := m.dataProvider.FetchMe()
					if err != nil {
						return initLoginMsg{err: fmt.Errorf("failed to fetch user info: %w", err)}
					}
					return initLoginMsg{me: me}
				}
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
		m.songColWidths = calculateColumnWidths(msg.Width, []float64{0.4, 0.3, 0.2, 0.1})

	case initLoginMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.me = msg.me

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
