package cmd

import (
	"errors"
	"testing"

	"github.com/CyberGrit/go-spotify-me/internal/spotify"
)

type mockSpotifyClient struct {
	mockArtists []spotify.Artist
	mockSongs   []spotify.Song
	mockMe      spotify.Me
	mockNext    string
	mockPrev    string
	mockErr     error
}

func (m *mockSpotifyClient) GetTopArtists(timeRange string) ([]spotify.Artist, string, string, error) {
	return m.mockArtists, m.mockNext, m.mockPrev, m.mockErr
}

func (m *mockSpotifyClient) GetTopSongs(timeRange string) ([]spotify.Song, string, string, error) {
	return m.mockSongs, m.mockNext, m.mockPrev, m.mockErr
}

func (m *mockSpotifyClient) GetMe() (spotify.Me, error) {
	return m.mockMe, m.mockErr
}

func (m *mockSpotifyClient) GetArtistsPage(url string) ([]spotify.Artist, string, string, error) {
	return m.mockArtists, m.mockNext, m.mockPrev, m.mockErr
}

func (m *mockSpotifyClient) GetSongsPage(url string) ([]spotify.Song, string, string, error) {
	return m.mockSongs, m.mockNext, m.mockPrev, m.mockErr
}

func TestFetchArtistsPage(t *testing.T) {
	origClient := spotifyClient
	defer func() { spotifyClient = origClient }()

	expectedArtists := []spotify.Artist{
		{Name: "Artist 1", Genres: "Pop", Popularity: 80},
	}
	spotifyClient = &mockSpotifyClient{
		mockArtists: expectedArtists,
		mockNext:    "next_url",
		mockPrev:    "prev_url",
	}

	resp, err := fetchArtistsPage("some_url")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Artists) != 1 || resp.Artists[0].Name != "Artist 1" {
		t.Errorf("unexpected artists: %+v", resp.Artists)
	}
	if resp.Next != "next_url" {
		t.Errorf("expected next url next_url, got %v", resp.Next)
	}
}

func TestFetchSongsPage(t *testing.T) {
	origClient := spotifyClient
	defer func() { spotifyClient = origClient }()

	expectedSongs := []spotify.Song{
		{Name: "Song 1", Artist: "Artist 1", Album: "Album 1", Popularity: 90},
	}
	spotifyClient = &mockSpotifyClient{
		mockSongs: expectedSongs,
		mockNext:  "next_url_s",
		mockPrev:  "prev_url_s",
	}

	resp, err := fetchSongsPage("some_url")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Songs) != 1 || resp.Songs[0].Name != "Song 1" {
		t.Errorf("unexpected songs: %+v", resp.Songs)
	}
}

func TestFetchMe(t *testing.T) {
	origClient := spotifyClient
	defer func() { spotifyClient = origClient }()

	expectedMe := spotify.Me{
		DisplayName: "Test User",
		Email:       "test@example.com",
	}
	spotifyClient = &mockSpotifyClient{
		mockMe: expectedMe,
	}

	me, err := fetchMe()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if me.DisplayName != "Test User" {
		t.Errorf("expected DisplayName Test User, got %v", me.DisplayName)
	}
}

func TestFetchArtistsPage_Error(t *testing.T) {
	origClient := spotifyClient
	defer func() { spotifyClient = origClient }()

	spotifyClient = &mockSpotifyClient{
		mockErr: errors.New("mock error"),
	}

	_, err := fetchArtistsPage("some_url")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
