package auth

import "time"

// TokenStore defines the interface for token storage operations
type TokenStore interface {
	GetRefreshToken() (string, error)
	SaveTokens(accessToken, refreshToken string, expiresAt time.Time) error
	GetValidAccessToken() (string, bool)
}
