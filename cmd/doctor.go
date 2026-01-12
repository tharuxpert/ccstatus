package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"ccstatus/internal/config"
	"ccstatus/internal/statusline"
	"ccstatus/internal/ui"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check ccstatus configuration and dependencies",
	Long: `Doctor runs diagnostic checks to verify ccstatus is properly configured.

Checks performed:
  - Claude Code configuration exists
  - ccstatus is properly configured in settings
  - ccstatus binary is in PATH
  - OAuth token is available in Keychain
  - Anthropic API endpoint is reachable`,
	RunE: runDoctor,
}

type checkResult struct {
	name    string
	status  string
	message string
	ok      bool
}

func runDoctor(cmd *cobra.Command, args []string) error {
	ui.CompactTitle("ccstatus doctor")

	// Run all checks with spinners
	s := ui.NewSpinner("Running diagnostics...")
	s.Start()

	checks := []checkResult{
		checkConfigExists(),
		checkStatuslineConfigured(),
		checkBinaryInPath(),
		checkOAuthToken(),
		checkAPIEndpoint(),
	}

	s.Stop()

	// Print results
	fmt.Println()
	ui.Bold.Println("  Diagnostics")
	ui.Divider()
	fmt.Println()

	passCount := 0
	failCount := 0

	for _, check := range checks {
		if check.ok {
			ui.StatusOK(check.name, check.message)
			passCount++
		} else {
			ui.StatusError(check.name, check.message)
			failCount++
		}
	}

	// Summary
	fmt.Println()
	ui.Divider()

	if failCount == 0 {
		ui.SuccessMessage("All checks passed!", "ccstatus is ready to use.")
	} else {
		ui.ErrorMessage(
			fmt.Sprintf("%d of %d checks failed", failCount, len(checks)),
			"",
		)
		fmt.Println()
		ui.Bold.Println("  Quick fixes:")
		fmt.Println()
		ui.Bullet("Run " + ui.InfoBold.Sprint("ccstatus install") + " to configure the statusline")
		ui.Bullet("Ensure ccstatus is in your PATH")
		ui.Bullet("Sign in to Claude Code to generate OAuth credentials")
	}

	fmt.Println()
	return nil
}

func checkConfigExists() checkResult {
	result := checkResult{
		name: "Claude Code configuration",
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		result.ok = false
		result.message = fmt.Sprintf("Cannot determine config path: %v", err)
		return result
	}

	exists, err := config.ConfigExists()
	if err != nil {
		result.ok = false
		result.message = fmt.Sprintf("Error checking config: %v", err)
		return result
	}

	if !exists {
		result.ok = false
		result.message = fmt.Sprintf("Not found at %s", configPath)
		return result
	}

	result.ok = true
	result.message = configPath
	return result
}

func checkStatuslineConfigured() checkResult {
	result := checkResult{
		name: "Statusline configuration",
	}

	settings, err := config.ReadSettings()
	if err != nil {
		result.ok = false
		result.message = fmt.Sprintf("Cannot read settings: %v", err)
		return result
	}

	cmd := config.GetStatuslineCommand(settings)

	if cmd == "" {
		result.ok = false
		result.message = "Not configured"
		return result
	}

	if cmd != "ccstatus" {
		result.ok = false
		result.message = fmt.Sprintf("Different command: %s", cmd)
		return result
	}

	result.ok = true
	result.message = "ccstatus"
	return result
}

func checkBinaryInPath() checkResult {
	result := checkResult{
		name: "Binary in PATH",
	}

	path, err := exec.LookPath("ccstatus")
	if err != nil {
		result.ok = false
		result.message = "ccstatus not found in PATH"
		return result
	}

	result.ok = true
	result.message = path
	return result
}

func checkOAuthToken() checkResult {
	result := checkResult{
		name: "OAuth token",
	}

	token, err := statusline.GetAccessToken()
	if err != nil {
		result.ok = false
		result.message = fmt.Sprintf("Cannot retrieve: %v", err)
		return result
	}

	if token == "" {
		result.ok = false
		result.message = "Empty token - sign in to Claude Code"
		return result
	}

	// Mask token for display
	maskedToken := token[:8] + "..." + token[len(token)-4:]
	result.ok = true
	result.message = maskedToken
	return result
}

func checkAPIEndpoint() checkResult {
	result := checkResult{
		name: "Anthropic API",
	}

	// First check if we have a token
	token, err := statusline.GetAccessToken()
	if err != nil || token == "" {
		result.ok = false
		result.message = "Skipped (no token)"
		return result
	}

	// Try to reach the API
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		result.ok = false
		result.message = fmt.Sprintf("Request error: %v", err)
		return result
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")

	resp, err := client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			result.ok = false
			result.message = "Request timed out"
			return result
		}
		result.ok = false
		result.message = fmt.Sprintf("Connection failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		result.ok = false
		result.message = "Token rejected (401)"
		return result
	}

	if resp.StatusCode != 200 {
		result.ok = false
		result.message = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return result
	}

	result.ok = true
	result.message = "Reachable"
	return result
}
