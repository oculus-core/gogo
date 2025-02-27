package config

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

// NewDefaultProjectConfig creates a new project config with sensible defaults
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
