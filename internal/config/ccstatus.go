// Package config - ccstatus-specific configuration
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// CCStatusConfigFile is the ccstatus config filename
	CCStatusConfigFile = "ccstatus.json"
)

// CCStatusConfig represents ccstatus-specific configuration options
type CCStatusConfig struct {
	ShowSessionUsage bool `json:"show_session_usage"`
	ShowWeeklyUsage  bool `json:"show_weekly_usage"`
	ShowResetTimes   bool `json:"show_reset_times"`
	ShowGitBranch    bool `json:"show_git_branch"`
}

// DefaultCCStatusConfig returns the default configuration
func DefaultCCStatusConfig() *CCStatusConfig {
	return &CCStatusConfig{
		ShowSessionUsage: true,
		ShowWeeklyUsage:  true,
		ShowResetTimes:   true,
		ShowGitBranch:    false,
	}
}

// GetCCStatusConfigPath returns the path to ~/.claude/ccstatus.json
func GetCCStatusConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ConfigDir, CCStatusConfigFile), nil
}

// LoadCCStatusConfig loads the ccstatus configuration from disk
func LoadCCStatusConfig() (*CCStatusConfig, error) {
	path, err := GetCCStatusConfigPath()
	if err != nil {
		return DefaultCCStatusConfig(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return defaults if no config file exists
			return DefaultCCStatusConfig(), nil
		}
		return DefaultCCStatusConfig(), fmt.Errorf("cannot read config file: %w", err)
	}

	var cfg CCStatusConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultCCStatusConfig(), fmt.Errorf("cannot parse config file: %w", err)
	}

	return &cfg, nil
}

// SaveCCStatusConfig saves the ccstatus configuration to disk
func SaveCCStatusConfig(cfg *CCStatusConfig) error {
	path, err := GetCCStatusConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("cannot write config file: %w", err)
	}

	return nil
}
