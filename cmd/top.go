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
	"github.com/spf13/cobra"
	tea "github.com/charmbracelet/bubbletea"
)

// Artist represents an artist's details
type Artist struct {
	Name       string
	Genres     string
	Popularity int
}

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Display your top Spotify artists in a TUI",
	Run: func(cmd *cobra.Command, args []string) {
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

		// Fetch initial page of artists
		url := "https://api.spotify.com/v1/me/top/artists"
		response, err := fetchArtistsPage(url)
		if err != nil {
			fmt.Println("Error fetching artists:", err)
			return
		}

		// Start the TUI
		p := tea.NewProgram(initialModel(response))
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

// fetchArtistsPage fetches a page of artists from the Spotify API
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

// parseArtists parses the API response to extract artist details
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

		// Convert genres from []interface{} to []string
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
