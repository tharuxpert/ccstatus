// Package config handles Claude Code configuration file operations.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// ConfigDir is the Claude config directory name
	ConfigDir = ".claude"
	// SettingsFile is the settings filename
	SettingsFile = "settings.json"
	// BackupPrefix is the prefix for backup files
	BackupPrefix = "settings.backup"
)

// Settings represents the Claude Code settings.json structure.
// We use map[string]any to preserve unknown fields.
type Settings map[string]any

// StatuslineConfig represents the statusline configuration
type StatuslineConfig struct {
	Command string `json:"command"`
}

// GetConfigPath returns the path to ~/.claude/settings.json
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ConfigDir, SettingsFile), nil
}

// GetConfigDir returns the path to ~/.claude/
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ConfigDir), nil
}

// ConfigExists checks if the settings file exists
func ConfigExists() (bool, error) {
	path, err := GetConfigPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// ReadSettings reads and parses the settings file
func ReadSettings() (Settings, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(Settings), nil
		}
		return nil, fmt.Errorf("cannot read settings file: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("cannot parse settings file: %w", err)
	}

	if settings == nil {
		settings = make(Settings)
	}

	return settings, nil
}

// WriteSettings writes settings to the config file
func WriteSettings(settings Settings) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal settings: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("cannot write settings file: %w", err)
	}

	return nil
}

// CreateBackup creates a timestamped backup of the settings file
func CreateBackup() (string, error) {
	path, err := GetConfigPath()
	if err != nil {
		return "", err
	}

	// Check if source file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil // No file to backup
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(filepath.Dir(path), fmt.Sprintf("%s.%s.json", BackupPrefix, timestamp))

	// Read original file
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read settings file for backup: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return "", fmt.Errorf("cannot write backup file: %w", err)
	}

	return backupPath, nil
}

// GetLatestBackup returns the path to the most recent backup file
func GetLatestBackup() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return "", fmt.Errorf("cannot read config directory: %w", err)
	}

	var latestBackup string
	var latestTime time.Time

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if len(name) > len(BackupPrefix) && name[:len(BackupPrefix)] == BackupPrefix {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
				latestBackup = filepath.Join(configDir, name)
			}
		}
	}

	if latestBackup == "" {
		return "", fmt.Errorf("no backup files found")
	}

	return latestBackup, nil
}

// RestoreFromBackup restores settings from a backup file
func RestoreFromBackup(backupPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("cannot read backup file: %w", err)
	}

	// Validate it's valid JSON
	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("backup file is not valid JSON: %w", err)
	}

	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("cannot restore settings: %w", err)
	}

	return nil
}

// GetStatuslineCommand returns the current statusline command if configured
func GetStatuslineCommand(settings Settings) string {
	statusline, ok := settings["statusline"]
	if !ok {
		return ""
	}

	statuslineMap, ok := statusline.(map[string]any)
	if !ok {
		return ""
	}

	cmd, ok := statuslineMap["command"].(string)
	if !ok {
		return ""
	}

	return cmd
}

// SetStatuslineCommand sets the statusline command in settings
func SetStatuslineCommand(settings Settings, command string) {
	settings["statusline"] = map[string]any{
		"command": command,
	}
}

// RemoveStatusline removes the statusline configuration from settings
func RemoveStatusline(settings Settings) {
	delete(settings, "statusline")
}

// IsStatuslineConfigured checks if ccstatus is already configured
func IsStatuslineConfigured(settings Settings) bool {
	return GetStatuslineCommand(settings) == "ccstatus"
}
