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
var configFile string
var appType string
var useWizard bool
var moduleName string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project",
	Long: `Create a new Go project with a structured layout.
Launches an interactive wizard to configure your project,
or uses default settings if you skip the wizard.

You can also specify a configuration file with --config
or a project type with --type (cli, api, library).`,
	Args: cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		// Initialize config based on provided options
		if configFile != "" {
			// Load config from file
			var err error
			projectConfig, err = config.LoadConfigFromFile(configFile)
			if err != nil {
				fmt.Printf("Error loading config file: %v\n", err)
				return
			}
			fmt.Printf("Loaded configuration from %s\n", configFile)
		} else if appType != "" {
			// Initialize config based on project type
			switch appType {
			case string(config.TypeCLI):
				projectConfig = config.NewCLIProjectConfig()
			case string(config.TypeAPI):
				projectConfig = config.NewAPIProjectConfig()
			case string(config.TypeLibrary):
				projectConfig = config.NewLibraryProjectConfig()
			default:
				fmt.Printf("Unknown project type: %s. Using default.\n", appType)
				projectConfig = config.NewDefaultProjectConfig()
			}
			fmt.Printf("Using %s project template\n", appType)
		} else {
			// Initialize default config
			projectConfig = config.NewDefaultProjectConfig()
		}

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
	newCmd.Flags().StringVarP(&configFile, "config", "c", "", "path to configuration file")
	newCmd.Flags().StringVarP(&appType, "type", "t", "", "project type (cli, api, library)")
	newCmd.Flags().BoolVarP(&useWizard, "wizard", "w", true, "use interactive wizard")
	newCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Go module name")
}
