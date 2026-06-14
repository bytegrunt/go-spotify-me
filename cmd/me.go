package cmd

import (
	"github.com/CyberGrit/go-spotify-me/internal/auth"
)

// Me represents the user information from the /me endpoint
type Me struct {
	Country     string
	DisplayName string
	Email       string
	Product     string
	ProfileURL  string
}

// fetchMe fetches the user's information from the /me endpoint
func fetchMe(client SpotifyClient) (Me, error) {
	store := auth.NewOSTokenStore()
	token, _ := store.GetValidAccessToken()
	response, err := client.Get(token, "https://api.spotify.com/v1/me")
	if err != nil {
		return Me{}, err
	}

	country, _ := response["country"].(string)
	displayName, _ := response["display_name"].(string)
	email, _ := response["email"].(string)
	product, _ := response["product"].(string)
	externalURLs, _ := response["external_urls"].(map[string]interface{})
	profileURL, _ := externalURLs["spotify"].(string)

	return Me{
		Country:     country,
		DisplayName: displayName,
		Email:       email,
		Product:     product,
		ProfileURL:  profileURL,
	}, nil
}
