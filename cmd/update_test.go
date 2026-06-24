package cmd

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type FakeDataProvider struct {
	CalledFetchTopArtists bool
	CalledFetchTopSongs   bool
	TimeRange             string
}

func (f *FakeDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	f.CalledFetchTopArtists = true
	f.TimeRange = timeRange
	return func() tea.Msg {
		return switchToArtistsMsg{APIResponse{Artists: []Artist{{Name: "Fake Artist"}}}}
	}
}

func (f *FakeDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	f.CalledFetchTopSongs = true
	f.TimeRange = timeRange
	return func() tea.Msg {
		return switchToSongsMsg{APIResponse{Songs: []Song{{Name: "Fake Song"}}}}
	}
}

func (f *FakeDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return nil
}

func (f *FakeDataProvider) FetchSongsPage(url string) tea.Cmd {
	return nil
}

func TestUpdate_FetchTopArtists(t *testing.T) {
	fakeProvider := &FakeDataProvider{}
	model := appModel{
		currentView:  viewMenu,
		dataProvider: fakeProvider,
		artistTable:  table.New(table.WithColumns([]table.Column{{Title: "Name"}, {Title: "Genres"}, {Title: "Popularity"}})),
	}

	// simulate 'a' keypress
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, cmd := model.Update(msg)

	if !fakeProvider.CalledFetchTopArtists {
		t.Errorf("Expected FetchTopArtists to be called")
	}

	if fakeProvider.TimeRange != "medium_term" {
		t.Errorf("Expected medium_term time_range, got %s", fakeProvider.TimeRange)
	}

	if cmd == nil {
		t.Fatalf("Expected cmd to be returned")
	}

	resMsg := cmd()
	_, ok := resMsg.(switchToArtistsMsg)
	if !ok {
		t.Errorf("Expected cmd to return switchToArtistsMsg")
	}

	// Update with switchToArtistsMsg to ensure the model handles it correctly
	newModel, _ = newModel.Update(resMsg)
	updatedModel := newModel.(appModel)
	if updatedModel.currentView != viewArtists {
		t.Errorf("Expected currentView to be viewArtists, got %v", updatedModel.currentView)
	}

	if len(updatedModel.artists.Artists) == 0 || updatedModel.artists.Artists[0].Name != "Fake Artist" {
		t.Errorf("Expected model to contain fake artist")
	}
}

func TestUpdate_FetchTopSongs(t *testing.T) {
	fakeProvider := &FakeDataProvider{}
	model := appModel{
		currentView:  viewMenu,
		dataProvider: fakeProvider,
		songTable:    table.New(table.WithColumns([]table.Column{{Title: "Name"}, {Title: "Artist"}, {Title: "Album"}, {Title: "Popularity"}})),
	}

	// simulate 's' keypress
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newModel, cmd := model.Update(msg)

	if !fakeProvider.CalledFetchTopSongs {
		t.Errorf("Expected FetchTopSongs to be called")
	}

	if fakeProvider.TimeRange != "medium_term" {
		t.Errorf("Expected medium_term time_range, got %s", fakeProvider.TimeRange)
	}

	if cmd == nil {
		t.Fatalf("Expected cmd to be returned")
	}

	resMsg := cmd()
	_, ok := resMsg.(switchToSongsMsg)
	if !ok {
		t.Errorf("Expected cmd to return switchToSongsMsg")
	}

	newModel, _ = newModel.Update(resMsg)
	updatedModel := newModel.(appModel)
	if updatedModel.currentView != viewSongs {
		t.Errorf("Expected currentView to be viewSongs, got %v", updatedModel.currentView)
	}

	if len(updatedModel.songs.Songs) == 0 || updatedModel.songs.Songs[0].Name != "Fake Song" {
		t.Errorf("Expected model to contain fake song")
	}
}
