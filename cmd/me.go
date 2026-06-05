package cmd

import "github.com/CyberGrit/go-spotify-me/internal/spotify"

// fetchMe fetches the user's information from the /me endpoint
func fetchMe() (spotify.Me, error) {
	return spotifyClient.GetMe()
}
