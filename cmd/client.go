package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/CyberGrit/go-spotify-me/internal/auth"
)

type SpotifyClient interface {
	GetTopArtists(timeRange string) ([]Artist, string, string, error)
	GetTopSongs(timeRange string) ([]Song, string, string, error)
	GetMe() (Me, error)
	GetArtistsPage(url string) ([]Artist, string, string, error)
	GetSongsPage(url string) ([]Song, string, string, error)
}

type DefaultSpotifyClient struct {
	httpClient    *http.Client
	tokenProvider func() (string, bool)
}

func NewSpotifyClient() *DefaultSpotifyClient {
	return &DefaultSpotifyClient{
		httpClient:    &http.Client{},
		tokenProvider: auth.GetValidAccessToken,
	}
}

func (c *DefaultSpotifyClient) doRequest(url string) (map[string]interface{}, error) {
	token, valid := c.tokenProvider()
	if !valid {
		return nil, fmt.Errorf("failed to get a valid access token")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
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

func (c *DefaultSpotifyClient) GetTopArtists(timeRange string) ([]Artist, string, string, error) {
	url := "https://api.spotify.com/v1/me/top/artists?time_range=" + timeRange
	return c.GetArtistsPage(url)
}

func (c *DefaultSpotifyClient) GetTopSongs(timeRange string) ([]Song, string, string, error) {
	url := "https://api.spotify.com/v1/me/top/tracks?time_range=" + timeRange
	return c.GetSongsPage(url)
}

func (c *DefaultSpotifyClient) GetMe() (Me, error) {
	url := "https://api.spotify.com/v1/me"
	response, err := c.doRequest(url)
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

func (c *DefaultSpotifyClient) GetArtistsPage(url string) ([]Artist, string, string, error) {
	response, err := c.doRequest(url)
	if err != nil {
		return nil, "", "", err
	}

	artists := parseArtists(response)
	next, _ := response["next"].(string)
	prev, _ := response["previous"].(string)

	return artists, next, prev, nil
}

func (c *DefaultSpotifyClient) GetSongsPage(url string) ([]Song, string, string, error) {
	response, err := c.doRequest(url)
	if err != nil {
		return nil, "", "", err
	}

	songs := parseSongs(response)
	next, _ := response["next"].(string)
	prev, _ := response["previous"].(string)

	return songs, next, prev, nil
}
