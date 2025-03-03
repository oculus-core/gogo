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

		// Create a .gitkeep file to ensure the directory is tracked by Git
		gitkeepPath := filepath.Join(dirPath, ".gitkeep")
		if err := os.WriteFile(gitkeepPath, []byte(""), 0600); err != nil {
			return fmt.Errorf("failed to create .gitkeep in %s: %v", dir, err)
		}
	}

	// Generate initial code based on application type
	if err := generateInitialCodeByType(cfg, projectDir); err != nil {
		return err
	}

	// Generate config file
	if err := generateConfigFile(cfg, projectDir); err != nil {
		return err
	}

	// Generate go.mod file
	if err := generateGoMod(cfg, projectDir); err != nil {
		return err
	}

	// Generate GitHub Actions workflows if enabled
	if cfg.UseGitHubActions {
		if err := generateGitHubWorkflows(cfg, projectDir); err != nil {
			return err
		}
	}

	// Generate linter configuration if enabled
	if cfg.UseLinters {
		if err := generateLinterConfig(cfg, projectDir); err != nil {
			return err
		}
	}

	// Generate pre-commit hooks configuration if enabled
	if cfg.UsePreCommitHooks {
		if err := generatePreCommitConfig(cfg, projectDir); err != nil {
			return err
		}
	}

	return nil
}

// generateInitialCodeByType generates initial code based on the application type
func generateInitialCodeByType(cfg *config.ProjectConfig, projectDir string) error {
	switch cfg.Type {
	case config.TypeCLI:
		return generateCLICode(cfg, projectDir)
	case config.TypeAPI:
		return generateAPICode(cfg, projectDir)
	case config.TypeLibrary:
		return generateLibraryCode(cfg, projectDir)
	default:
		return generateDefaultCode(cfg, projectDir)
	}
}

// generateCLICode generates code for a CLI application
func generateCLICode(cfg *config.ProjectConfig, projectDir string) error {
	// Create cmd directory structure
	cmdDir := filepath.Join(projectDir, "cmd", cfg.Name)
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		return fmt.Errorf("failed to create cmd directory: %v", err)
	}

	// Generate main.go
	mainPath := filepath.Join(cmdDir, "main.go")
	mainContent := fmt.Sprintf(`package main

import (
	"fmt"
	"os"

	"%s/cmd/%s/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
`, cfg.Module, cfg.Name)

	if err := os.WriteFile(mainPath, []byte(mainContent), 0600); err != nil {
		return fmt.Errorf("failed to create main.go: %v", err)
	}

	// Create cmd package directory
	cmdPkgDir := filepath.Join(cmdDir, "cmd")
	if err := os.MkdirAll(cmdPkgDir, 0755); err != nil {
		return fmt.Errorf("failed to create cmd package directory: %v", err)
	}

	// Generate root.go with Cobra
	rootPath := filepath.Join(cmdPkgDir, "root.go")
	rootContent := fmt.Sprintf(`package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "%s",
	Short: "A brief description of your application",
	Long: `+"`"+`A longer description that spans multiple lines and likely contains
examples and usage of using your application.`+"`"+`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.%s.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".%s" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".%s")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
`, cfg.Name, cfg.Name, cfg.Name, cfg.Name)

	if err := os.WriteFile(rootPath, []byte(rootContent), 0600); err != nil {
		return fmt.Errorf("failed to create root.go: %v", err)
	}

	// Generate version.go
	versionPath := filepath.Join(cmdPkgDir, "version.go")
	versionContent := fmt.Sprintf(`package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `+"`"+`Print the version, commit, and build date information for your application.`+"`"+`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %%s (%%s) built on %%s\n", Version, Commit, BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
`, cfg.Name)

	if err := os.WriteFile(versionPath, []byte(versionContent), 0600); err != nil {
		return fmt.Errorf("failed to create version.go: %v", err)
	}

	return nil
}

// generateAPICode generates code for an API application
func generateAPICode(cfg *config.ProjectConfig, projectDir string) error {
	// Create cmd directory structure
	cmdDir := filepath.Join(projectDir, "cmd", cfg.Name)
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		return fmt.Errorf("failed to create cmd directory: %v", err)
	}

	// Generate main.go
	mainPath := filepath.Join(cmdDir, "main.go")
	mainContent := fmt.Sprintf(`package main

import (
	"log"

	"%s/internal/api"
	"%s/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %%v", err)
	}

	server := api.NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %%v", err)
	}
}
`, cfg.Module, cfg.Module)

	if err := os.WriteFile(mainPath, []byte(mainContent), 0600); err != nil {
		return fmt.Errorf("failed to create main.go: %v", err)
	}

	// Create internal/config directory
	configDir := filepath.Join(projectDir, "internal", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create internal/config directory: %v", err)
	}

	// Generate config.go
	configPath := filepath.Join(configDir, "config.go")
	configContent := `package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server ServerConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port int
	Host string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT: %v", err)
		}
	}

	host := "localhost"
	if hostEnv := os.Getenv("HOST"); hostEnv != "" {
		host = hostEnv
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
			Host: host,
		},
	}, nil
}
`

	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		return fmt.Errorf("failed to create config.go: %v", err)
	}

	// Create internal/api directory
	apiDir := filepath.Join(projectDir, "internal", "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create internal/api directory: %v", err)
	}

	// Generate server.go
	serverPath := filepath.Join(apiDir, "server.go")
	serverContent := fmt.Sprintf(`package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"%s/internal/config"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	cfg    *config.Config
}

// NewServer creates a new API server
func NewServer(cfg *config.Config) *Server {
	router := gin.Default()

	server := &Server{
		router: router,
		cfg:    cfg,
	}

	server.registerRoutes()

	return server
}

// Run starts the server
func (s *Server) Run() error {
	addr := fmt.Sprintf("%%s:%%d", s.cfg.Server.Host, s.cfg.Server.Port)
	return s.router.Run(addr)
}

// registerRoutes sets up the API routes
func (s *Server) registerRoutes() {
	s.router.GET("/health", s.healthCheck)

	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/hello", s.helloWorld)
	}
}

// healthCheck handles the health check endpoint
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// helloWorld handles the hello world endpoint
func (s *Server) helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}
`, cfg.Module)

	if err := os.WriteFile(serverPath, []byte(serverContent), 0600); err != nil {
		return fmt.Errorf("failed to create server.go: %v", err)
	}

	return nil
}

