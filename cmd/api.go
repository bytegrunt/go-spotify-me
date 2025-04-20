/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bytegrunt/go-spotify-me/internal/auth"

)

// Artist represents an artist's details
type Artist struct {
	Name       string
	Genres     string
	Popularity int
}

type Song struct {
	Name       string
	Artist     string
	Album      string
	Popularity int
}

type APIResponse struct {
	Artists []Artist
	Songs   []Song
	Next    string
	Prev    string
}


// MakeAPIRequest makes a GET request to the Spotify API and returns the response or an error
func MakeAPIRequest(token string, url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return response, nil
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
