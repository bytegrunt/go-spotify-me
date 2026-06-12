package cmd

type DataProvider interface {
	FetchArtists(url string) (APIResponse, error)
	FetchSongs(url string) (APIResponse, error)
	FetchMe() (Me, error)
}

type DefaultDataProvider struct{}

func NewDataProvider() *DefaultDataProvider {
	return &DefaultDataProvider{}
}

func (d *DefaultDataProvider) FetchArtists(url string) (APIResponse, error) {
	return fetchArtistsPage(url)
}

func (d *DefaultDataProvider) FetchSongs(url string) (APIResponse, error) {
	return fetchSongsPage(url)
}

func (d *DefaultDataProvider) FetchMe() (Me, error) {
	return fetchMe()
}
