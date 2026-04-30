package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"ccstatus/internal/ui"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the ccstatus version number and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		ui.Primary.Printf("  ccstatus ")
		ui.Bold.Println(DisplayVersion())
		ui.Dim.Printf("  %s/%s\n", runtime.GOOS, runtime.GOARCH)
		if commit != "" {
			ui.Dim.Printf("  commit %s\n", commit)
		}
		if date != "" {
			ui.Dim.Printf("  built %s\n", date)
		}
		fmt.Println()
	},
}

// GetVersion returns the version without a leading v prefix.
func GetVersion() string {
	if version != "" && version != "dev" {
		return strings.TrimPrefix(version, "v")
	}

	if described := gitDescribeVersion(); described != "" {
		return strings.TrimPrefix(described, "v")
	}

	return "dev"
}

// DisplayVersion returns the user-facing version string.
func DisplayVersion() string {
	current := GetVersion()
	if current == "" || current == "dev" {
		return "dev"
	}
	return "v" + strings.TrimPrefix(current, "v")
}

func gitDescribeVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
