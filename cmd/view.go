package cmd

import (
	"fmt"
)

func (m appModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
	}

	switch m.currentView {
	case viewMenu:
		return m.renderMenu()
	case viewArtists:
        return m.artistTable.View() + "\n\n[←] Previous Page | [→] Next Page | [q] Back to Menu"
    case viewSongs:
        return m.songTable.View() + "\n\n[←] Previous Page | [→] Next Page | [q] Back to Menu"
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
		m.me.Country,
		m.me.ProfileURL,
	)
}
