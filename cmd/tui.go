package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type viewType int

const (
	viewMenu viewType = iota
	viewArtists
	viewSongs
)

type appModel struct {
	currentView viewType
	artists     APIResponse
	songs       APIResponse
	windowSize  tea.WindowSizeMsg
	err         error
}

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

func (m appModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.ClearScreen,
		tea.WindowSize(),
	)
}

func initialAppModel() appModel {
	return appModel{
		currentView: viewMenu,
	}
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, func() tea.Msg {
				response, err := fetchArtistsPage("https://api.spotify.com/v1/me/top/artists")
				if err != nil {
					return errMsg{err}
				}
				return switchToArtistsMsg{response}
			}

		case "s", "S":
			return m, func() tea.Msg {
				songs, err := fetchSongsPage("https://api.spotify.com/v1/me/top/tracks")
				if err != nil {
					return errMsg{err}
				}
				return switchToSongsMsg{songs}
			}

		case "n":
			if m.currentView == viewArtists && m.artists.Next != "" {
				response, err := fetchArtistsPage(m.artists.Next)
				if err != nil {
					return m, nil
				}
				m.artists = response
			} else if m.currentView == viewSongs && m.songs.Next != "" {
				response, err := fetchSongsPage(m.songs.Next)
				if err != nil {
					return m, nil
				}
				m.songs = response
			}
		
		case "p":
			if m.currentView == viewArtists && m.artists.Prev != "" {
				response, err := fetchArtistsPage(m.artists.Prev)
				if err != nil {
					return m, nil
				}
				m.artists = response
			} else if m.currentView == viewSongs && m.songs.Prev != "" {
				response, err := fetchSongsPage(m.songs.Prev)
				if err != nil {
					return m, nil
				}
				m.songs = response
			}
		
		}

	// Handle window resize
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
	return m, nil
}

func (m appModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
	}

	switch m.currentView {
	case viewMenu:
		return "\nWelcome to go-spotify\n\nPress A for Top Artists\nPress S for Top Songs\nPress Q to quit"
	case viewArtists:
		return m.renderArtists()
	case viewSongs:
		return m.renderSongs()
	default:
		return "Unknown view"
	}
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


// func initialRootModel() rootModel {
// 	return rootModel{
// 		message: "Welcome to go-spotify\n\nPress A for Top Artists\nPress S for Top Songs\nPress Q to quit",
// 	}
// }

// func (m rootModel) Init() tea.Cmd {
// 	return nil
// }

// func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "a", "A":
// 			return m, func() tea.Msg {
// 				// Fetch and return the top artists
// 				response, err := fetchArtistsPage("https://api.spotify.com/v1/me/top/artists")
// 				if err != nil {
// 					return errMsg{err}
// 				}
// 				return switchToTopArtistsMsg{response}
// 			}
// 		case "s", "S":
// 			// You would create a `fetchTopSongsPage` function for this
// 			return m, func() tea.Msg {
// 				// Placeholder - you'd fetch and return the top songs here
// 				return errMsg{fmt.Errorf("Top songs not yet implemented")}
// 			}
// 		case "q", "Q", "esc":
// 			return m, tea.Quit
// 		}
// 	case switchToTopArtistsMsg:
// 		return initialModel(msg.response), nil
// 	case errMsg:
// 		return m, func() tea.Msg {
// 			fmt.Println("Error:", msg.err)
// 			return tea.Quit()
// 		}
// 	}
// 	return m, nil
// }

// func (m rootModel) View() string {
// 	return m.message
// }

// func initialModel(response APIResponse) model {
//     return model{
//         artists: response.Artists,
//         next:    response.Next,
//         prev:    response.Prev,
//     }
// }

// func (m model) Init() tea.Cmd {
//     return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//     switch msg := msg.(type) {
//     case tea.KeyMsg:
//         switch msg.String() {
//         case "q", "esc":
//             return m, tea.Quit
//         case "n": // Fetch next page
//             if m.next != "" {
//                 response, err := fetchArtistsPage(m.next)
//                 if err != nil {
//                     m.err = err
//                     return m, nil
//                 }
//                 m.artists = response.Artists
//                 m.next = response.Next
//                 m.prev = response.Prev
//             }
//         case "p": // Fetch previous page
//             if m.prev != "" {
//                 response, err := fetchArtistsPage(m.prev)
//                 if err != nil {
//                     m.err = err
//                     return m, nil
//                 }
//                 m.artists = response.Artists
//                 m.next = response.Next
//                 m.prev = response.Prev
//             }
//         }
//     }
//     return m, nil
// }

// func (m model) View() string {
//     if m.err != nil {
//         return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
//     }

//     var s strings.Builder

//     // Add table header
//     s.WriteString(fmt.Sprintf("%-30s %-60s %-10s\n", "Name", "Genres", "Popularity"))
//     s.WriteString(strings.Repeat("-", 100) + "\n") // Add a separator line

//     // Add table rows
//     for _, artist := range m.artists {
//         s.WriteString(fmt.Sprintf("%-30s %-60s %-10d\n", artist.Name, artist.Genres, artist.Popularity))
//     }

//     // Add navigation instructions
//     s.WriteString("\n[Press 'n' for next page, 'p' for previous page, 'q' to quit]\n")
//     return s.String()
// }
