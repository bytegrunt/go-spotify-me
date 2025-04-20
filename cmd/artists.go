package cmd

import (
	"strings"

	"github.com/bytegrunt/go-spotify-me/internal/auth"
)

// Artist represents an artist's details
type Artist struct {
	Name       string
	Genres     string
	Popularity int
}

func fetchArtistsPage(url string) (APIResponse, error) {
	token, _ := auth.GetValidAccessToken()
	response, err := MakeAPIRequest(token, url)
	if err != nil {
		return APIResponse{}, err
	}

	artists := parseArtists(response)
	next, _ := response["next"].(string)
	prev, _ := response["previous"].(string)

	return APIResponse{
		Artists: artists,
		Next:    next,
		Prev:    prev,
	}, nil
}

func parseArtists(response map[string]interface{}) []Artist {
	items, ok := response["items"].([]interface{})
	if !ok {
		return nil
	}

	var artists []Artist
	for _, item := range items {
		artist, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name := artist["name"].(string)

		genresInterface, ok := artist["genres"].([]interface{})
		if !ok {
			genresInterface = []interface{}{}
		}
		var genres []string
		for _, genre := range genresInterface {
			genres = append(genres, genre.(string))
		}

		popularity := int(artist["popularity"].(float64))

		artists = append(artists, Artist{
			Name:       name,
			Genres:     strings.Join(genres, ", "),
			Popularity: popularity,
		})
	}

	return artists
}
