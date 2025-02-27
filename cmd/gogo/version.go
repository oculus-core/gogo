package gogo

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information - will be set during build via ldflags
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Long:  `Display version, commit, and build date information.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Gogo CLI")
		fmt.Println("--------")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Commit:     %s\n", Commit)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
