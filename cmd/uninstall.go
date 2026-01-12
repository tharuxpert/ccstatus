package cmd

import (
	"fmt"
	"time"

	"ccstatus/internal/config"
	"ccstatus/internal/ui"

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
	ui.CompactTitle("ccstatus uninstall")

	// Step 1: Check configuration
	s := ui.NewSpinner("Checking current configuration...")
	s.Start()
	time.Sleep(300 * time.Millisecond)

	exists, err := config.ConfigExists()
	if err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to check config", err.Error())
		return nil
	}

	if !exists {
		s.Stop()
		ui.WarningMessage("No configuration found", "Claude Code settings file does not exist.")
		fmt.Println()
		ui.Dim.Println("  Nothing to uninstall.")
		fmt.Println()
		return nil
	}

	settings, err := config.ReadSettings()
	if err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to read settings", err.Error())
		return nil
	}

	s.Stop()

	// Step 2: Check current statusline configuration
	currentCmd := config.GetStatuslineCommand(settings)

	fmt.Println()
	ui.Bold.Println("  Current Status")
	ui.Divider()
	fmt.Println()

	if currentCmd == "" {
		ui.StatusInfo("Statusline", "Not configured")
		fmt.Println()
		ui.Dim.Println("  Nothing to uninstall.")
		fmt.Println()
		return nil
	}

	if currentCmd != "ccstatus" {
		ui.StatusWarning("Statusline", currentCmd)
		fmt.Println()
		ui.Dim.Println("  ccstatus is not the configured statusline.")
		ui.Dim.Println("  Nothing to uninstall.")
		fmt.Println()
		return nil
	}

	ui.StatusOK("Statusline", "ccstatus (installed)")

	// Step 3: Check for backups
	s = ui.NewSpinner("Checking for backups...")
	s.Start()
	time.Sleep(200 * time.Millisecond)

	backupPath, backupErr := config.GetLatestBackup()

	s.Stop()

	// Step 4: Present options
	fmt.Println()
	ui.Bold.Println("  Uninstall Options")
	ui.Divider()

	if backupErr == nil {
		fmt.Println()
		ui.StatusOK("Backup found", backupPath)

		options := []string{
			"Restore from backup (recommended)",
			"Remove statusline configuration only",
			"Cancel",
		}

		choice := ui.PromptChoice("How would you like to proceed?", options)

		switch choice {
		case 1:
			return restoreFromBackupStyled(backupPath)
		case 2:
			return removeStatuslineStyled(settings)
		case 3:
			fmt.Println()
			ui.WarningMessage("Uninstall cancelled", "No changes were made.")
			fmt.Println()
			return nil
		}
	} else {
		fmt.Println()
		ui.StatusWarning("No backup found", "")
		fmt.Println()
		ui.Dim.Println("  The statusline configuration will be removed.")
		ui.Dim.Println("  Other settings will remain unchanged.")

		if !ui.Confirm("Remove statusline configuration?") {
			fmt.Println()
			ui.WarningMessage("Uninstall cancelled", "No changes were made.")
			fmt.Println()
			return nil
		}

		return removeStatuslineStyled(settings)
	}

	return nil
}

func restoreFromBackupStyled(backupPath string) error {
	if !ui.Confirm("Restore configuration from backup?") {
		fmt.Println()
		ui.WarningMessage("Uninstall cancelled", "No changes were made.")
		fmt.Println()
		return nil
	}

	fmt.Println()
	s := ui.NewProgressSpinner("Restoring from backup...")
	s.Start()
	time.Sleep(400 * time.Millisecond)

	if err := config.RestoreFromBackup(backupPath); err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to restore from backup", err.Error())
		return nil
	}

	s.Stop()
	ui.StatusOK("Configuration restored", "")

	ui.SuccessMessage("Uninstall complete!", "")
	fmt.Println()
	ui.InfoBox(
		"Your previous configuration has been restored.",
		"",
		"Restart Claude Code to see the changes.",
	)
	fmt.Println()

	return nil
}

func removeStatuslineStyled(settings config.Settings) error {
	fmt.Println()

	// Create backup before removing
	s := ui.NewProgressSpinner("Creating safety backup...")
	s.Start()
	time.Sleep(300 * time.Millisecond)

	backupPath, err := config.CreateBackup()
	if err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to create backup", err.Error())
		return nil
	}

	s.Stop()

	if backupPath != "" {
		ui.StatusOK("Safety backup created", backupPath)
	}

	// Remove statusline config
	s = ui.NewProgressSpinner("Removing statusline configuration...")
	s.Start()
	time.Sleep(300 * time.Millisecond)

	config.RemoveStatusline(settings)

	if err := config.WriteSettings(settings); err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to write settings", err.Error())
		return nil
	}

	s.Stop()
	ui.StatusOK("Statusline configuration removed", "")

	ui.SuccessMessage("Uninstall complete!", "")
	fmt.Println()
	ui.InfoBox(
		"ccstatus has been removed from your configuration.",
		"",
		"Restart Claude Code to see the changes.",
	)
	fmt.Println()

	return nil
}
