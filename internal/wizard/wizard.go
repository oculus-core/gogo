package wizard

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/charmbracelet/lipgloss"

	"github.com/oculus-core/gogo/pkg/config"
)

var (
	// Define some styles with muted colors
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#8A7B9D")). // Muted purple/mauve
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#6B8E6B")). // Muted sage green
			MarginTop(1).
			MarginBottom(1)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#778899")) // Light slate gray
)

// RunWizard runs the interactive project setup wizard
func RunWizard(cfg *config.ProjectConfig) error {
	fmt.Println() // Add blank line before the welcome banner
	fmt.Println(titleStyle.Render("🚀 Welcome to the Gogo Project Generator Wizard"))
	fmt.Println("This wizard will help you set up a new Go project with best practices")
	fmt.Println()

	// Project information section
	fmt.Println(sectionStyle.Render("📋 Project Information"))

	// Project name
	namePrompt := &survey.Input{
		Message: "Project name:",
		Default: cfg.Name,
	}
	if err := survey.AskOne(namePrompt, &cfg.Name); err != nil {
		if err == terminal.InterruptErr {
			return fmt.Errorf("wizard cancelled")
		}
		return err
	}

	// Module path
	modulePrompt := &survey.Input{
		Message: "Module path:",
		Default: cfg.Module,
	}
	if err := survey.AskOne(modulePrompt, &cfg.Module); err != nil {
		if err == terminal.InterruptErr {
			return fmt.Errorf("wizard cancelled")
		}
		return err
	}

	// Description
	descPrompt := &survey.Input{
		Message: "Description:",
		Default: cfg.Description,
	}
	if err := survey.AskOne(descPrompt, &cfg.Description); err != nil {
		if err == terminal.InterruptErr {
			return fmt.Errorf("wizard cancelled")
		}
		return err
	}

	// Author
	authorPrompt := &survey.Input{
		Message: "Author:",
		Default: cfg.Author,
	}
	if err := survey.AskOne(authorPrompt, &cfg.Author); err != nil {
		if err == terminal.InterruptErr {
			return fmt.Errorf("wizard cancelled")
		}
		return err
	}

	// License
	licensePrompt := &survey.Select{
		Message: "License:",
		Options: []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "None"},
		Default: cfg.License,
	}
	if err := survey.AskOne(licensePrompt, &cfg.License); err != nil {
		if err == terminal.InterruptErr {
			return fmt.Errorf("wizard cancelled")
		}
		return err
	}

	// Now ask for project details using survey
	fmt.Println(highlightStyle.Render("\nProject Details:"))

	// Project Type
	appTypePrompt := &survey.Select{
		Message: "Project Type:",
		Options: []string{
			string(config.TypeDefault),
			string(config.TypeCLI),
			string(config.TypeAPI),
			string(config.TypeLibrary),
		},
		Default: string(cfg.Type),
		Description: func(value string, _ int) string {
			switch value {
			case string(config.TypeCLI):
				return "Command-line application (includes Cobra and Viper)"
			case string(config.TypeAPI):
				return "API/Web service (includes Gin)"
			case string(config.TypeLibrary):
				return "Library/Package (no cmd directory)"
			default:
				return "Generic Go project"
			}
		},
	}

	var appTypeStr string
	if err := survey.AskOne(appTypePrompt, &appTypeStr); err != nil {
		if err == terminal.InterruptErr {
			return fmt.Errorf("wizard cancelled")
		}
		return err
	}

	// Update the config based on the selected project type
	prevType := cfg.Type
	cfg.Type = config.ProjectType(appTypeStr)

	// Only apply type-specific settings if the type has changed
	if prevType != cfg.Type {
		switch cfg.Type {
		case config.TypeCLI:
			cfg.UseCobra = true
			cfg.UseViper = true
		case config.TypeAPI:
			cfg.UseGin = true
		case config.TypeLibrary:
			cfg.UseCmd = false
		}
	}

	// Project structure section
	fmt.Println(sectionStyle.Render("📁 Project Structure"))

	structurePrompt := &survey.MultiSelect{
		Message: "Select project directories to include:",
		Options: []string{
			"cmd (application entrypoints)",
			"internal (private packages)",
			"pkg (public packages)",
			"test (test utilities)",
			"docs (documentation)",
		},
		Default: getStructureDefaults(cfg),
	}

	var selectedStructure []string
	if err := survey.AskOne(structurePrompt, &selectedStructure); err != nil {
		return err
	}

	// Update config based on selections
	cfg.UseCmd = contains(selectedStructure, "cmd (application entrypoints)")
	cfg.UseInternal = contains(selectedStructure, "internal (private packages)")
	cfg.UsePkg = contains(selectedStructure, "pkg (public packages)")
	cfg.UseTest = contains(selectedStructure, "test (test utilities)")
	cfg.UseDocs = contains(selectedStructure, "docs (documentation)")

	// Files section
	fmt.Println(sectionStyle.Render("📝 Project Files"))

	filesPrompt := &survey.MultiSelect{
		Message: "Select files to generate:",
		Options: []string{
			"README.md",
			"LICENSE",
			"Makefile",
		},
		Default: getFilesDefaults(cfg),
	}

	var selectedFiles []string
	if err := survey.AskOne(filesPrompt, &selectedFiles); err != nil {
		return err
	}

	// Update config based on selections
	cfg.CreateReadme = contains(selectedFiles, "README.md")
	cfg.CreateLicense = contains(selectedFiles, "LICENSE")
	cfg.CreateMakefile = contains(selectedFiles, "Makefile")

	// Code quality tools section
	fmt.Println(sectionStyle.Render("🛠️ Code Quality Tools"))

	toolsPrompt := &survey.MultiSelect{
		Message: "Select code quality tools to include:",
		Options: []string{
			"Linters (golangci-lint)",
			"Pre-commit hooks",
			"Git hooks",
		},
		Default: getToolsDefaults(cfg),
	}

	var selectedTools []string
	if err := survey.AskOne(toolsPrompt, &selectedTools); err != nil {
		return err
	}

	// Update config based on selections
	cfg.UseLinters = contains(selectedTools, "Linters (golangci-lint)")
	cfg.UsePreCommitHooks = contains(selectedTools, "Pre-commit hooks")
	cfg.UseGitHooks = contains(selectedTools, "Git hooks")

	// Dependencies section
	fmt.Println(sectionStyle.Render("📦 Dependencies"))

	depsPrompt := &survey.MultiSelect{
		Message: "Select dependencies to include:",
		Options: []string{
			"Cobra (CLI framework)",
			"Viper (configuration)",
		},
		Default: getDepsDefaults(cfg),
	}

	var selectedDeps []string
	if err := survey.AskOne(depsPrompt, &selectedDeps); err != nil {
		return err
	}

	// Update config based on selections
	cfg.UseCobra = contains(selectedDeps, "Cobra (CLI framework)")
	cfg.UseViper = contains(selectedDeps, "Viper (configuration)")

	// CI/CD section
	fmt.Println(sectionStyle.Render("🔄 CI/CD"))

	cicdPrompt := &survey.Confirm{
		Message: "Set up GitHub Actions for CI/CD?",
		Default: cfg.UseGitHubActions,
	}
	if err := survey.AskOne(cicdPrompt, &cfg.UseGitHubActions); err != nil {
		return err
	}

	// Summary
	fmt.Println(sectionStyle.Render("✅ Configuration Summary"))
	fmt.Println(highlightStyle.Render("Project:"), cfg.Name)
	fmt.Println(highlightStyle.Render("Module:"), cfg.Module)
	fmt.Println(highlightStyle.Render("Description:"), cfg.Description)
	fmt.Println(highlightStyle.Render("Author:"), cfg.Author)
	fmt.Println(highlightStyle.Render("License:"), cfg.License)

	fmt.Println(highlightStyle.Render("Directories:"))
	if cfg.UseCmd {
		fmt.Println("  - cmd")
	}
	if cfg.UseInternal {
		fmt.Println("  - internal")
	}
	if cfg.UsePkg {
		fmt.Println("  - pkg")
	}
	if cfg.UseTest {
		fmt.Println("  - test")
	}
	if cfg.UseDocs {
		fmt.Println("  - docs")
	}

	fmt.Println(highlightStyle.Render("Files:"))
	if cfg.CreateReadme {
		fmt.Println("  - README.md")
	}
	if cfg.CreateLicense {
		fmt.Println("  - LICENSE")
	}
	if cfg.CreateMakefile {
		fmt.Println("  - Makefile")
	}

	fmt.Println(highlightStyle.Render("Tools:"))
	if cfg.UseLinters {
		fmt.Println("  - Linters")
	}
	if cfg.UsePreCommitHooks {
		fmt.Println("  - Pre-commit hooks")
	}
	if cfg.UseGitHooks {
		fmt.Println("  - Git hooks")
	}

	fmt.Println(highlightStyle.Render("Dependencies:"))
	if cfg.UseCobra {
		fmt.Println("  - Cobra")
	}
	if cfg.UseViper {
		fmt.Println("  - Viper")
	}

	fmt.Println(highlightStyle.Render("CI/CD:"))
	if cfg.UseGitHubActions {
		fmt.Println("  - GitHub Actions")
	}

	// Confirm generation
	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "Generate project with these settings?",
		Default: true,
	}
	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("project generation cancelled")
	}

	return nil
}

