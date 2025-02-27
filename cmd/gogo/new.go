package gogo

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/oculus-core/gogo/internal/wizard"
	"github.com/oculus-core/gogo/pkg/config"
)

var projectConfig *config.ProjectConfig
var outputDir string
var skipWizard bool

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project",
	Long: `Create a new Go project with a structured layout.
Launches an interactive wizard to configure your project,
or uses default settings if you skip the wizard.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		// Initialize default config
		projectConfig = config.NewDefaultProjectConfig()

		// If a project name is provided, use it
		if len(args) > 0 {
			projectConfig.Name = args[0]
		}

		if !skipWizard {
			// Run the interactive wizard
			if err := wizard.RunWizard(projectConfig); err != nil {
				fmt.Printf("Error in wizard: %v\n", err)
				return
			}
		}

		// Generate the project
		if err := wizard.GenerateProject(projectConfig, outputDir); err != nil {
			fmt.Printf("Error generating project: %v\n", err)
			return
		}

		// Get absolute path for display
		absPath, err := filepath.Abs(outputDir)
		if err != nil {
			// Fallback to the relative path if there's an error
			absPath = outputDir
		}

		fmt.Printf("\nSuccessfully created project %s in %s\n", projectConfig.Name, absPath)
		fmt.Println("\nNext steps:")
		fmt.Println("  1. cd", outputDir)
		fmt.Println("  2. git init")
		fmt.Println("  3. go mod tidy")
		fmt.Println("  4. make build")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Flags for the new command
	newCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "output directory for the project")
	newCmd.Flags().BoolVarP(&skipWizard, "skip-wizard", "s", false, "skip the interactive wizard and use defaults")
}
