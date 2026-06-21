package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type DataProvider interface {
	FetchTopArtists(timeRange string) tea.Cmd
	FetchTopSongs(timeRange string) tea.Cmd
	FetchArtistsPage(url string) tea.Cmd
	FetchSongsPage(url string) tea.Cmd
}

type DefaultDataProvider struct{}

func NewDataProvider() DefaultDataProvider {
	return DefaultDataProvider{}
}

func (p DefaultDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	return p.FetchArtistsPage(fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?time_range=%s", timeRange))
}

func (p DefaultDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	return p.FetchSongsPage(fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?time_range=%s", timeRange))
}

func (p DefaultDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchArtistsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (p DefaultDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchSongsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}
