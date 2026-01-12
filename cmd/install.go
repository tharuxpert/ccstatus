package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"ccstatus/internal/config"
	"ccstatus/internal/ui"

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
	ui.CompactTitle("ccstatus install")

	// Step 1: Check configuration
	s := ui.NewSpinner("Detecting Claude Code configuration...")
	s.Start()
	time.Sleep(300 * time.Millisecond) // Brief pause for visual feedback

	configPath, err := config.GetConfigPath()
	if err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to determine config path", err.Error())
		return nil
	}

	exists, err := config.ConfigExists()
	if err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to check config", err.Error())
		return nil
	}

	s.Stop()

	// Show config status
	fmt.Println()
	ui.Bold.Println("  Configuration")
	ui.Divider()
	fmt.Println()

	ui.PrintPath("Location", configPath)
	if exists {
		ui.StatusOK("Config file", "Found")
	} else {
		ui.StatusInfo("Config file", "Will be created")
	}

	// Step 2: Read current settings
	s = ui.NewSpinner("Reading current settings...")
	s.Start()
	time.Sleep(200 * time.Millisecond)

	settings, err := config.ReadSettings()
	if err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to read settings", err.Error())
		return nil
	}

	s.Stop()

	// Step 3: Check if already configured
	if config.IsStatuslineConfigured(settings) {
		ui.SuccessMessage("Already configured!", "ccstatus is already set as the statusline command.")
		fmt.Println()
		ui.Dim.Println("  No changes needed.")
		fmt.Println()
		return nil
	}

	// Check if another statusline is configured
	currentCmd := config.GetStatuslineCommand(settings)
	if currentCmd != "" {
		fmt.Println()
		ui.StatusWarning("Existing configuration", "")
		ui.PrintKeyValue("Current command", currentCmd)
	}

	// Step 4: Show what will change
	fmt.Println()
	ui.Bold.Println("  Changes to be made")
	ui.Divider()
	fmt.Println()

	stepNum := 1
	if exists {
		ui.Step(stepNum, "Create a backup of current settings")
		stepNum++
	}

	if currentCmd != "" {
		// Another command is configured, show replacement
		ui.Step(stepNum, fmt.Sprintf("Update statusLine.command: %s %s %s",
			ui.Error.Sprint(currentCmd),
			ui.Dim.Sprint(ui.IconArrow),
			ui.Success.Sprint("ccstatus")))
	} else if config.HasStatuslineObject(settings) {
		// statusLine exists but command is empty/missing
		ui.Step(stepNum, "Set statusLine.command:")
		fmt.Println()
		ui.Dim.Println("     Existing statusLine object will be preserved.")
		ui.Dim.Print("     Setting: ")
		ui.Info.Println("\"command\": \"ccstatus\"")
	} else {
		// No statusLine object, will create new one
		ui.Step(stepNum, "Add statusLine configuration:")
		fmt.Println()
		preview := map[string]any{
			"statusLine": map[string]string{
				"command": "ccstatus",
			},
		}
		previewJSON, _ := json.MarshalIndent(preview, "", "  ")
		ui.CodeBlock(string(previewJSON))
	}

	// Step 5: Ask for confirmation
	if !ui.Confirm("Apply these changes?") {
		fmt.Println()
		ui.WarningMessage("Installation cancelled", "No changes were made.")
		fmt.Println()
		return nil
	}

	// Step 6: Create backup
	fmt.Println()
	if exists {
		s = ui.NewProgressSpinner("Creating backup...")
		s.Start()
		time.Sleep(300 * time.Millisecond)

		backupPath, err := config.CreateBackup()
		if err != nil {
			s.Stop()
			ui.ErrorMessage("Failed to create backup", err.Error())
			return nil
		}

		s.Stop()
		ui.StatusOK("Backup created", backupPath)
	}

	// Step 7: Update settings
	s = ui.NewProgressSpinner("Updating configuration...")
	s.Start()
	time.Sleep(300 * time.Millisecond)

	config.SetStatuslineCommand(settings, "ccstatus")

	if err := config.WriteSettings(settings); err != nil {
		s.Stop()
		ui.ErrorMessage("Failed to write settings", err.Error())
		return nil
	}

	s.Stop()
	ui.StatusOK("Configuration updated", "")

	// Step 8: Success message
	ui.SuccessMessage("Installation complete!", "")
	fmt.Println()
	ui.InfoBox(
		"ccstatus is now your Claude Code statusline.",
		"",
		"Restart Claude Code to see the changes.",
	)
	fmt.Println()

	return nil
}
