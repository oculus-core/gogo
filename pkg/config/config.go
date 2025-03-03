package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ProjectType represents the type of project to generate
type ProjectType string

const (
	// TypeCLI is for command-line interface projects
	TypeCLI ProjectType = "cli"
	// TypeAPI is for REST API projects
	TypeAPI ProjectType = "api"
	// TypeLibrary is for library/package projects
	TypeLibrary ProjectType = "library"
	// TypeDefault is the default project type
	TypeDefault ProjectType = "default"
)

// ProjectConfig represents the configuration for a gogo project
type ProjectConfig struct {
	// General project information
	Name        string      `yaml:"name" json:"name"`
	Module      string      `yaml:"module" json:"module"`
	Description string      `yaml:"description" json:"description"`
	License     string      `yaml:"license" json:"license"`
	Author      string      `yaml:"author" json:"author"`
	Type        ProjectType `yaml:"type" json:"type"`

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
	UseGin   bool `yaml:"use_gin" json:"use_gin"`

	// CI/CD
	UseGitHubActions bool `yaml:"use_github_actions" json:"use_github_actions"`
}

// NewDefaultProjectConfig creates a new project config with sensible defaults
func NewDefaultProjectConfig() *ProjectConfig {
	return &ProjectConfig{
		Name:              "my-project",
		Module:            "github.com/username/my-project",
		Description:       "A Go project",
		License:           "MIT",
		Author:            "",
		Type:              TypeDefault,
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
		UseGin:            false,
		UseGitHubActions:  true,
	}
}

// NewCLIProjectConfig creates a new project config for CLI applications
func NewCLIProjectConfig() *ProjectConfig {
	cfg := NewDefaultProjectConfig()
	cfg.Type = TypeCLI
	cfg.UseCobra = true
	cfg.UseViper = true
	return cfg
}

// NewAPIProjectConfig creates a new project config for API applications
func NewAPIProjectConfig() *ProjectConfig {
	cfg := NewDefaultProjectConfig()
	cfg.Type = TypeAPI
	cfg.UseGin = true
	return cfg
}

// NewLibraryProjectConfig creates a new project config for library projects
func NewLibraryProjectConfig() *ProjectConfig {
	cfg := NewDefaultProjectConfig()
	cfg.Type = TypeLibrary
	cfg.UseCmd = false
	return cfg
}

// GetProjectConfigForType returns a project config for the specified project type
func GetProjectConfigForType(projType ProjectType) *ProjectConfig {
	switch projType {
	case TypeCLI:
		return NewCLIProjectConfig()
	case TypeAPI:
		return NewAPIProjectConfig()
	case TypeLibrary:
		return NewLibraryProjectConfig()
	default:
		return NewDefaultProjectConfig()
	}
}

// LoadConfigFromFile loads a project configuration from a YAML file
func LoadConfigFromFile(filePath string) (*ProjectConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &cfg, nil
}

// SaveConfigToFile saves a project configuration to a YAML file
func SaveConfigToFile(cfg *ProjectConfig, filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
