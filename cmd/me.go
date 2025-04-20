package cmd

import (
	"github.com/bytegrunt/go-spotify-me/internal/auth"
)

// Me represents the user information from the /me endpoint
type Me struct {
    DisplayName string
    Email       string
    Product     string
}

// fetchMe fetches the user's information from the /me endpoint
func fetchMe() (Me, error) {
    token, _ := auth.GetValidAccessToken()
    response, err := MakeAPIRequest(token, "https://api.spotify.com/v1/me")
    if err != nil {
        return Me{}, err
    }

    displayName, _ := response["display_name"].(string)
    email, _ := response["email"].(string)
    product, _ := response["product"].(string)

    return Me{
        DisplayName: displayName,
        Email:       email,
        Product:     product,
    }, nil
}