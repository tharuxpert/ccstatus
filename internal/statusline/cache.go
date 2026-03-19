package statusline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheFile   = "ccstatus-cache.json"
	defaultTTL  = 5 * time.Minute
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
	return filepath.Join(home, ".claude", cacheFile), nil
}

// loadCache returns a cached UsageResponse if it exists and is fresh.
func loadCache() (*UsageResponse, bool) {
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

	if time.Since(cached.FetchedAt) > defaultTTL {
		return nil, false
	}

	return &cached.Usage, true
}

// loadCacheIgnoringTTL returns a cached UsageResponse even if expired.
// Used as a fallback when the API is unavailable.
func loadCacheIgnoringTTL() (*UsageResponse, bool) {
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

	return &cached.Usage, true
}

// saveCache writes a UsageResponse to the cache file.
func saveCache(usage *UsageResponse) {
	path, err := getCachePath()
	if err != nil {
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
