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
	rows := [][]string{
		{"Name", m.me.DisplayName},
		{"Email", m.me.Email},
		{"Product", m.me.Product},
		{"Country", m.me.Country},
		{"Profile URL", m.me.ProfileURL},
	}

	if m.windowSize.Width < 20 {
		m.windowSize.Width = 20
	}

	colWidths := calculateColumnWidths(m.windowSize.Width, []float64{0.3, 0.7})

	var renderedRows []string
	header := theme.RenderRow([]string{"Field", "Value"}, colWidths, theme.HeaderStyle)
	for _, row := range rows {
		renderedRows = append(renderedRows, theme.RenderRow(row, colWidths, theme.RowStyle))
	}

	table := header + "\n" + strings.Join(renderedRows, "\n")
	return theme.TableContainerStyle.Render(table) + "\n" + theme.HelpStyle.Render("[A] Top Artists  [S] Top Songs  [Q] Quit")
}
