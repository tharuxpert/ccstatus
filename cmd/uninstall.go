package cmd

import (
	"fmt"
	"time"

	"ccstatus/internal/config"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove ccstatus from Claude Code statusline configuration",
	Long: `Uninstall removes ccstatus from the Claude Code statusline configuration.

This command will:
  1. Check if ccstatus is currently configured
  2. Offer to restore from a previous backup, or
  3. Remove the statusline configuration entirely

You will be asked to confirm before any changes are made.`,
	RunE: runUninstall,
}

func runUninstall(cmd *cobra.Command, args []string) error {
	fmt.Println("ccstatus uninstaller")
	fmt.Println()

	// Step 1: Check if config exists
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Checking Claude Code configuration..."
	s.Start()

	exists, err := config.ConfigExists()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to check config: %w", err)
	}

	if !exists {
		s.Stop()
		fmt.Println("  No Claude Code configuration found.")
		fmt.Println("  Nothing to uninstall.")
		return nil
	}

	settings, err := config.ReadSettings()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to read settings: %w", err)
	}

	s.Stop()

	// Step 2: Check if ccstatus is configured
	currentCmd := config.GetStatuslineCommand(settings)

	if currentCmd == "" {
		fmt.Println("  No statusline is currently configured.")
		fmt.Println("  Nothing to uninstall.")
		return nil
	}

	if currentCmd != "ccstatus" {
		fmt.Printf("  Current statusline command: %s\n", currentCmd)
		fmt.Println("  ccstatus is not the configured statusline.")
		fmt.Println("  Nothing to uninstall.")
		return nil
	}

	fmt.Println("  ccstatus is currently configured as the statusline.")
	fmt.Println()

	// Step 3: Check for available backup
	s.Suffix = " Checking for backups..."
	s.Start()

	backupPath, backupErr := config.GetLatestBackup()

	s.Stop()

	// Step 4: Offer options
	if backupErr == nil {
		fmt.Printf("  Found backup: %s\n", backupPath)
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  1. Restore from backup (recommended)")
		fmt.Println("  2. Remove statusline configuration only")
		fmt.Println("  3. Cancel")
		fmt.Println()

		choice := promptChoice("Select option", []string{"1", "2", "3"})

		switch choice {
		case "1":
			return restoreFromBackup(backupPath)
		case "2":
			return removeStatusline(settings)
		case "3":
			fmt.Println("Uninstall cancelled.")
			return nil
		}
	} else {
		fmt.Println("  No backup found.")
		fmt.Println()
		fmt.Println("The statusline configuration will be removed.")
		fmt.Println("Other settings will remain unchanged.")
		fmt.Println()

		if !confirm("Do you want to proceed?") {
			fmt.Println("Uninstall cancelled.")
			return nil
		}

		return removeStatusline(settings)
	}

	return nil
}

func restoreFromBackup(backupPath string) error {
	fmt.Println()

	if !confirm("Restore configuration from backup?") {
		fmt.Println("Uninstall cancelled.")
		return nil
	}

	fmt.Println()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Restoring from backup..."
	s.Start()

	if err := config.RestoreFromBackup(backupPath); err != nil {
		s.Stop()
		return fmt.Errorf("failed to restore from backup: %w", err)
	}

	s.Stop()

	fmt.Println("  Configuration restored from backup.")
	fmt.Println()
	fmt.Println("Uninstall complete!")
	fmt.Println("Restart Claude Code to see the changes.")

	return nil
}

func removeStatusline(settings config.Settings) error {
	fmt.Println()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Creating backup..."
	s.Start()

	backupPath, err := config.CreateBackup()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to create backup: %w", err)
	}

	s.Stop()

	if backupPath != "" {
		fmt.Printf("  Backup created: %s\n", backupPath)
	}

	s.Suffix = " Removing statusline configuration..."
	s.Start()

	config.RemoveStatusline(settings)

	if err := config.WriteSettings(settings); err != nil {
		s.Stop()
		return fmt.Errorf("failed to write settings: %w", err)
	}

	s.Stop()

	fmt.Println("  Statusline configuration removed.")
	fmt.Println()
	fmt.Println("Uninstall complete!")
	fmt.Println("Restart Claude Code to see the changes.")

	return nil
}

// promptChoice prompts the user to select from a list of valid choices
func promptChoice(prompt string, validChoices []string) string {
	for {
		fmt.Printf("%s [%s]: ", prompt, joinChoices(validChoices))

		var response string
		fmt.Scanln(&response)

		for _, choice := range validChoices {
			if response == choice {
				return response
			}
		}

		fmt.Printf("  Invalid choice. Please enter one of: %s\n", joinChoices(validChoices))
	}
}

func joinChoices(choices []string) string {
	result := ""
	for i, c := range choices {
		if i > 0 {
			result += "/"
		}
		result += c
	}
	return result
}
