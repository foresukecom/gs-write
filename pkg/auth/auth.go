package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	// ConfigDir is the directory where config files are stored
	ConfigDir = ".config/gs-write"
	// AuthFile is the name of the authentication file
	AuthFile = "auth.json"
)

// AuthConfig represents the OAuth2 authentication configuration
type AuthConfig struct {
	Credentials *oauth2.Config `json:"credentials"`
	Token       *oauth2.Token  `json:"token"`
}

// GetAuthPath returns the full path to the auth file
func GetAuthPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ConfigDir, AuthFile), nil
}

// LoadAuthConfig loads the authentication configuration from the auth file
func LoadAuthConfig() (*AuthConfig, error) {
	authPath, err := GetAuthPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(authPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("auth file not found. Please run 'gs-write auth' first")
		}
		return nil, fmt.Errorf("failed to read auth file: %w", err)
	}

	var config AuthConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse auth file: %w", err)
	}

	return &config, nil
}

// SaveAuthConfig saves the authentication configuration to the auth file
func SaveAuthConfig(config *AuthConfig) error {
	authPath, err := GetAuthPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(authPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal auth config: %w", err)
	}

	if err := os.WriteFile(authPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write auth file: %w", err)
	}

	return nil
}

// ParseCredentials parses the credentials JSON
func ParseCredentials(credentialsJSON []byte) (*oauth2.Config, error) {
	config, err := google.ConfigFromJSON(credentialsJSON, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}
	return config, nil
}

// GetAuthURL generates the authorization URL
func GetAuthURL(config *oauth2.Config) string {
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

// ExchangeToken exchanges the authorization code for a token
func ExchangeToken(ctx context.Context, config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	return token, nil
}

// GetClient creates an HTTP client with the OAuth2 token
func GetClient(ctx context.Context) (*oauth2.Config, *oauth2.Token, error) {
	config, err := LoadAuthConfig()
	if err != nil {
		return nil, nil, err
	}

	// Check if token is valid
	if config.Token.Valid() {
		return config.Credentials, config.Token, nil
	}

	// Try to refresh the token
	tokenSource := config.Credentials.TokenSource(ctx, config.Token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to refresh token. Please run 'gs-write auth' again: %w", err)
	}

	// Save the new token
	config.Token = newToken
	if err := SaveAuthConfig(config); err != nil {
		return nil, nil, fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return config.Credentials, newToken, nil
}
