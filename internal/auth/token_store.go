package auth

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CyberGrit/go-spotify-me/internal/logging"
	"github.com/zalando/go-keyring"
	"go.uber.org/zap"
)

type TokenStore interface {
	GetValidAccessToken() (string, bool)
	GetRefreshToken() (string, error)
	SaveAccessToken(accessToken, refreshToken string, expirationTime time.Time) error
}

type OSTokenStore struct {
	logger *zap.Logger
}

func NewOSTokenStore(logger *zap.Logger) *OSTokenStore {
	return &OSTokenStore{
		logger: logger,
	}
}

func (s *OSTokenStore) GetValidAccessToken() (string, bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		if s.logger != nil {
			s.logger.Error("Failed to get user home directory", zap.Error(err))
		}
		return "", false
	}

	filePath := filepath.Join(homeDir, ".go-spotify-me-cli")

	if !strings.HasPrefix(filePath, homeDir) {
		if s.logger != nil {
			s.logger.Error("Invalid file path", zap.String("filePath", filePath))
		}
		return "", false
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o600)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("Failed to open file", zap.Error(err))
		}
		return "", false
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("Failed to read token file", zap.Error(err))
		}
		return "", false
	}

	lines := strings.Split(string(data), "\n")
	var accessToken, expiresAtStr string
	for _, line := range lines {
		if strings.HasPrefix(line, "access_token=") {
			accessToken = strings.TrimPrefix(line, "access_token=")
		} else if strings.HasPrefix(line, "expires_at=") {
			expiresAtStr = strings.TrimPrefix(line, "expires_at=")
		}
	}

	if accessToken == "" || expiresAtStr == "" {
		return "", false
	}

	expirationTime, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("Failed to parse expiration time", zap.Error(err))
		}
		return "", false
	}

	if time.Now().After(expirationTime) {
		return "", false
	}

	return accessToken, true
}

func (s *OSTokenStore) GetRefreshToken() (string, error) {
	token, err := keyring.Get("go-spotify-me-cli", "refresh_token")
	if err == nil && token != "" {
		return token, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, ".go-spotify-me-cli")
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o600)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "refresh_token=") {
			return strings.TrimPrefix(line, "refresh_token="), nil
		}
	}

	return "", fmt.Errorf("refresh token not found")
}

func (s *OSTokenStore) SaveAccessToken(accessToken, refreshToken string, expirationTime time.Time) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, ".go-spotify-me-cli")

	if !strings.HasPrefix(filePath, homeDir) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	var data string

	err = keyring.Set("go-spotify-me-cli", "refresh_token", refreshToken)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("Failed to store refresh token in keyring", zap.Error(err))
			s.logger.Info("Falling back to saving the refresh token in the hidden file.")
		}
		data = fmt.Sprintf("access_token=%s\nrefresh_token=%s\nexpires_at=%s\n", accessToken, refreshToken, expirationTime.Format(time.RFC3339))
	} else {
		data = fmt.Sprintf("access_token=%s\nexpires_at=%s\n", accessToken, expirationTime.Format(time.RFC3339))
	}

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	logging.DebugLog("Access token saved to %s", filePath)
	return nil
}
