// Package cmd contains all CLI commands for ccstatus.
package cmd

import (
	"os"

	"ccstatus/internal/statusline"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags
var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ccstatus",
	Short: "Claude Code statusline utility",
	Long: `ccstatus is a statusline utility for Claude Code.

When run without arguments, it outputs the current usage status
for use in Claude Code's statusline feature.

Use subcommands for installation and diagnostics.`,
	// Run the statusline output when no subcommand is provided
	Run: func(cmd *cobra.Command, args []string) {
		statusline.Run()
	},
	// Disable completion command
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	// Silence usage on errors for cleaner output
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Disable the default help command to keep CLI clean
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	// Add subcommands
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(versionCmd)
}
