package cmd

import (
	"fmt"
	"strings"
)

func (m appModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
	}

	switch m.currentView {
	case viewMenu:
		return m.renderMenu()
	case viewArtists:
		return m.renderArtists()
	case viewSongs:
		return m.renderSongs()
	case viewEnterClientID:
		return m.renderEnterClientID()
	default:
		return "Unknown view"
	}
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
		m.me.ProfileURL,
		m.me.Country,
	)
}

func (m appModel) renderArtists() string {
	var s strings.Builder

	// Dynamic column widths
	totalWidth := m.windowSize.Width
	if totalWidth == 0 {
		totalWidth = 100 // fallback
	}
	colName := 30
	colGenre := totalWidth - colName - 12 - 5 // leave space for padding and popularity
	colPop := 5

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
	s.WriteString("\n[←] previous, [→] next, [q] back to menu\n")
	return s.String()
}

func (m appModel) renderSongs() string {
	var s strings.Builder

	totalWidth := m.windowSize.Width
	if totalWidth == 0 {
		totalWidth = 100
	}
	colName := 40
	colArtist := 30
	colAlbum := totalWidth - colName - colArtist - 12 - 5
	colPop := 5

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
	s.WriteString("\n[←] previous, [→] next, [q] back to menu\n")
	return s.String()
}
