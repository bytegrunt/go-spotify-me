package cmd

func fetchSongsPage(url string) (APIResponse, error) {
	songs, next, prev, err := spotifyClient.GetSongsPage(url)
	if err != nil {
		return APIResponse{}, err
	}
	return APIResponse{
		Songs: songs,
		Next:  next,
		Prev:  prev,
	}, nil
}
