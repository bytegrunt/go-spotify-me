package cmd

func fetchArtistsPage(url string) (APIResponse, error) {
	artists, next, prev, err := spotifyClient.GetArtistsPage(url)
	if err != nil {
		return APIResponse{}, err
	}
	return APIResponse{
		Artists: artists,
		Next:    next,
		Prev:    prev,
	}, nil
}
