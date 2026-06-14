package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Using BubbleTea's Cmd pattern
type DataProvider interface {
	FetchTopArtists(timeRange string) tea.Cmd
	FetchTopSongs(timeRange string) tea.Cmd
	FetchArtistsPage(url string) tea.Cmd
	FetchSongsPage(url string) tea.Cmd
}

type SpotifyDataProvider struct {
	client SpotifyClient
}

func NewSpotifyDataProvider(client SpotifyClient) *SpotifyDataProvider {
	return &SpotifyDataProvider{
		client: client,
	}
}

func (p *SpotifyDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?time_range=%s", timeRange)
	return p.FetchArtistsPage(url)
}

func (p *SpotifyDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?time_range=%s", timeRange)
	return p.FetchSongsPage(url)
}

func (p *SpotifyDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchArtistsPage(p.client, url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (p *SpotifyDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchSongsPage(p.client, url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}