// Helper functions to set default selections in the wizard
func getStructureDefaults(cfg *config.ProjectConfig) []string {
	var defaults []string
	if cfg.UseCmd {
		defaults = append(defaults, "cmd (application entrypoints)")
	}
	if cfg.UseInternal {
		defaults = append(defaults, "internal (private packages)")
	}
	if cfg.UsePkg {
		defaults = append(defaults, "pkg (public packages)")
	}
	if cfg.UseTest {
		defaults = append(defaults, "test (test utilities)")
	}
	if cfg.UseDocs {
		defaults = append(defaults, "docs (documentation)")
	}
	return defaults
}

func getFilesDefaults(cfg *config.ProjectConfig) []string {
	var defaults []string
	if cfg.CreateReadme {
		defaults = append(defaults, "README.md")
	}
	if cfg.CreateLicense {
		defaults = append(defaults, "LICENSE")
	}
	if cfg.CreateMakefile {
		defaults = append(defaults, "Makefile")
	}
	return defaults
}

func getToolsDefaults(cfg *config.ProjectConfig) []string {
	var defaults []string
	if cfg.UseLinters {
		defaults = append(defaults, "Linters (golangci-lint)")
	}
	if cfg.UsePreCommitHooks {
		defaults = append(defaults, "Pre-commit hooks")
	}
	if cfg.UseGitHooks {
		defaults = append(defaults, "Git hooks")
	}
	return defaults
}

func getDepsDefaults(cfg *config.ProjectConfig) []string {
	var defaults []string
	if cfg.UseCobra {
		defaults = append(defaults, "Cobra (CLI framework)")
	}
	if cfg.UseViper {
		defaults = append(defaults, "Viper (configuration)")
	}
	return defaults
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == strings.TrimSpace(item) {
			return true
		}
	}
	return false
}
