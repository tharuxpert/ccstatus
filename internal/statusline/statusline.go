// Package statusline provides the core statusline output functionality.
// This logic MUST remain unchanged to preserve the statusline output format.
package statusline

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Input represents the JSON input from Claude Code
type Input struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
}

// Credentials represents the OAuth credentials from Keychain
type Credentials struct {
	ClaudeAiOauth struct {
		AccessToken string `json:"accessToken"`
	} `json:"claudeAiOauth"`
}

// UsageResponse represents the API response from Anthropic
type UsageResponse struct {
	FiveHour struct {
		Utilization float64 `json:"utilization"`
		ResetsAt    string  `json:"resets_at"`
	} `json:"five_hour"`
	SevenDay struct {
		Utilization float64 `json:"utilization"`
		ResetsAt    string  `json:"resets_at"`
	} `json:"seven_day"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Run executes the statusline logic and prints output to stdout.
// This function MUST NOT be modified - it produces the exact statusline format.
func Run() {
	// Read model info from stdin
	model := readModelFromStdin()

	// Get OAuth token from macOS Keychain
	token, err := GetAccessToken()
	if err != nil || token == "" {
		printFallback(model)
		return
	}

	// Fetch usage data from Anthropic API
	usage, err := FetchUsage(token)
	if err != nil || usage.Error != nil {
		printFallback(model)
		return
	}

	// Format and print statusline
	printStatusLine(model, usage)
}

// readModelFromStdin reads and parses the JSON input from stdin
func readModelFromStdin() string {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "Unknown"
	}

	var input Input
	if err := json.Unmarshal(data, &input); err != nil {
		return "Unknown"
	}

	if input.Model.DisplayName == "" {
		return "Unknown"
	}
	return input.Model.DisplayName
}

// GetAccessToken retrieves the OAuth token from macOS Keychain
func GetAccessToken() (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", "Claude Code-credentials", "-w")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	credsJSON := strings.TrimSpace(string(output))
	if credsJSON == "" {
		return "", fmt.Errorf("empty credentials")
	}

	var creds Credentials
	if err := json.Unmarshal([]byte(credsJSON), &creds); err != nil {
		return "", err
	}

	return creds.ClaudeAiOauth.AccessToken, nil
}

// FetchUsage retrieves usage data from the Anthropic API
func FetchUsage(token string) (*UsageResponse, error) {
	req, err := http.NewRequest("GET", "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var usage UsageResponse
	if err := json.Unmarshal(body, &usage); err != nil {
		return nil, err
	}

	return &usage, nil
}

// formatResetTime converts an ISO timestamp to local 12-hour format (e.g., "3:45pm")
func formatResetTime(isoTime string) string {
	if isoTime == "" {
		return "--"
	}

	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		// Try parsing with fractional seconds
		t, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", isoTime)
		if err != nil {
			return "--"
		}
	}

	// Convert to local timezone
	local := t.Local()

	// Format as "3:04pm" (12-hour format, lowercase am/pm)
	hour := local.Hour()
	minute := local.Minute()
	ampm := "am"
	if hour >= 12 {
		ampm = "pm"
	}
	if hour > 12 {
		hour -= 12
	}
	if hour == 0 {
		hour = 12
	}

	return fmt.Sprintf("%d:%02d%s", hour, minute, ampm)
}

// formatWeeklyResetTime converts an ISO timestamp to local format with date (e.g., "Jan 15 3:45pm")
func formatWeeklyResetTime(isoTime string) string {
	if isoTime == "" {
		return "--"
	}

	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		// Try parsing with fractional seconds
		t, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", isoTime)
		if err != nil {
			return "--"
		}
	}

	// Convert to local timezone
	local := t.Local()

	// Format as "Jan 2 3:04pm"
	month := local.Format("Jan")
	day := local.Day()
	hour := local.Hour()
	minute := local.Minute()
	ampm := "am"
	if hour >= 12 {
		ampm = "pm"
	}
	if hour > 12 {
		hour -= 12
	}
	if hour == 0 {
		hour = 12
	}

	return fmt.Sprintf("%s %d %d:%02d%s", month, day, hour, minute, ampm)
}

// printFallback prints the statusline with placeholder values
func printFallback(model string) {
	fmt.Printf("%s | Session: --%% | Week: --%%", model)
}

// printStatusLine formats and prints the full statusline
func printStatusLine(model string, usage *UsageResponse) {
	sessionPct := int(usage.FiveHour.Utilization)
	weeklyPct := int(usage.SevenDay.Utilization)
	sessionReset := formatResetTime(usage.FiveHour.ResetsAt)
	weeklyReset := formatWeeklyResetTime(usage.SevenDay.ResetsAt)

	fmt.Printf("%s | Session: %d%% (resets %s) | Week: %d%% (resets %s)",
		model, sessionPct, sessionReset, weeklyPct, weeklyReset)
}
