package statusline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ccstatus/internal/config"
)

func TestClaudeCodeUserAgent(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{name: "with version", version: "1.0.80", want: "claude-code/1.0.80"},
		{name: "without version", version: "", want: "claude-code"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := claudeCodeUserAgent(tt.version); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestSaveCacheCreatesClaudeDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	usage := &UsageResponse{}
	usage.FiveHour.Utilization = 42

	saveCache(usage)

	path := filepath.Join(home, config.ConfigDir, cacheFile)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected cache file to be written: %v", err)
	}

	cached, ok := loadCache()
	if !ok {
		t.Fatal("expected fresh cache to load")
	}
	if got := cached.FiveHour.Utilization; got != 42 {
		t.Fatalf("expected cached utilization 42, got %v", got)
	}
}

func TestLoadCacheRejectsExpiredCacheButAllowsStaleFallback(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cacheDir := filepath.Join(home, config.ConfigDir)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatal(err)
	}

	cached := cachedUsage{
		Usage:     UsageResponse{},
		FetchedAt: time.Now().Add(-defaultTTL - time.Minute),
	}
	cached.Usage.SevenDay.Utilization = 73

	data, err := json.Marshal(cached)
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(cacheDir, cacheFile)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	if _, ok := loadCache(); ok {
		t.Fatal("expected expired cache to be rejected")
	}

	stale, ok := loadStaleCache()
	if !ok {
		t.Fatal("expected stale cache fallback to load")
	}
	if got := stale.SevenDay.Utilization; got != 73 {
		t.Fatalf("expected stale utilization 73, got %v", got)
	}
}
