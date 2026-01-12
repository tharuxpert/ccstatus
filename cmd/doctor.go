package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"ccstatus/internal/config"
	"ccstatus/internal/statusline"

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
	fmt.Println("ccstatus doctor")
	fmt.Println()

	checks := []checkResult{
		checkConfigExists(),
		checkStatuslineConfigured(),
		checkBinaryInPath(),
		checkOAuthToken(),
		checkAPIEndpoint(),
	}

	// Print results
	allOK := true
	for _, check := range checks {
		statusIcon := "[OK]"
		if !check.ok {
			statusIcon = "[!!]"
			allOK = false
		}

		fmt.Printf("  %s %s\n", statusIcon, check.name)
		if check.message != "" {
			fmt.Printf("      %s\n", check.message)
		}
	}

	fmt.Println()

	if allOK {
		fmt.Println("All checks passed. ccstatus is ready to use.")
	} else {
		fmt.Println("Some checks failed. See above for details.")
		fmt.Println()
		fmt.Println("Common fixes:")
		fmt.Println("  - Run 'ccstatus install' to configure the statusline")
		fmt.Println("  - Ensure ccstatus is in your PATH")
		fmt.Println("  - Sign in to Claude Code to generate OAuth credentials")
	}

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
		result.message = fmt.Sprintf("Config not found at %s", configPath)
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
		result.message = "No statusline configured. Run 'ccstatus install' to set up."
		return result
	}

	if cmd != "ccstatus" {
		result.ok = false
		result.message = fmt.Sprintf("Different command configured: %s", cmd)
		return result
	}

	result.ok = true
	result.message = "ccstatus is configured"
	return result
}

func checkBinaryInPath() checkResult {
	result := checkResult{
		name: "ccstatus in PATH",
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
		result.message = fmt.Sprintf("Cannot retrieve token: %v", err)
		return result
	}

	if token == "" {
		result.ok = false
		result.message = "Token is empty. Sign in to Claude Code."
		return result
	}

	// Mask token for display
	maskedToken := token[:8] + "..." + token[len(token)-4:]
	result.ok = true
	result.message = fmt.Sprintf("Found (%s)", maskedToken)
	return result
}

func checkAPIEndpoint() checkResult {
	result := checkResult{
		name: "Anthropic API endpoint",
	}

	// First check if we have a token
	token, err := statusline.GetAccessToken()
	if err != nil || token == "" {
		result.ok = false
		result.message = "Skipped (no OAuth token available)"
		return result
	}

	// Try to reach the API
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		result.ok = false
		result.message = fmt.Sprintf("Cannot create request: %v", err)
		return result
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")

	resp, err := client.Do(req)
	if err != nil {
		// Check if it's a network error vs other error
		if os.IsTimeout(err) {
			result.ok = false
			result.message = "Request timed out"
			return result
		}
		result.ok = false
		result.message = fmt.Sprintf("Cannot reach API: %v", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		result.ok = false
		result.message = "Token rejected (401). Try signing in to Claude Code again."
		return result
	}

	if resp.StatusCode != 200 {
		result.ok = false
		result.message = fmt.Sprintf("Unexpected status: %d", resp.StatusCode)
		return result
	}

	result.ok = true
	result.message = "Reachable (status 200)"
	return result
}
