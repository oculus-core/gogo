# Gogo CLI Technical Specification

## Project Overview

The "gogo" CLI is a tool designed to scaffold Go projects with best practices and project structure,
providing an interactive wizard for project configuration. This specification outlines the technical
details for implementing the CLI.

## Architecture

### High-Level Architecture

```mermaid
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│    Gogo CLI     │────▶│     Wizard      │────▶│    Generator    │
│                 │     │                 │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                                               │
        │                                               │
        ▼                                               ▼
┌─────────────────┐                           ┌─────────────────┐
│                 │                           │                 │
│  Configuration  │                           │   Templates     │
│                 │                           │                 │
└─────────────────┘                           └─────────────────┘
```

### Components

1. **CLI Layer**: Built with Cobra, handles command-line arguments, flags, and user interaction.
2. **Configuration Layer**: Managed by Viper, handles loading and saving configuration settings.
3. **Wizard Layer**: Interactive wizard for project configuration.
4. **Generator Layer**: Business logic for project generation.
5. **Templates Layer**: Templates for generating project files.

## Directory Structure

```bash
gogo/
├── cmd/
│   └── gogo/          # CLI commands
│       ├── root.go    # Root command
│       ├── new.go     # New project command
│       └── version.go # Version command
├── pkg/
│   ├── config/        # Configuration management
│   │   └── config.go  # Configuration structure
│   ├── prompt/        # Prompt utilities
│   └── template/      # Template utilities
├── internal/
│   ├── wizard/        # Wizard implementation
│   │   ├── wizard.go  # Interactive wizard
│   │   └── generator.go # Project generation
│   └── generator/     # Code generation
├── bin/               # Build artifacts
├── main.go            # Entry point
├── go.mod             # Go modules
├── go.sum             # Go modules checksums
├── .golangci.yml      # Linter configuration
├── .pre-commit-config.yaml # Pre-commit hooks
├── .goreleaser.yml    # GoReleaser configuration
├── Makefile           # Build automation
├── .gitignore         # Git ignore
├── README.md          # Documentation
└── LICENSE            # License
```

## Data Models

### Project Configuration Model

```go
// ProjectConfig represents the configuration for a gogo project
type ProjectConfig struct {
    // General project information
    Name        string `yaml:"name" json:"name"`
    Module      string `yaml:"module" json:"module"`
    Description string `yaml:"description" json:"description"`
    License     string `yaml:"license" json:"license"`
    Author      string `yaml:"author" json:"author"`

    // Project structure options
    UseCmd         bool `yaml:"use_cmd" json:"use_cmd"`
    UseInternal    bool `yaml:"use_internal" json:"use_internal"`
    UsePkg         bool `yaml:"use_pkg" json:"use_pkg"`
    UseTest        bool `yaml:"use_test" json:"use_test"`
    UseDocs        bool `yaml:"use_docs" json:"use_docs"`
    CreateReadme   bool `yaml:"create_readme" json:"create_readme"`
    CreateLicense  bool `yaml:"create_license" json:"create_license"`
    CreateMakefile bool `yaml:"create_makefile" json:"create_makefile"`

    // Code quality tools
    UseLinters        bool `yaml:"use_linters" json:"use_linters"`
    UsePreCommitHooks bool `yaml:"use_pre_commit_hooks" json:"use_pre_commit_hooks"`
    UseGitHooks       bool `yaml:"use_git_hooks" json:"use_git_hooks"`

    // Dependencies
    UseCobra bool `yaml:"use_cobra" json:"use_cobra"`
    UseViper bool `yaml:"use_viper" json:"use_viper"`

    // CI/CD
    UseGitHubActions bool `yaml:"use_github_actions" json:"use_github_actions"`
}
```

## API Design

### Wizard Service

```go
// Wizard service for interactive project configuration
type Wizard interface {
    // RunWizard runs the interactive wizard for project configuration
    RunWizard(cfg *config.ProjectConfig) error

    // GenerateProject generates a new project based on the configuration
    GenerateProject(cfg *config.ProjectConfig, outputDir string) error
}
```

### Project Generator

