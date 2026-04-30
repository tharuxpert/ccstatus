package statusline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"ccstatus/internal/config"
)

const (
	cacheFile  = "ccstatus-cache.json"
	defaultTTL = 5 * time.Minute
)

type cachedUsage struct {
	Usage     UsageResponse `json:"usage"`
	FetchedAt time.Time     `json:"fetched_at"`
}

func getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, config.ConfigDir, cacheFile), nil
}

func readCache() (*cachedUsage, bool) {
	path, err := getCachePath()
	if err != nil {
		return nil, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var cached cachedUsage
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, false
	}

	return &cached, true
}

func loadCache() (*UsageResponse, bool) {
	cached, ok := readCache()
	if !ok || time.Since(cached.FetchedAt) > defaultTTL {
		return nil, false
	}
	return &cached.Usage, true
}

func loadStaleCache() (*UsageResponse, bool) {
	cached, ok := readCache()
	if !ok {
		return nil, false
	}
	return &cached.Usage, true
}

func saveCache(usage *UsageResponse) {
	if usage == nil {
		return
	}

	path, err := getCachePath()
	if err != nil {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return
	}

	cached := cachedUsage{
		Usage:     *usage,
		FetchedAt: time.Now(),
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return
	}

	_ = os.WriteFile(path, data, 0600)
}
