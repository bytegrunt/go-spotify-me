package cmd

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type FakeDataProvider struct{}

func (f FakeDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	return func() tea.Msg {
		return switchToArtistsMsg{response: APIResponse{}}
	}
}

func (f FakeDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	return func() tea.Msg {
		return switchToSongsMsg{response: APIResponse{}}
	}
}

func (f FakeDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg {
		return switchToArtistsMsg{response: APIResponse{}}
	}
}

func (f FakeDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg {
		return switchToSongsMsg{response: APIResponse{}}
	}
}

func TestAppModelUpdateWithFakeProvider_Artists(t *testing.T) {
	dp := FakeDataProvider{}
	m := appModel{
		currentView:  viewMenu,
		dataProvider: dp,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, cmd := m.Update(msg)

	teaMsg := cmd()
	if _, ok := teaMsg.(switchToArtistsMsg); !ok {
		t.Errorf("Expected switchToArtistsMsg, got %T", teaMsg)
	}

	finalModel, _ := newModel.Update(teaMsg)

	app, ok := finalModel.(appModel)
	if !ok {
		t.Fatalf("Expected appModel, got %T", finalModel)
	}
	if app.currentView != viewArtists {
		t.Errorf("Expected currentView to be viewArtists, got %v", app.currentView)
	}
}

func TestAppModelUpdateWithFakeProvider_Songs(t *testing.T) {
	dp := FakeDataProvider{}
	m := appModel{
		currentView:  viewMenu,
		dataProvider: dp,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newModel, cmd := m.Update(msg)

	teaMsg := cmd()
	if _, ok := teaMsg.(switchToSongsMsg); !ok {
		t.Errorf("Expected switchToSongsMsg, got %T", teaMsg)
	}

	finalModel, _ := newModel.Update(teaMsg)

	app, ok := finalModel.(appModel)
	if !ok {
		t.Fatalf("Expected appModel, got %T", finalModel)
	}
	if app.currentView != viewSongs {
		t.Errorf("Expected currentView to be viewSongs, got %v", app.currentView)
	}
}
