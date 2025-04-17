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
	"github.com/bytegrunt/go-spotify-me/internal/logging"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
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

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Display your top Spotify artists in a TUI",
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure user is logged in
		_, isValid := auth.GetValidAccessToken()
		if !isValid {
			logging.DebugLog("Access token is not valid. Running login command...")
			loginCmd, _, err := rootCmd.Find([]string{"login"})
			if err != nil {
				fmt.Println("Failed to find login command:", err)
				return
			}
			loginCmd.Run(loginCmd, args)
			_, isValid = auth.GetValidAccessToken()
			if !isValid {
				fmt.Println("Failed to obtain a valid access token after login.")
				return
			}
		}

		// Run the main app TUI
		p := tea.NewProgram(initialAppModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running TUI:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
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

	items, ok := response["items"].([]interface{})
	if !ok {
		return APIResponse{}, fmt.Errorf("invalid response format")
	}

	var artists []Artist
	for _, item := range items {
		artist, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := artist["name"].(string)
		popularity := int(artist["popularity"].(float64))

		genres := []string{}
		if genreList, ok := artist["genres"].([]interface{}); ok {
			for _, genre := range genreList {
				if g, ok := genre.(string); ok {
					genres = append(genres, g)
				}
			}
		}

		artists = append(artists, Artist{
			Name:       name,
			Genres:     strings.Join(genres, ", "),
			Popularity: popularity,
		})
	}

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

	items, ok := response["items"].([]interface{})
	if !ok {
		return APIResponse{}, fmt.Errorf("invalid response format")
	}

	var songs []Song
	for _, item := range items {
		track, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := track["name"].(string)
		popularity := int(track["popularity"].(float64))

		// Album name
		albumName := ""
		if album, ok := track["album"].(map[string]interface{}); ok {
			albumName, _ = album["name"].(string)
		}

		// Artist name
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

	next, _ := response["next"].(string)
	prev, _ := response["previous"].(string)

	return APIResponse{
		Songs: songs,
		Next:  next,
		Prev:  prev,
	}, nil
}
