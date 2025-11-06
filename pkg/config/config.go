package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	// ConfigDir is the directory where config files are stored
	ConfigDir = ".config/gs-write"
	// ConfigFile is the name of the user config file
	ConfigFile = "config.toml"
)

// UserConfig represents the user configuration
type UserConfig struct {
	Freeze FreezeConfig `toml:"freeze"`
	Filter FilterConfig `toml:"filter"`
}

// FreezeConfig represents freeze panes configuration
type FreezeConfig struct {
	Rows *int `toml:"rows,omitempty"`
	Cols *int `toml:"cols,omitempty"`
}

// FilterConfig represents basic filter configuration
type FilterConfig struct {
	HeaderRow *int `toml:"header_row,omitempty"`
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ConfigDir, ConfigFile), nil
}

// LoadUserConfig loads the user configuration from config.toml
func LoadUserConfig() (*UserConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &UserConfig{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config UserConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveUserConfig saves the user configuration to config.toml
func SaveUserConfig(config *UserConfig) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetFreezeRows returns the freeze rows setting from config
func (c *UserConfig) GetFreezeRows() (int, bool) {
	if c.Freeze.Rows != nil {
		return *c.Freeze.Rows, true
	}
	return 0, false
}

// GetFreezeCols returns the freeze cols setting from config
func (c *UserConfig) GetFreezeCols() (int, bool) {
	if c.Freeze.Cols != nil {
		return *c.Freeze.Cols, true
	}
	return 0, false
}

// SetFreezeRows sets the freeze rows setting
func (c *UserConfig) SetFreezeRows(rows int) {
	c.Freeze.Rows = &rows
}

// SetFreezeCols sets the freeze cols setting
func (c *UserConfig) SetFreezeCols(cols int) {
	c.Freeze.Cols = &cols
}

// UnsetFreezeRows removes the freeze rows setting
func (c *UserConfig) UnsetFreezeRows() {
	c.Freeze.Rows = nil
}

// UnsetFreezeCols removes the freeze cols setting
func (c *UserConfig) UnsetFreezeCols() {
	c.Freeze.Cols = nil
}

// GetFilterHeaderRow returns the filter header row setting from config
func (c *UserConfig) GetFilterHeaderRow() (int, bool) {
	if c.Filter.HeaderRow != nil {
		return *c.Filter.HeaderRow, true
	}
	return 0, false
}

// SetFilterHeaderRow sets the filter header row setting
func (c *UserConfig) SetFilterHeaderRow(row int) {
	c.Filter.HeaderRow = &row
}

// UnsetFilterHeaderRow removes the filter header row setting
func (c *UserConfig) UnsetFilterHeaderRow() {
	c.Filter.HeaderRow = nil
}
