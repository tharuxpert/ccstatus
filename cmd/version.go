package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the ccstatus version number and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ccstatus version %s\n", Version)
	},
}
