package cmd

import tea "github.com/charmbracelet/bubbletea"

// DataProvider defines an interface for fetching data asynchronously
type DataProvider interface {
	FetchTopArtists(timeRange string) tea.Cmd
	FetchTopSongs(timeRange string) tea.Cmd
	FetchArtistsPage(url string) tea.Cmd
	FetchSongsPage(url string) tea.Cmd
}

// DefaultDataProvider implements DataProvider by calling actual API fetch functions
type DefaultDataProvider struct{}

func (d DefaultDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	return func() tea.Msg {
		url := "https://api.spotify.com/v1/me/top/artists?time_range=" + timeRange
		response, err := fetchArtistsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (d DefaultDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	return func() tea.Msg {
		url := "https://api.spotify.com/v1/me/top/tracks?time_range=" + timeRange
		response, err := fetchSongsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}

func (d DefaultDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchArtistsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (d DefaultDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchSongsPage(url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}