// generateLibraryCode generates code for a library
func generateLibraryCode(cfg *config.ProjectConfig, projectDir string) error {
	// Create pkg directory structure
	pkgDir := filepath.Join(projectDir, "pkg", cfg.Name)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("failed to create pkg directory: %v", err)
	}

	// Generate library.go
	libPath := filepath.Join(pkgDir, fmt.Sprintf("%s.go", cfg.Name))
	libContent := fmt.Sprintf(`package %s

// Version is the current version of the library
const Version = "0.1.1"

// Hello returns a greeting message
func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return "Hello, " + name + "!"
}
`, cfg.Name)

	if err := os.WriteFile(libPath, []byte(libContent), 0600); err != nil {
		return fmt.Errorf("failed to create library file: %v", err)
	}

	// Generate test file
	testPath := filepath.Join(pkgDir, fmt.Sprintf("%s_test.go", cfg.Name))
	testContent := fmt.Sprintf(`package %s

import "testing"

func TestHello(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty name",
			input:    "",
			expected: "Hello, World!",
		},
		{
			name:     "with name",
			input:    "Gopher",
			expected: "Hello, Gopher!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hello(tt.input); got != tt.expected {
				t.Errorf("Hello() = %%q, want %%q", got, tt.expected)
			}
		})
	}
}
`, cfg.Name)

	if err := os.WriteFile(testPath, []byte(testContent), 0600); err != nil {
		return fmt.Errorf("failed to create test file: %v", err)
	}

	return nil
}

// generateDefaultCode generates code for a default project
func generateDefaultCode(cfg *config.ProjectConfig, projectDir string) error {
	// Create a simple main.go in the project root
	mainPath := filepath.Join(projectDir, "main.go")
	mainContent := fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Hello from %s!")
}
`, cfg.Name)

	if err := os.WriteFile(mainPath, []byte(mainContent), 0600); err != nil {
		return fmt.Errorf("failed to create main.go: %v", err)
	}

	return nil
}

// generateConfigFile creates the gogo.yaml configuration file
func generateConfigFile(cfg *config.ProjectConfig, projectDir string) error {
	configPath := filepath.Join(projectDir, "gogo.yaml")

	configContent := fmt.Sprintf(`# Gogo Project Configuration
# Generated on: %s

# Project Information
project:
  name: %q
  module: %q
  description: %q
  license: %q
  author: %q

# Project Structure
structure:
  use_cmd: %t
  use_internal: %t
  use_pkg: %t
  use_test: %t
  use_docs: %t

# Generated Files
files:
  create_readme: %t
  create_license: %t
  create_makefile: %t

# Code Quality
quality:
  use_linters: %t
  use_pre_commit_hooks: %t
  use_git_hooks: %t

# Dependencies
dependencies:
  use_cobra: %t
  use_viper: %t

# CI/CD
cicd:
  use_github_actions: %t
`,
		time.Now().Format(time.RFC3339),
		cfg.Name,
		cfg.Module,
		cfg.Description,
		cfg.License,
		cfg.Author,
		cfg.UseCmd,
		cfg.UseInternal,
		cfg.UsePkg,
		cfg.UseTest,
		cfg.UseDocs,
		cfg.CreateReadme,
		cfg.CreateLicense,
		cfg.CreateMakefile,
		cfg.UseLinters,
		cfg.UsePreCommitHooks,
		cfg.UseGitHooks,
		cfg.UseCobra,
		cfg.UseViper,
		cfg.UseGitHubActions,
	)

	return os.WriteFile(configPath, []byte(configContent), 0600)
}

// generateRootFiles creates the basic files at the project root
func generateRootFiles(cfg *config.ProjectConfig, projectDir string) error {
	// Generate config file
	if err := generateConfigFile(cfg, projectDir); err != nil {
		return fmt.Errorf("failed to generate config file: %w", err)
	}

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

	// Create the workflow directory if it doesn't exist
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflow directory: %v", err)
	}

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

// TODO: Add template generation in a future version
// generateTemplates creates code templates for the project
//
//nolint:unused
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
