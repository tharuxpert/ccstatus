package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"ccstatus/internal/config"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Configure ccstatus as the Claude Code statusline",
	Long: `Install configures ccstatus as the statusline command in Claude Code.

This command will:
  1. Detect your Claude Code configuration
  2. Create a timestamped backup of your settings
  3. Add the statusline configuration

You will be asked to confirm before any changes are made.`,
	RunE: runInstall,
}

func runInstall(cmd *cobra.Command, args []string) error {
	fmt.Println("ccstatus installer")
	fmt.Println()

	// Step 1: Check if config exists
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Checking Claude Code configuration..."
	s.Start()

	configPath, err := config.GetConfigPath()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to determine config path: %w", err)
	}

	exists, err := config.ConfigExists()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to check config: %w", err)
	}

	s.Stop()

	if !exists {
		fmt.Printf("  Config location: %s\n", configPath)
		fmt.Println("  Status: Not found (will be created)")
	} else {
		fmt.Printf("  Config location: %s\n", configPath)
		fmt.Println("  Status: Found")
	}
	fmt.Println()

	// Step 2: Read current settings
	s.Suffix = " Reading current settings..."
	s.Start()

	settings, err := config.ReadSettings()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to read settings: %w", err)
	}

	s.Stop()

	// Step 3: Check if already configured
	if config.IsStatuslineConfigured(settings) {
		fmt.Println("  ccstatus is already configured as the statusline command.")
		fmt.Println("  No changes needed.")
		return nil
	}

	// Check if another statusline is configured
	currentCmd := config.GetStatuslineCommand(settings)
	if currentCmd != "" {
		fmt.Printf("  Current statusline command: %s\n", currentCmd)
		fmt.Println()
	}

	// Step 4: Show what will change
	fmt.Println("The following changes will be made:")
	fmt.Println()

	if exists {
		fmt.Println("  1. Create a backup of your current settings")
	}

	if currentCmd != "" {
		fmt.Printf("  2. Replace statusline command: %s -> ccstatus\n", currentCmd)
	} else {
		fmt.Println("  2. Add statusline configuration:")
		preview := map[string]any{
			"statusline": map[string]string{
				"command": "ccstatus",
			},
		}
		previewJSON, _ := json.MarshalIndent(preview, "     ", "  ")
		fmt.Printf("     %s\n", string(previewJSON))
	}
	fmt.Println()

	// Step 5: Ask for confirmation
	if !confirm("Do you want to proceed?") {
		fmt.Println("Installation cancelled.")
		return nil
	}
	fmt.Println()

	// Step 6: Create backup
	if exists {
		s.Suffix = " Creating backup..."
		s.Start()

		backupPath, err := config.CreateBackup()
		if err != nil {
			s.Stop()
			return fmt.Errorf("failed to create backup: %w", err)
		}

		s.Stop()
		fmt.Printf("  Backup created: %s\n", backupPath)
	}

	// Step 7: Update settings
	s.Suffix = " Updating configuration..."
	s.Start()

	config.SetStatuslineCommand(settings, "ccstatus")

	if err := config.WriteSettings(settings); err != nil {
		s.Stop()
		return fmt.Errorf("failed to write settings: %w", err)
	}

	s.Stop()
	fmt.Println("  Configuration updated successfully.")
	fmt.Println()

	// Step 8: Success message
	fmt.Println("Installation complete!")
	fmt.Println()
	fmt.Println("ccstatus is now configured as your Claude Code statusline.")
	fmt.Println("Restart Claude Code to see the changes.")

	return nil
}

// confirm prompts the user for a yes/no confirmation
func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [y/N]: ", prompt)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
