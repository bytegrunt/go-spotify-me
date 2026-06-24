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

type SpotifyDataProvider struct{}

func NewSpotifyDataProvider() *SpotifyDataProvider {
	return &SpotifyDataProvider{}
}

func (s *SpotifyDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?time_range=%s", timeRange)
		response, err := fetchArtistsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (s *SpotifyDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?time_range=%s", timeRange)
		response, err := fetchSongsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}

func (s *SpotifyDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchArtistsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (s *SpotifyDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchSongsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}
