package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// DataProvider defines an interface for fetching data asynchronously.
type DataProvider interface {
	FetchTopArtists(timeRange string) tea.Cmd
	FetchTopSongs(timeRange string) tea.Cmd
	FetchArtistsPage(url string) tea.Cmd
	FetchSongsPage(url string) tea.Cmd
}

// DefaultDataProvider is the default implementation of DataProvider.
type DefaultDataProvider struct{}

// FetchTopArtists fetches the top artists for a given time range asynchronously.
func (p DefaultDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?time_range=%s", timeRange)
	return p.FetchArtistsPage(url)
}

// FetchTopSongs fetches the top songs for a given time range asynchronously.
func (p DefaultDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?time_range=%s", timeRange)
	return p.FetchSongsPage(url)
}

// FetchArtistsPage fetches an artists page asynchronously.
func (p DefaultDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchArtistsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

// FetchSongsPage fetches a songs page asynchronously.
func (p DefaultDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchSongsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}
