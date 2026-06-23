package cmd

// Me represents the user information from the /me endpoint
type Me struct {
	Country     string
	DisplayName string
	Email       string
	Product     string
	ProfileURL  string
}

// fetchMe fetches the user's information from the /me endpoint
func fetchMe() (Me, error) {
	client := NewSpotifyClient()
	return client.GetMe()
}
