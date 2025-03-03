package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oculus-core/gogo/pkg/config"
)

// TestProjectTypeSettings tests the type-specific settings logic
func TestProjectTypeSettings(t *testing.T) {
	tests := []struct {
		name           string
		startType      config.ProjectType
		targetType     config.ProjectType
		expectUseCobra bool
		expectUseViper bool
		expectUseGin   bool
		expectUseCmd   bool
	}{
		{
			name:           "Change to CLI",
			startType:      config.TypeDefault,
			targetType:     config.TypeCLI,
			expectUseCobra: true,
			expectUseViper: true,
			expectUseGin:   false,
			expectUseCmd:   true,
		},
		{
			name:           "Change to API",
			startType:      config.TypeDefault,
			targetType:     config.TypeAPI,
			expectUseCobra: false,
			expectUseViper: false,
			expectUseGin:   true,
			expectUseCmd:   true,
		},
		{
			name:           "Change to Library",
			startType:      config.TypeDefault,
			targetType:     config.TypeLibrary,
			expectUseCobra: false,
			expectUseViper: false,
			expectUseGin:   false,
			expectUseCmd:   false,
		},
		{
			name:           "Already CLI",
			startType:      config.TypeCLI,
			targetType:     config.TypeCLI,
			expectUseCobra: false,
			expectUseViper: false,
			expectUseGin:   false,
			expectUseCmd:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create test config
			cfg := config.NewDefaultProjectConfig()
			cfg.Type = tc.startType

			// Simulate the type change logic from the wizard
			prevType := cfg.Type
			cfg.Type = tc.targetType

			// Apply type-specific settings when type changes
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

			// Verify settings were applied correctly
			assert.Equal(t, tc.targetType, cfg.Type)
			assert.Equal(t, tc.expectUseCobra, cfg.UseCobra)
			assert.Equal(t, tc.expectUseViper, cfg.UseViper)
			assert.Equal(t, tc.expectUseGin, cfg.UseGin)
			assert.Equal(t, tc.expectUseCmd, cfg.UseCmd)
		})
	}
}

// TestGetStructureDefaults tests the structure defaults based on project type
func TestGetStructureDefaults(t *testing.T) {
	tests := []struct {
		name           string
		appType        config.ProjectType
		expectCmd      bool
		expectPkg      bool
		expectInternal bool
	}{
		{
			name:           "CLI Structure",
			appType:        config.TypeCLI,
			expectCmd:      true,
			expectPkg:      true,
			expectInternal: true,
		},
		{
			name:           "API Structure",
			appType:        config.TypeAPI,
			expectCmd:      true,
			expectPkg:      true,
			expectInternal: true,
		},
		{
			name:           "Library Structure",
			appType:        config.TypeLibrary,
			expectCmd:      false,
			expectPkg:      true,
			expectInternal: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create config with the given project type
			cfg := config.GetProjectConfigForType(tc.appType)

			// Get structure defaults
			defaults := getStructureDefaults(cfg)

			// Check if cmd is included
			cmdIncluded := contains(defaults, "cmd (application entrypoints)")
			assert.Equal(t, tc.expectCmd, cmdIncluded)

			// Check if pkg is included
			pkgIncluded := contains(defaults, "pkg (public packages)")
			assert.Equal(t, tc.expectPkg, pkgIncluded)

			// Check if internal is included
			internalIncluded := contains(defaults, "internal (private packages)")
			assert.Equal(t, tc.expectInternal, internalIncluded)
		})
	}
}

// TestGetDepsDefaults tests the dependencies defaults based on project type
func TestGetDepsDefaults(t *testing.T) {
	tests := []struct {
		name        string
		appType     config.ProjectType
		expectCobra bool
		expectViper bool
		expectGin   bool
	}{
		{
			name:        "CLI Dependencies",
			appType:     config.TypeCLI,
			expectCobra: true,
			expectViper: true,
			expectGin:   false,
		},
		{
			name:        "API Dependencies",
			appType:     config.TypeAPI,
			expectCobra: false,
			expectViper: false,
			expectGin:   true,
		},
		{
			name:        "Library Dependencies",
			appType:     config.TypeLibrary,
			expectCobra: false,
			expectViper: false,
			expectGin:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create config with the given project type
			cfg := config.GetProjectConfigForType(tc.appType)

			// Get dependencies defaults
			defaults := getDepsDefaults(cfg)

			// Check if cobra is included
			cobraIncluded := contains(defaults, "Cobra (CLI framework)")
			assert.Equal(t, tc.expectCobra, cobraIncluded)

			// Check if viper is included
			viperIncluded := contains(defaults, "Viper (configuration)")
			assert.Equal(t, tc.expectViper, viperIncluded)

			// Check if gin is included in the module
			ginIncluded := cfg.UseGin
			assert.Equal(t, tc.expectGin, ginIncluded)
		})
	}
}
