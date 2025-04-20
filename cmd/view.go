package cmd

import (
	"fmt"
	"strings"

	"github.com/bytegrunt/go-spotify-me/internal/theme"
	"github.com/charmbracelet/bubbles/table"
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
		return m.renderTable(m.artistTable, m.artistColWidths) + "\n" + footer
	case viewSongs:
		return m.renderTable(m.songTable, m.songColWidths) + "\n" + footer
	case viewEnterClientID:
		return m.renderEnterClientID()
	default:
		return "Unknown view"
	}
}

func (m appModel) renderTable(t table.Model, colWidths []int) string {
	var rows []string

	// Header
	headers := make([]string, len(t.Columns()))
	for i, col := range t.Columns() {
		headers[i] = col.Title
	}
	header := theme.RenderRow(headers, colWidths, theme.HeaderStyle)

	// Body
	for i, row := range t.Rows() {
		style := theme.RowStyle
		if i == t.Cursor() {
			style = theme.SelectedRowStyle
		}
		rows = append(rows, theme.RenderRow(row, colWidths, style))
	}

	body := header + "\n" + strings.Join(rows, "\n")
	return theme.TableContainerStyle.Render(body)
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
