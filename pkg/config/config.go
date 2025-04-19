package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	URL          string            `json:"url"`
	Name         string            `json:"name"`
	Icon         string            `json:"icon"`
	Width        int               `json:"width"`
	Height       int               `json:"height"`
	HideTitleBar bool              `json:"hideTitleBar"`
	Transparent  bool              `json:"transparent"`
	AlwaysOnTop  bool              `json:"alwaysOnTop"`
	UserAgent    string            `json:"userAgent"`
	Headers      map[string]string `json:"headers"`
	InjectCSS    []string          `json:"injectCSS"`
	InjectJS     []string          `json:"injectJS"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Width:        1024,
		Height:       768,
		HideTitleBar: false,
		Transparent:  false,
		AlwaysOnTop:  false,
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		Headers:      make(map[string]string),
		InjectCSS:    make([]string, 0),
		InjectJS:     make([]string, 0),
	}
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()

	// If the file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, path string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
