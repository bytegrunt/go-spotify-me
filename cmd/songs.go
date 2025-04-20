package cmd

import (
	"github.com/bytegrunt/go-spotify-me/internal/auth"
)

type Song struct {
	Name       string
	Artist     string
	Album      string
	Popularity int
}

func fetchSongsPage(url string) (APIResponse, error) {
	token, _ := auth.GetValidAccessToken()
	response, err := MakeAPIRequest(token, url)
	if err != nil {
		return APIResponse{}, err
	}

	songs := parseSongs(response)
	next, _ := response["next"].(string)
	prev, _ := response["previous"].(string)

	return APIResponse{
		Songs: songs,
		Next:  next,
		Prev:  prev,
	}, nil
}

func parseSongs(response map[string]interface{}) []Song {
	items, ok := response["items"].([]interface{})
	if !ok {
		return nil
	}

	var songs []Song
	for _, item := range items {
		track, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name := track["name"].(string)
		popularity := int(track["popularity"].(float64))

		albumName := ""
		if album, ok := track["album"].(map[string]interface{}); ok {
			albumName, _ = album["name"].(string)
		}

		artistName := ""
		if artistList, ok := track["artists"].([]interface{}); ok && len(artistList) > 0 {
			if firstArtist, ok := artistList[0].(map[string]interface{}); ok {
				artistName, _ = firstArtist["name"].(string)
			}
		}

		songs = append(songs, Song{
			Name:       name,
			Artist:     artistName,
			Album:      albumName,
			Popularity: popularity,
		})
	}

	return songs
}
