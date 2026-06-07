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
	FetchMe() (Me, error)
	Login() error
}

type DefaultDataProvider struct{}

func (d *DefaultDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	return d.FetchArtistsPage(fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?time_range=%s", timeRange))
}

func (d *DefaultDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	return d.FetchSongsPage(fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?time_range=%s", timeRange))
}

func (d *DefaultDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchArtistsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (d *DefaultDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchSongsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}

func (d *DefaultDataProvider) FetchMe() (Me, error) {
	return fetchMe()
}

func (d *DefaultDataProvider) Login() error {
	return Login()
}
