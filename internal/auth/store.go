package auth

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zalando/go-keyring"
	"go.uber.org/zap"

	"github.com/CyberGrit/go-spotify-me/internal/logging"
)

type TokenStore interface {
	GetValidAccessToken() (string, bool)
	GetRefreshToken() (string, error)
	SaveTokens(accessToken, refreshToken string, expirationTime time.Time) error
}

type osTokenStore struct{}

func NewOSTokenStore() TokenStore {
	return &osTokenStore{}
}

func (s *osTokenStore) getFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, ".go-spotify-me-cli")

	if !strings.HasPrefix(filePath, homeDir) {
		return "", fmt.Errorf("invalid file path: %s", filePath)
	}
	return filePath, nil
}

func (s *osTokenStore) GetValidAccessToken() (string, bool) {
	filePath, err := s.getFilePath()
	if err != nil {
		if logger != nil {
			logger.Error("Failed to get file path", zap.Error(err))
		}
		return "", false
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		if logger != nil {
			logger.Error("Failed to open file", zap.Error(err))
		}
		return "", false
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		if logger != nil {
			logger.Error("Failed to read token file", zap.Error(err))
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
		if logger != nil {
			logger.Error("Failed to parse expiration time", zap.Error(err))
		}
		return "", false
	}

	if time.Now().After(expirationTime) {
		return "", false
	}

	return accessToken, true
}

func (s *osTokenStore) GetRefreshToken() (string, error) {
	refreshToken, err := keyring.Get("go-spotify-me-cli", "refresh_token")
	if err == nil && refreshToken != "" {
		return refreshToken, nil
	}

	filePath, err := s.getFilePath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "refresh_token=") {
			return strings.TrimPrefix(line, "refresh_token="), nil
		}
	}

	return "", fmt.Errorf("refresh token not found")
}

func (s *osTokenStore) SaveTokens(accessToken, refreshToken string, expirationTime time.Time) error {
	filePath, err := s.getFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	var data string

	if refreshToken != "" {
		err = keyring.Set("go-spotify-me-cli", "refresh_token", refreshToken)
		if err != nil {
			if logger != nil {
				logger.Error("Failed to store refresh token in keyring", zap.Error(err))
				logger.Info("Falling back to saving the refresh token in the hidden file.")
			}
			data = fmt.Sprintf("access_token=%s\nrefresh_token=%s\nexpires_at=%s\n", accessToken, refreshToken, expirationTime.Format(time.RFC3339))
		} else {
			data = fmt.Sprintf("access_token=%s\nexpires_at=%s\n", accessToken, expirationTime.Format(time.RFC3339))
		}
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
