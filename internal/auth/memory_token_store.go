package auth

import (
	"fmt"
	"time"
)

// MemoryTokenStore is an in-memory implementation of the TokenStore interface for testing.
type MemoryTokenStore struct {
	refreshToken string
	accessToken  string
	expiresAt    time.Time
}

// NewMemoryTokenStore returns a new MemoryTokenStore.
func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{}
}

// GetRefreshToken returns the stored refresh token.
func (m *MemoryTokenStore) GetRefreshToken() (string, error) {
	if m.refreshToken == "" {
		return "", fmt.Errorf("refresh token not found")
	}
	return m.refreshToken, nil
}

// SaveTokens stores the access token, refresh token, and expiration time.
func (m *MemoryTokenStore) SaveTokens(accessToken, refreshToken string, expiresAt time.Time) error {
	m.accessToken = accessToken
	m.refreshToken = refreshToken
	m.expiresAt = expiresAt
	return nil
}

// GetValidAccessToken returns the access token if it is still valid, otherwise returns an empty string and false.
func (m *MemoryTokenStore) GetValidAccessToken() (string, bool) {
	if m.accessToken == "" || m.expiresAt.IsZero() {
		return "", false
	}
	if time.Now().After(m.expiresAt) {
		return "", false
	}
	return m.accessToken, true
}