```go
// Generator service for project generation
type Generator interface {
    // GenerateProject generates a new project based on the configuration
    GenerateProject(cfg *config.ProjectConfig, outputDir string) error

    // GenerateRootFiles generates the root files for the project
    GenerateRootFiles(cfg *config.ProjectConfig, projectDir string) error

    // GenerateGoMod generates the go.mod file for the project
    GenerateGoMod(cfg *config.ProjectConfig, projectDir string) error

    // GenerateGitHubWorkflows generates GitHub Actions workflows
    GenerateGitHubWorkflows(cfg *config.ProjectConfig, projectDir string) error

    // GenerateLinterConfig generates linter configuration
    GenerateLinterConfig(cfg *config.ProjectConfig, projectDir string) error

    // GeneratePreCommitConfig generates pre-commit hooks configuration
    GeneratePreCommitConfig(cfg *config.ProjectConfig, projectDir string) error

    // GenerateInitialCode generates initial code for the project
    GenerateInitialCode(cfg *config.ProjectConfig, projectDir string) error

    // GenerateTemplates generates template files for the project
    GenerateTemplates(cfg *config.ProjectConfig, projectDir string) error
}
```

## Wizard Flow

1. User runs `gogo new [project-name]` command
2. CLI initializes default project configuration
3. If not skipped, the interactive wizard is launched
4. User configures project information (name, module, description, license, author)
5. User selects project structure (cmd, internal, pkg, test, docs)
6. User selects files to generate (README, LICENSE, Makefile)
7. User configures code quality tools (linters, pre-commit hooks, git hooks)
8. User configures dependencies (Cobra, Viper)
9. User configures CI/CD (GitHub Actions)
10. User reviews and confirms configuration
11. Project is generated based on the configuration
12. User is presented with next steps

## Error Handling Strategy

1. Proper error propagation through the call stack
2. Contextual error messages for user-friendly feedback
3. Graceful degradation when errors occur
4. Detailed error logging for debugging

## Testing Strategy

### Unit Tests

- Test individual components in isolation
- Mock dependencies using interfaces
- Focus on business logic in the generator layer
- Aim for high code coverage

### Integration Tests

- Test interaction between components
- Test file system operations in a temporary directory
- Test interactive wizard with simulated input

### Acceptance Tests

- Test end-to-end scenarios
- Test using CLI commands as a user would

## Project Generation

The project generation process includes:

1. **Directory Structure Creation**

   - Create project directory
   - Create subdirectories based on configuration (cmd, internal, pkg, test, docs)

2. **File Generation**

   - Generate README.md with project information
   - Generate LICENSE based on selected license
   - Generate .gitignore with common Go patterns
   - Generate Makefile with common targets
   - Generate go.mod file with project module path
   - Generate GitHub Actions workflows
   - Generate linter configuration
   - Generate pre-commit hooks configuration
   - Generate initial code with project structure

3. **Template Generation**
   - Generate example templates for the project

## Dependencies

1. **Required Dependencies**

   - Go 1.16+
   - Cobra v1.9.1+ (CLI framework)
   - Viper v1.19.0+ (Configuration management)
   - Survey v2.3.7+ (Interactive prompts)
   - LipGloss v1.0.0+ (Terminal styling)

2. **Development Dependencies**
   - GoReleaser (Release automation)
   - GolangCI-Lint (Linting)
   - Pre-commit (Git hooks)

## Configuration Defaults

```go
// Default project configuration
func NewDefaultProjectConfig() *ProjectConfig {
    return &ProjectConfig{
        Name:              "my-project",
        Module:            "github.com/username/my-project",
        Description:       "A Go project",
        License:           "MIT",
        Author:            "",
        UseCmd:            true,
        UseInternal:       true,
        UsePkg:            true,
        UseTest:           true,
        UseDocs:           true,
        CreateReadme:      true,
        CreateLicense:     true,
        CreateMakefile:    true,
        UseLinters:        true,
        UsePreCommitHooks: true,
        UseGitHooks:       true,
        UseCobra:          false,
        UseViper:          false,
        UseGitHubActions:  true,
    }
}
```

## Command Structure

```bash
gogo
├── new [project-name]  # Create a new Go project
│   ├── --output, -o    # Output directory for the project
│   └── --skip-wizard, -s # Skip the interactive wizard
└── version             # Show version information
```

## Implementation Guidelines

1. Use standard Go project layout
2. Follow Go code style and conventions
3. Implement proper error handling
4. Document all public APIs
5. Use semantic versioning for releases
6. Implement comprehensive testing
7. Provide user-friendly error messages
8. Ensure cross-platform compatibility
9. Support configuration via command-line flags and environment variables

## Deployment and Distribution

1. **Packaging**

   - Build binaries for multiple platforms (Linux, macOS, Windows)
   - Create distribution packages (Homebrew)
   - Sign binaries and packages

2. **Updates**

   - Implement self-update capability
   - Check for updates periodically
   - Notify users of new versions

3. **Installation**
   - Go install: `go install github.com/oculus-core/gogo@latest`
   - Homebrew: `brew install oculus-core/homebrew-gogo/gogo`
   - Download binary from GitHub releases
