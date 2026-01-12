package cmd

import (
	"fmt"
	"runtime"

	"ccstatus/internal/ui"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the ccstatus version number and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		ui.Primary.Printf("  ccstatus ")
		ui.Bold.Printf("v%s\n", Version)
		ui.Dim.Printf("  %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println()
	},
}
