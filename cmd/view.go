package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
    headerStyle = lipgloss.NewStyle().
            Bold(true).
            Foreground(lipgloss.Color("5")).
            Background(lipgloss.Color("236")).
            Padding(0, 1)

    rowStyle = lipgloss.NewStyle().
            Foreground(lipgloss.Color("7")).
            Background(lipgloss.Color("0")).
            Padding(0, 1)

    selectedRowStyle = lipgloss.NewStyle().
            Foreground(lipgloss.Color("0")).
            Background(lipgloss.Color("5")).
            Padding(0, 1)
)

func (m appModel) View() string {
    if m.err != nil {
        return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
    }

    switch m.currentView {
    case viewMenu:
        return m.renderMenu()
    case viewArtists:
        return m.artistTable.View() + "\n\n[↑/↓] Navigate | [←] Previous Page | [→] Next Page | [q] Back to Menu"
    case viewSongs:
        return m.songTable.View() + "\n\n[↑/↓] Navigate | [←] Previous Page | [→] Next Page | [q] Back to Menu"
    default:
        return "Unknown view"
    }
}

func (m appModel) renderTable(t table.Model) string {
    var rows []string
    for i, row := range t.Rows() {
        style := rowStyle
        if i == t.Cursor() {
            style = selectedRowStyle
        }
        rows = append(rows, style.Render(fmt.Sprintf("%s | %s | %s", row[0], row[1], row[2])))
    }

    headers := headerStyle.Render(fmt.Sprintf("%s | %s | %s", t.Columns()[0].Title, t.Columns()[1].Title, t.Columns()[2].Title))
    return headers + "\n" + strings.Join(rows, "\n")
}

func (m appModel) renderEnterClientID() string {
	return fmt.Sprintf(
		"Enter your Spotify Client ID:\n\n%s\n\nPress Enter to confirm, or Esc to quit.",
		m.textInput.View(),
	)
}

func (m appModel) renderMenu() string {
	return fmt.Sprintf(
		"Welcome, %s (%s)\nProduct: %s, Country:%s\n\n"+
			"Menu:\n"+
			"Press A for Top Artists\n"+
			"Press S for Top Songs\n"+
			"Press Q to quit\n\n"+
			"My Spotify Profile: %s\n",
		m.me.DisplayName,
		m.me.Email,
		m.me.Product,
		m.me.Country,
		m.me.ProfileURL,
	)
}
