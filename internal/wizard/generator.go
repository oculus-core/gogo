package wizard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oculus-core/gogo/pkg/config"
)

// GenerateProject creates a new Go project based on the provided configuration
func GenerateProject(cfg *config.ProjectConfig, outputDir string) error {
	// Create project directory if it doesn't exist
	projectDir := filepath.Join(outputDir, cfg.Name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	// Generate root files
	if err := generateRootFiles(cfg, projectDir); err != nil {
		return err
	}

	// Create standard directory structure
	dirs := []string{}

	if cfg.UseCmd {
		dirs = append(dirs, "cmd")
	}

	if cfg.UseInternal {
		dirs = append(dirs, "internal")
	}

	if cfg.UsePkg {
		dirs = append(dirs, "pkg")
	}

	if cfg.UseDocs {
		dirs = append(dirs, "docs")
	}

	if cfg.UseTest {
		dirs = append(dirs, "test")
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create .github directory and workflows if needed
	if cfg.UseGitHubActions {
		if err := generateGitHubWorkflows(cfg, projectDir); err != nil {
			return err
		}
	}

	// Generate Go module
	if err := generateGoMod(cfg, projectDir); err != nil {
		return err
	}

	// Generate initial code
	if err := generateInitialCode(cfg, projectDir); err != nil {
		return err
	}

	// Generate templates
	if err := generateTemplates(cfg, projectDir); err != nil {
		return err
	}

	// Generate linter configuration if needed
	if cfg.UseLinters {
		if err := generateLinterConfig(cfg, projectDir); err != nil {
			return err
		}
	}

	// Generate pre-commit hooks if needed
	if cfg.UsePreCommitHooks {
		if err := generatePreCommitConfig(cfg, projectDir); err != nil {
			return err
		}
	}

	return nil
}

// generateRootFiles creates the basic files at the project root
func generateRootFiles(cfg *config.ProjectConfig, projectDir string) error {
	// Generate README.md
	if cfg.CreateReadme {
		readmePath := filepath.Join(projectDir, "README.md")

		// Fix: Split the string format to avoid backtick issues
		readmeContent := fmt.Sprintf("# %s\n\n%s\n\n## Overview\n\nTODO: Add project overview\n\n## Installation\n\n### Prerequisites\n\n- Go 1.16 or later\n\n### Building from Source\n\n", cfg.Name, cfg.Description)

		// Add code block separately to avoid backtick issues
		readmeContent += "```bash\n"
		readmeContent += fmt.Sprintf("# Clone the repository\ngit clone %s.git\ncd %s\n\n# Build the binary\ngo build -o bin/%s\n\n# Run tests\ngo test ./...\n", cfg.Module, cfg.Name, strings.ToLower(cfg.Name))
		readmeContent += "```\n\n"

		if cfg.CreateMakefile {
			readmeContent += "## Using Make\n\nThe project includes a Makefile to simplify common tasks:\n\n```bash\n"
			readmeContent += "# Build the binary\nmake build\n\n# Run tests\nmake test\n\n# Clean build artifacts\nmake clean\n"
			readmeContent += "```\n\nFor more details, run `make help` to see all available commands.\n"
		}

		if err := os.WriteFile(readmePath, []byte(readmeContent), 0600); err != nil {
			return err
		}
	}

	// Generate LICENSE
	if cfg.CreateLicense && cfg.License != "None" {
		licensePath := filepath.Join(projectDir, "LICENSE")
		year := time.Now().Year()

		var licenseContent string
		switch cfg.License {
		case "MIT":
			licenseContent = fmt.Sprintf("MIT License\n\nCopyright (c) %d %s\n\n"+
				"Permission is hereby granted, free of charge, to any person obtaining a copy\n"+
				"of this software and associated documentation files (the \"Software\"), to deal\n"+
				"in the Software without restriction, including without limitation the rights\n"+
				"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell\n"+
				"copies of the Software, and to permit persons to whom the Software is\n"+
				"furnished to do so, subject to the following conditions:\n\n"+
				"The above copyright notice and this permission notice shall be included in all\n"+
				"copies or substantial portions of the Software.\n\n"+
				"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\n"+
				"IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\n"+
				"FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\n"+
				"AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\n"+
				"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\n"+
				"OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\n"+
				"SOFTWARE.\n", year, cfg.Author)

		// Add more license templates as needed for other license types
		default:
			licenseContent = fmt.Sprintf("Copyright (c) %d %s\n\n"+
				"This project is licensed under the %s License.\n"+
				"Please see https://choosealicense.com/licenses/ for more information.\n",
				year, cfg.Author, cfg.License)
		}

		if err := os.WriteFile(licensePath, []byte(licenseContent), 0600); err != nil {
			return err
		}
	}

	// Generate .gitignore
	gitignorePath := filepath.Join(projectDir, ".gitignore")
	gitignoreContent := "# Binaries for programs and plugins\n" +
		"*.exe\n" +
		"*.exe~\n" +
		"*.dll\n" +
		"*.so\n" +
		"*.dylib\n" +
		"bin/\n\n" +
		"# Test binary, built with 'go test -c'\n" +
		"*.test\n\n" +
		"# Output of the go coverage tool\n" +
		"*.out\n" +
		"coverage.html\n\n" +
		"# Dependency directories (remove the comment below to include it)\n" +
		"# vendor/\n\n" +
		"# Go workspace file\n" +
		"go.work\n\n" +
		"# IDE specific files\n" +
		".idea/\n" +
		".vscode/\n" +
		"*.swp\n" +
		"*.swo\n\n" +
		"# OS specific files\n" +
		".DS_Store\n" +
		".DS_Store?\n" +
		"._*\n" +
		".Spotlight-V100\n" +
		".Trashes\n" +
		"ehthumbs.db\n" +
		"Thumbs.db\n"

	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0600); err != nil {
		return err
	}

	// Generate Makefile
	if cfg.CreateMakefile {
		makefilePath := filepath.Join(projectDir, "Makefile")
		makefileContent := fmt.Sprintf(".PHONY: all build clean test\n\n"+
			"# Binary name\n"+
			"BINARY_NAME=%s\n"+
			"# Binary directory\n"+
			"BIN_DIR=./bin\n\n"+
			"# Go commands\n"+
			"GO ?= go\n"+
			"GOBUILD = $(GO) build\n"+
			"GOCLEAN = $(GO) clean\n"+
			"GOTEST = $(GO) test\n"+
			"GOGET = $(GO) get\n\n"+
			"# Version info from git\n"+
			"GIT_COMMIT=$(shell git rev-parse --short HEAD || echo \"unknown\")\n"+
			"GIT_DIRTY=$(shell test -n \"`git status --porcelain`\" && echo \"+DIRTY\" || echo \"\")\n"+
			"GIT_TAG=$(shell git describe --tags --abbrev=0 2>/dev/null || echo \"v0.0.0\")\n"+
			"BUILD_DATE=$(shell date '+%%Y-%%m-%%d-%%H:%%M:%%S')\n\n"+
			"# Get the module name from go.mod\n"+
			"MODULE_NAME=$(shell grep \"^module\" go.mod | awk '{print $$2}')\n\n"+
			"# Linker flags\n"+
			"LDFLAGS=-ldflags \"-X $(MODULE_NAME)/cmd.Version=$(GIT_TAG) \\\n"+
			"-X $(MODULE_NAME)/cmd.Commit=$(GIT_COMMIT)$(GIT_DIRTY) \\\n"+
			"-X $(MODULE_NAME)/cmd.BuildDate=$(BUILD_DATE)\"\n\n"+
			"# Default target (build binary)\n"+
			"all: build\n\n"+
			"# Build binary\n"+
			"build:\n"+
			"\t@echo \"Building $(BINARY_NAME)...\"\n"+
			"\t@echo \"Git commit: $(GIT_COMMIT)$(GIT_DIRTY)\"\n"+
			"\t@echo \"Git tag: $(GIT_TAG)\"\n"+
			"\t@echo \"Build date: $(BUILD_DATE)\"\n"+
			"\t@mkdir -p $(BIN_DIR)\n"+
			"\t$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)\n"+
			"\t@echo \"Build complete: $(BIN_DIR)/$(BINARY_NAME)\"\n\n"+
			"# Clean build artifacts\n"+
			"clean:\n"+
			"\t@echo \"Cleaning...\"\n"+
			"\t@$(GOCLEAN)\n"+
			"\t@rm -rf $(BIN_DIR)\n"+
			"\t@rm -f coverage.out coverage.html\n"+
			"\t@echo \"Clean complete\"\n\n"+
			"# Run tests\n"+
			"test:\n"+
			"\t@echo \"Running tests...\"\n"+
			"\t$(GOTEST) -v ./...\n"+
			"\t@echo \"Tests complete\"\n\n"+
			"# Run tests with coverage\n"+
			"test-coverage:\n"+
			"\t@echo \"Running tests with coverage...\"\n"+
			"\t$(GOTEST) -v ./... -coverprofile=coverage.out\n"+
			"\t$(GO) tool cover -html=coverage.out -o coverage.html\n"+
			"\t@echo \"Coverage report generated at coverage.html\"\n\n"+
			"# Install dependencies\n"+
			"deps:\n"+
			"\t@echo \"Installing dependencies...\"\n"+
			"\t$(GOGET) -v ./...\n"+
			"\t@echo \"Dependencies installed\"\n\n"+
			"# Lint the code\n"+
			"lint:\n"+
			"\t@echo \"Linting code...\"\n"+
			"\tgolangci-lint run ./...\n"+
			"\t@echo \"Lint complete\"\n\n"+
			"# Help target\n"+
			"help:\n"+
			"\t@echo \"Available targets:\"\n"+
			"\t@echo \"  all               - Default target, builds the binary\"\n"+
			"\t@echo \"  build             - Build the binary to $(BIN_DIR)/$(BINARY_NAME)\"\n"+
			"\t@echo \"  clean             - Clean build artifacts\"\n"+
			"\t@echo \"  test              - Run tests\"\n"+
			"\t@echo \"  test-coverage     - Run tests with coverage reporting\"\n"+
			"\t@echo \"  deps              - Install dependencies\"\n"+
			"\t@echo \"  lint              - Lint the code\"\n",
			strings.ToLower(cfg.Name))

		if err := os.WriteFile(makefilePath, []byte(makefileContent), 0600); err != nil {
			return err
		}
	}

	return nil
}

// generateGoMod creates the go.mod file
func generateGoMod(cfg *config.ProjectConfig, projectDir string) error {
	goModPath := filepath.Join(projectDir, "go.mod")
	goModContent := fmt.Sprintf("module %s\n\ngo 1.19\n", cfg.Module)

	if cfg.UseCobra || cfg.UseViper {
		goModContent += "\nrequire (\n"
		if cfg.UseCobra {
			goModContent += "\tgithub.com/spf13/cobra v1.9.1\n"
		}
		if cfg.UseViper {
			goModContent += "\tgithub.com/spf13/viper v1.19.0\n"
		}
		goModContent += ")\n"
	}

	return os.WriteFile(goModPath, []byte(goModContent), 0600)
}

// generateGitHubWorkflows creates GitHub Actions workflow files
func generateGitHubWorkflows(cfg *config.ProjectConfig, projectDir string) error {
	workflowDir := filepath.Join(projectDir, ".github", "workflows")

	// CI workflow
	ciWorkflowPath := filepath.Join(workflowDir, "ci.yml")
	ciWorkflowContent := "name: CI\n\n" +
		"on:\n" +
		"  push:\n" +
		"    branches: [ main ]\n" +
		"  pull_request:\n" +
		"    branches: [ main ]\n\n" +
		"jobs:\n" +
		"  build:\n" +
		"    runs-on: ubuntu-latest\n" +
		"    steps:\n" +
		"    - uses: actions/checkout@v3\n\n" +
		"    - name: Set up Go\n" +
		"      uses: actions/setup-go@v4\n" +
		"      with:\n" +
		"        go-version: '1.19'\n\n" +
		"    - name: Build\n" +
		"      run: go build -v ./...\n\n" +
		"    - name: Test\n" +
		"      run: go test -v ./...\n"

	if err := os.WriteFile(ciWorkflowPath, []byte(ciWorkflowContent), 0600); err != nil {
		return err
	}

	// Lint workflow
	if cfg.UseLinters {
		lintWorkflowPath := filepath.Join(workflowDir, "lint.yml")
		lintWorkflowContent := "name: Lint\n\n" +
			"on:\n" +
			"  push:\n" +
			"    branches: [ main ]\n" +
			"  pull_request:\n" +
			"    branches: [ main ]\n\n" +
			"jobs:\n" +
			"  golangci:\n" +
			"    name: lint\n" +
			"    runs-on: ubuntu-latest\n" +
			"    steps:\n" +
			"      - uses: actions/checkout@v3\n" +
			"      - name: golangci-lint\n" +
			"        uses: golangci/golangci-lint-action@v3\n" +
			"        with:\n" +
			"          version: latest\n"

		if err := os.WriteFile(lintWorkflowPath, []byte(lintWorkflowContent), 0600); err != nil {
			return err
		}
	}

	return nil
}

// generateLinterConfig creates the golangci-lint configuration
func generateLinterConfig(cfg *config.ProjectConfig, projectDir string) error {
	linterConfigPath := filepath.Join(projectDir, ".golangci.yml")
	linterConfigContent := "run:\n" +
		"  timeout: 5m\n" +
		"linters:\n" +
		"  disable-all: true\n" +
		"  enable:\n" +
		"    - errcheck\n" +
		"    - gosimple\n" +
		"    - govet\n" +
		"    - ineffassign\n" +
		"    - staticcheck\n" +
		"    - unused\n" +
		"    - gofmt\n" +
		"    - goimports\n" +
		"    - gosec\n" +
		"    - misspell\n" +
		"    - revive\n" +
		"    - unused\n" +
		"    - whitespace\n" +
		"linters-settings:\n" +
		"  goimports:\n" +
		"    local-prefixes: " + cfg.Module + "\n" +
		"issues:\n" +
		"  exclude-rules:\n" +
		"    - path: _test\\.go\n" +
		"      linters:\n" +
		"        - gosec\n"

	return os.WriteFile(linterConfigPath, []byte(linterConfigContent), 0600)
}

// generatePreCommitConfig creates the pre-commit hooks configuration
func generatePreCommitConfig(cfg *config.ProjectConfig, projectDir string) error {
	preCommitConfigPath := filepath.Join(projectDir, ".pre-commit-config.yaml")

	// Create a custom hook name based on the project name
	customHookName := strings.ToLower(cfg.Name) + "-build"

	preCommitConfigContent := "repos:\n" +
		"  - repo: https://github.com/pre-commit/pre-commit-hooks\n" +
		"    rev: v4.5.0\n" +
		"    hooks:\n" +
		"      - id: trailing-whitespace\n" +
		"      - id: end-of-file-fixer\n" +
		"      - id: check-yaml\n" +
		"      - id: check-added-large-files\n" +
		"      - id: check-json\n" +
		"      - id: check-merge-conflict\n" +
		"  # Commit message validation for conventional commits\n" +
		"  - repo: https://github.com/compilerla/conventional-pre-commit\n" +
		"    rev: v2.1.1\n" +
		"    hooks:\n" +
		"      - id: conventional-pre-commit\n" +
		"        stages: [commit-msg]\n" +
		"        args: [] # Add custom args here if needed\n" +
		"  # Primary Go linting and formatting\n" +
		"  - repo: https://github.com/golangci/golangci-lint\n" +
		"    rev: v1.64.5\n" +
		"    hooks:\n" +
		"      - id: golangci-lint\n" +
		"        args: [--timeout=5m]\n" +
		"  # Additional Go tools\n" +
		"  - repo: https://github.com/dnephin/pre-commit-golang\n" +
		"    rev: v0.5.1\n" +
		"    hooks:\n" +
		"      - id: go-fmt\n" +
		"      - id: go-mod-tidy\n" +
		"      - id: go-unit-tests\n" +
		"  # Local hooks\n" +
		"  - repo: local\n" +
		"    hooks:\n" +
		"      - id: " + customHookName + "\n" +
		"        name: " + customHookName + "\n" +
		"        entry: bash -c 'go build ./...'\n" +
		"        language: system\n" +
		"        pass_filenames: false\n" +
		"        description: Run go build on all packages\n"

	// Also create a .commitlintrc.yaml file for additional configuration
	commitlintPath := filepath.Join(projectDir, ".commitlintrc.yaml")
	commitlintContent := "# .commitlintrc.yaml\n" +
		"extends:\n" +
		"  - conventional\n" +
		"rules:\n" +
		"  header-max-length: [2, always, 100]\n" +
		"  body-max-line-length: [2, always, 100]\n" +
		"  type-enum:\n" +
		"    - 2\n" +
		"    - always\n" +
		"    - - feat     # A new feature\n" +
		"      - fix      # A bug fix\n" +
		"      - docs     # Documentation only changes\n" +
		"      - style    # Changes that do not affect the meaning of the code\n" +
		"      - refactor # A code change that neither fixes a bug nor adds a feature\n" +
		"      - perf     # A code change that improves performance\n" +
		"      - test     # Adding missing tests or correcting existing tests\n" +
		"      - build    # Changes that affect the build system or external dependencies\n" +
		"      - ci       # Changes to CI configuration files and scripts\n" +
		"      - chore    # Other changes that don't modify src or test files\n" +
		"      - revert   # Reverts a previous commit\n"

	if err := os.WriteFile(preCommitConfigPath, []byte(preCommitConfigContent), 0600); err != nil {
		return err
	}

	return os.WriteFile(commitlintPath, []byte(commitlintContent), 0600)
}

// generateInitialCode creates the initial Go code for the project
func generateInitialCode(cfg *config.ProjectConfig, projectDir string) error {
	// Create main.go
	mainGoPath := filepath.Join(projectDir, "main.go")
	var mainGoContent string

	if cfg.UseCobra {
		// If using Cobra, create a CLI-style main.go
		mainGoContent = fmt.Sprintf("package main\n\n"+
			"import (\n"+
			"\t\"fmt\"\n"+
			"\t\"os\"\n\n"+
			"\t\"%s/cmd\"\n"+
			")\n\n"+
			"func main() {\n"+
			"\tif err := cmd.Execute(); err != nil {\n"+
			"\t\tfmt.Fprintln(os.Stderr, err)\n"+
			"\t\tos.Exit(1)\n"+
			"\t}\n"+
			"}\n", cfg.Module)

		// Create cmd directory structure for Cobra
		cmdDir := filepath.Join(projectDir, "cmd")
		if err := os.MkdirAll(cmdDir, 0755); err != nil {
			return err
		}

		// Create root.go for Cobra
		rootGoPath := filepath.Join(cmdDir, "root.go")

		// Start with imports
		rootGoContent := "package cmd\n\n" +
			"import (\n" +
			"\t\"fmt\"\n" +
			"\t\"os\"\n\n" +
			"\t\"github.com/spf13/cobra\"\n"

		// Add viper import if needed
		if cfg.UseViper {
			rootGoContent += "\t\"github.com/spf13/viper\"\n"
		}
		rootGoContent += ")\n\n"

		// Add config variable if viper is used
		if cfg.UseViper {
			rootGoContent += "var cfgFile string\n\n"
		}

		// Add root command
		rootGoContent += "// rootCmd represents the base command when called without any subcommands\n" +
			"var rootCmd = &cobra.Command{\n" +
			fmt.Sprintf("\tUse:   \"%s\",\n", strings.ToLower(cfg.Name)) +
			fmt.Sprintf("\tShort: \"%s\",\n", cfg.Description) +
			fmt.Sprintf("\tLong: `%s`,\n", cfg.Description) +
			"\t// Run: func(cmd *cobra.Command, args []string) { },\n" +
			"}\n\n" +
			"// Execute adds all child commands to the root command and sets flags appropriately.\n" +
			"// This is called by main.main(). It only needs to happen once to the rootCmd.\n" +
			"func Execute() error {\n" +
			"\treturn rootCmd.Execute()\n" +
			"}\n\n" +
			"func init() {\n"

		// Add viper initialization if needed
		if cfg.UseViper {
			rootGoContent += "\tcobra.OnInitialize(initConfig)\n\n" +
				"\t// Global flags\n" +
				fmt.Sprintf("\trootCmd.PersistentFlags().StringVar(&cfgFile, \"config\", \"\", \"config file (default is $HOME/.%s.yaml)\")\n",
					strings.ToLower(cfg.Name))
		} else {
			rootGoContent += "\t// Add your flags here\n"
		}

		rootGoContent += "}\n"

		// Add viper config function if needed
		if cfg.UseViper {
			rootGoContent += "\n// initConfig reads in config file and ENV variables if set.\n" +
				"func initConfig() {\n" +
				"\tif cfgFile != \"\" {\n" +
				"\t\t// Use config file from the flag.\n" +
				"\t\tviper.SetConfigFile(cfgFile)\n" +
				"\t} else {\n" +
				"\t\t// Find home directory.\n" +
				"\t\thome, err := os.UserHomeDir()\n" +
				"\t\tcobra.CheckErr(err)\n\n" +
				"\t\t// Search config in home directory with name \"." + strings.ToLower(cfg.Name) + "\" (without extension).\n" +
				"\t\tviper.AddConfigPath(home)\n" +
				"\t\tviper.SetConfigType(\"yaml\")\n" +
				"\t\tviper.SetConfigName(\"." + strings.ToLower(cfg.Name) + "\")\n" +
				"\t}\n\n" +
				"\tviper.AutomaticEnv() // read in environment variables that match\n\n" +
				"\t// If a config file is found, read it in.\n" +
				"\tif err := viper.ReadInConfig(); err == nil {\n" +
				"\t\tfmt.Fprintln(os.Stderr, \"Using config file:\", viper.ConfigFileUsed())\n" +
				"\t}\n" +
				"}\n"
		}

		if err := os.WriteFile(rootGoPath, []byte(rootGoContent), 0600); err != nil {
			return err
		}

		// Create version.go for Cobra
		versionGoPath := filepath.Join(cmdDir, "version.go")
		versionGoContent := "package cmd\n\n" +
			"import (\n" +
			"\t\"fmt\"\n\n" +
			"\t\"github.com/spf13/cobra\"\n" +
			")\n\n" +
			"// Version information - will be set during build via ldflags\n" +
			"var (\n" +
			"\tVersion   = \"dev\"\n" +
			"\tCommit    = \"none\"\n" +
			"\tBuildDate = \"unknown\"\n" +
			")\n\n" +
			"// versionCmd represents the version command\n" +
			"var versionCmd = &cobra.Command{\n" +
			"\tUse:   \"version\",\n" +
			"\tShort: \"Print the version information\",\n" +
			fmt.Sprintf("\tLong:  \"Display the version, commit, and build date information for the %s CLI\",\n", cfg.Name) +
			"\tRun: func(cmd *cobra.Command, args []string) {\n" +
			fmt.Sprintf("\t\tfmt.Println(\"%s CLI\")\n", cfg.Name) +
			"\t\tfmt.Println(\"--------\")\n" +
			"\t\tfmt.Printf(\"Version:    %%s\\n\", Version)\n" +
			"\t\tfmt.Printf(\"Commit:     %%s\\n\", Commit)\n" +
			"\t\tfmt.Printf(\"Build Date: %%s\\n\", BuildDate)\n" +
			"\t},\n" +
			"}\n\n" +
			"func init() {\n" +
			"\trootCmd.AddCommand(versionCmd)\n" +
			"}\n"

		if err := os.WriteFile(versionGoPath, []byte(versionGoContent), 0600); err != nil {
			return err
		}
	} else {
		// Simple main.go
		mainGoContent = fmt.Sprintf("package main\n\n"+
			"import (\n"+
			"\t\"fmt\"\n"+
			")\n\n"+
			"func main() {\n"+
			"\tfmt.Println(\"Hello from %s!\")\n"+
			"}\n", cfg.Name)
	}

	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0600); err != nil {
		return err
	}

	// Create main_test.go
	mainTestGoPath := filepath.Join(projectDir, "main_test.go")
	mainTestGoContent := "package main\n\n" +
		"import (\n" +
		"\t\"testing\"\n" +
		")\n\n" +
		"func TestMain(t *testing.T) {\n" +
		"\t// Add your tests here\n" +
		"}\n"

	return os.WriteFile(mainTestGoPath, []byte(mainTestGoContent), 0600)
}

// generateTemplates creates code templates for the project
func generateTemplates(cfg *config.ProjectConfig, projectDir string) error {
	// Create templates directory
	templatesDir := filepath.Join(projectDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return err
	}

	// Create example template file
	exampleTemplatePath := filepath.Join(templatesDir, "example.tmpl")
	exampleTemplateContent := "# Example Template for " + cfg.Name + "\n\n" +
		"This is an example template file that can be used with your application.\n\n" +
		"## Usage\n\n" +
		"```go\n" +
		"// Example code to use this template\n" +
		"tmpl, err := template.ParseFiles(\"templates/example.tmpl\")\n" +
		"if err != nil {\n" +
		"    log.Fatal(err)\n" +
		"}\n" +
		"data := map[string]interface{}{\n" +
		"    \"Title\": \"My Title\",\n" +
		"    \"Content\": \"My Content\",\n" +
		"}\n" +
		"err = tmpl.Execute(os.Stdout, data)\n" +
		"```\n"

	if err := os.WriteFile(exampleTemplatePath, []byte(exampleTemplateContent), 0600); err != nil {
		return err
	}

	return nil
}
