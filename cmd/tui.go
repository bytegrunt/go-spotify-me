package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type APIResponse struct {
    Artists []Artist
    Next    string
    Prev    string
}

type model struct {
    artists []Artist
    next    string
    prev    string
    err     error
}

func initialModel(response APIResponse) model {
    return model{
        artists: response.Artists,
        next:    response.Next,
        prev:    response.Prev,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "esc":
            return m, tea.Quit
        case "n": // Fetch next page
            if m.next != "" {
                response, err := fetchArtistsPage(m.next)
                if err != nil {
                    m.err = err
                    return m, nil
                }
                m.artists = response.Artists
                m.next = response.Next
                m.prev = response.Prev
            }
        case "p": // Fetch previous page
            if m.prev != "" {
                response, err := fetchArtistsPage(m.prev)
                if err != nil {
                    m.err = err
                    return m, nil
                }
                m.artists = response.Artists
                m.next = response.Next
                m.prev = response.Prev
            }
        }
    }
    return m, nil
}

func (m model) View() string {
    if m.err != nil {
        return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
    }

    var s strings.Builder

    // Add table header
    s.WriteString(fmt.Sprintf("%-30s %-60s %-10s\n", "Name", "Genres", "Popularity"))
    s.WriteString(strings.Repeat("-", 100) + "\n") // Add a separator line

    // Add table rows
    for _, artist := range m.artists {
        s.WriteString(fmt.Sprintf("%-30s %-60s %-10d\n", artist.Name, artist.Genres, artist.Popularity))
    }

    // Add navigation instructions
    s.WriteString("\n[Press 'n' for next page, 'p' for previous page, 'q' to quit]\n")
    return s.String()
}
