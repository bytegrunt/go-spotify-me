package cmd

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type MockDataProvider struct {
	LastMethodCalled string
	LastTimeRange    string
	LastURL          string
}

func (m *MockDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	m.LastMethodCalled = "FetchTopArtists"
	m.LastTimeRange = timeRange
	return func() tea.Msg { return nil }
}

func (m *MockDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	m.LastMethodCalled = "FetchTopSongs"
	m.LastTimeRange = timeRange
	return func() tea.Msg { return nil }
}

func (m *MockDataProvider) FetchArtistsPage(url string) tea.Cmd {
	m.LastMethodCalled = "FetchArtistsPage"
	m.LastURL = url
	return func() tea.Msg { return nil }
}

func (m *MockDataProvider) FetchSongsPage(url string) tea.Cmd {
	m.LastMethodCalled = "FetchSongsPage"
	m.LastURL = url
	return func() tea.Msg { return nil }
}

func TestUpdateWithMockProvider(t *testing.T) {
	mockProvider := &MockDataProvider{}
	model := appModel{
		currentView: viewMenu,
		artistTable: table.New(),
		songTable:   table.New(),
		provider:    mockProvider,
	}

	// Test "a"
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if mockProvider.LastMethodCalled != "FetchTopArtists" || mockProvider.LastTimeRange != "medium_term" {
		t.Errorf("Expected FetchTopArtists with medium_term, got %s with %s", mockProvider.LastMethodCalled, mockProvider.LastTimeRange)
	}

	// Test "s"
	mockProvider.LastMethodCalled = ""
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if mockProvider.LastMethodCalled != "FetchTopSongs" || mockProvider.LastTimeRange != "medium_term" {
		t.Errorf("Expected FetchTopSongs with medium_term, got %s with %s", mockProvider.LastMethodCalled, mockProvider.LastTimeRange)
	}

	// Change view to viewArtists and test "1"
	model.currentView = viewArtists
	mockProvider.LastMethodCalled = ""
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	if mockProvider.LastMethodCalled != "FetchTopArtists" || mockProvider.LastTimeRange != "short_term" {
		t.Errorf("Expected FetchTopArtists with short_term, got %s with %s", mockProvider.LastMethodCalled, mockProvider.LastTimeRange)
	}

	// Test "right" for next artists page
	model.artists.Next = "http://next-artists"
	mockProvider.LastMethodCalled = ""
	model.Update(tea.KeyMsg{Type: tea.KeyRight})
	if mockProvider.LastMethodCalled != "FetchArtistsPage" || mockProvider.LastURL != "http://next-artists" {
		t.Errorf("Expected FetchArtistsPage with url, got %s with %s", mockProvider.LastMethodCalled, mockProvider.LastURL)
	}
}
