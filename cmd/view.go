package cmd

import (
	"fmt"
	"strings"

	"github.com/bytegrunt/go-spotify-me/internal/theme"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var footer = theme.HelpStyle.Render("[↑/↓] Navigate  [←] Prev Page  [→] Next Page  [1] Short  [2] Medium  [3] Long  [q] Back")

func (m appModel) View() string {
    if m.err != nil {
        return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
    }

    switch m.currentView {
    case viewMenu:
        return m.renderMenu()
    case viewArtists:
		return m.renderTable(m.artistTable) + "\n" + footer
	case viewSongs:
		return m.renderTable(m.songTable) + "\n" + footer
	case viewEnterClientID:
		return m.renderEnterClientID()
	default:
        return "Unknown view"
    }
}

func (m appModel) renderTable(t table.Model) string {
    var rows []string

    // Headers
    var headers []string
    for _, col := range t.Columns() {
        headers = append(headers, theme.HeaderStyle.Render(col.Title))
    }

    // Rows
    for i, row := range t.Rows() {
        rowStr := lipgloss.JoinHorizontal(lipgloss.Top, row...)
        style := theme.RowStyle
        if i == t.Cursor() {
            style = theme.SelectedRowStyle
        }
        rows = append(rows, style.Render(rowStr))
    }

    tableBody := lipgloss.JoinVertical(lipgloss.Left,
        lipgloss.JoinHorizontal(lipgloss.Top, headers...),
        strings.Join(rows, "\n"),
    )

    return theme.TableContainerStyle.Render(tableBody)
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
