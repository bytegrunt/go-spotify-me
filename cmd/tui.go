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
	songs       []Song // define Song struct below
	err         error
}

type APIResponse struct {
    Artists []Artist
    Next    string
    Prev    string
}

type switchToArtistsMsg struct {
	response APIResponse
}

type switchToSongsMsg struct {
	songs []Song
}

type errMsg struct {
	err error
}

func (m appModel) Init() tea.Cmd {
	return nil
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
				songs, err := fetchSongs()
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
			}

		case "p":
			if m.currentView == viewArtists && m.artists.Prev != "" {
				response, err := fetchArtistsPage(m.artists.Prev)
				if err != nil {
					return m, nil
				}
				m.artists = response
			}
		}

	case switchToArtistsMsg:
		m.artists = msg.response
		m.currentView = viewArtists
	case switchToSongsMsg:
		m.songs = msg.songs
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
	s.WriteString(fmt.Sprintf("%-30s %-60s %-10s\n", "Name", "Genres", "Popularity"))
	s.WriteString(strings.Repeat("-", 100) + "\n")
	for _, artist := range m.artists.Artists {
		s.WriteString(fmt.Sprintf("%-30s %-60s %-10d\n", artist.Name, artist.Genres, artist.Popularity))
	}
	s.WriteString("\n[n] next, [p] previous, [q] back to menu\n")
	return s.String()
}

func (m appModel) renderSongs() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%-30s %-30s %-30s %-10s\n", "Name", "Artist", "Album", "Popularity"))
	s.WriteString(strings.Repeat("-", 100) + "\n")
	for _, song := range m.songs {
		s.WriteString(fmt.Sprintf("%-30s %-30s %-30s %-10d\n", song.Name, song.Artist, song.Album, song.Popularity))
	}
	s.WriteString("\n[q] back to menu\n")
	return s.String()
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
