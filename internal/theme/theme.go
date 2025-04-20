package theme

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Define your theme colors
var (
	colorPrimary    = lipgloss.Color("#1DB954") // Spotify green
	colorBackground = lipgloss.Color("#1E1E1E") // Dark background
	colorForeground = lipgloss.Color("#FFFFFF")
	colorMuted      = lipgloss.Color("#888888")
	colorAccent     = lipgloss.Color("#333333")
	colorHighlight  = lipgloss.Color("#3E3E3E")
	colorSelectedBG = lipgloss.Color("#2D46B9") // Spotify blue
	colorSelectedFG = lipgloss.Color("#FFFFFF")
)

// Header style for table columns
var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(colorPrimary).
	Background(colorAccent).
	MarginBottom(1).
	Padding(0, 1)

// Style for regular rows
var RowStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Background(colorBackground).
	Padding(0, 1)

// Style for selected row
var SelectedRowStyle = lipgloss.NewStyle().
	Foreground(colorSelectedFG).
	Background(colorSelectedBG).
	Bold(true).
	Padding(0, 1)

// Style to wrap the entire table
var TableContainerStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorPrimary).
	Margin(1, 2).
	Padding(1, 2)

// Style for the help/footer text
var HelpStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	MarginTop(1)

func RenderRow(cells []string, widths []int, style lipgloss.Style) string {
	rendered := make([]string, len(cells))
	for i, cell := range cells {
		rendered[i] = truncateOrPad(cell, widths[i])
	}
	return style.Render(strings.Join(rendered, " | "))
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
