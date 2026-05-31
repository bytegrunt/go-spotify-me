package cmd

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type MockDataProvider struct {
	Msg tea.Msg
}

func (m MockDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	return func() tea.Msg { return m.Msg }
}

func (m MockDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	return func() tea.Msg { return m.Msg }
}

func (m MockDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg { return m.Msg }
}

func (m MockDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg { return m.Msg }
}

func TestAppModelUpdate_FetchArtists(t *testing.T) {
	mockMsg := switchToArtistsMsg{
		response: APIResponse{
			Artists: []Artist{{Name: "Mock Artist"}},
		},
	}

	provider := MockDataProvider{Msg: mockMsg}

	model := appModel{
		dataProvider: provider,
		currentView:  viewMenu,
		artistTable:  table.New(),
	}

	// Simulate pressing 'a'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Fatal("Expected a command to be returned, got nil")
	}

	cmdMsg := cmd()
	if resultMsg, ok := cmdMsg.(switchToArtistsMsg); !ok {
		t.Errorf("Expected switchToArtistsMsg, got %T", cmdMsg)
	} else {
		if len(resultMsg.response.Artists) != 1 || resultMsg.response.Artists[0].Name != "Mock Artist" {
			t.Errorf("Expected mock artist data, got %+v", resultMsg.response)
		}
	}
}

func TestAppModelUpdate_FetchSongs(t *testing.T) {
	mockMsg := switchToSongsMsg{
		response: APIResponse{
			Songs: []Song{{Name: "Mock Song"}},
		},
	}

	provider := MockDataProvider{Msg: mockMsg}

	model := appModel{
		dataProvider: provider,
		currentView:  viewMenu,
		songTable:    table.New(),
	}

	// Simulate pressing 's'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Fatal("Expected a command to be returned, got nil")
	}

	cmdMsg := cmd()
	if resultMsg, ok := cmdMsg.(switchToSongsMsg); !ok {
		t.Errorf("Expected switchToSongsMsg, got %T", cmdMsg)
	} else {
		if len(resultMsg.response.Songs) != 1 || resultMsg.response.Songs[0].Name != "Mock Song" {
			t.Errorf("Expected mock song data, got %+v", resultMsg.response)
		}
	}
}
