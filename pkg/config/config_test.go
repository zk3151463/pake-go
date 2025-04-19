package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary test config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.json")

	// Test case 1: File does not exist
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if config.URL != "" {
		t.Errorf("Expected empty URL, got %s", config.URL)
	}
	if config.Width != 1024 {
		t.Errorf("Expected default width 1024, got %d", config.Width)
	}
	if config.Height != 768 {
		t.Errorf("Expected default height 768, got %d", config.Height)
	}

	// Test case 2: File exists with custom values
	testConfig := &Config{
		URL:    "https://test.com",
		Name:   "TestApp",
		Width:  800,
		Height: 600,
	}

	if err := SaveConfig(testConfig, configPath); err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if loadedConfig.URL != testConfig.URL {
		t.Errorf("Expected URL %s, got %s", testConfig.URL, loadedConfig.URL)
	}
	if loadedConfig.Name != testConfig.Name {
		t.Errorf("Expected Name %s, got %s", testConfig.Name, loadedConfig.Name)
	}
	if loadedConfig.Width != testConfig.Width {
		t.Errorf("Expected Width %d, got %d", testConfig.Width, loadedConfig.Width)
	}
	if loadedConfig.Height != testConfig.Height {
		t.Errorf("Expected Height %d, got %d", testConfig.Height, loadedConfig.Height)
	}
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.json")

	// Test case 1: Save new config
	testConfig := &Config{
		URL:    "https://test.com",
		Name:   "TestApp",
		Width:  800,
		Height: 600,
	}

	if err := SaveConfig(testConfig, configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created")
	}

	// Test case 2: Save to directory with nested subdirectories
	nestedPath := filepath.Join(tempDir, "nested", "subdirs", "config.json")
	if err := SaveConfig(testConfig, nestedPath); err != nil {
		t.Fatalf("Failed to save config to nested directory: %v", err)
	}

	// Verify nested directories were created
	if _, err := os.Stat(filepath.Dir(nestedPath)); os.IsNotExist(err) {
		t.Fatalf("Nested directories were not created")
	}

	// Test case 3: Save with invalid permissions
	if runtime.GOOS != "windows" {
		readOnlyDir := filepath.Join(tempDir, "readonly")
		if err := os.MkdirAll(readOnlyDir, 0444); err != nil {
			t.Fatalf("Failed to create read-only directory: %v", err)
		}
		invalidPath := filepath.Join(readOnlyDir, "config.json")
		err := SaveConfig(testConfig, invalidPath)
		if err == nil {
			t.Error("Expected error when saving to read-only directory")
		}
	}
}
