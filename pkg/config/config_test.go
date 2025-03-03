package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	// Setup test directory and files
	testdataDir := "testdata"
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}
	defer os.RemoveAll(testdataDir)

	// Create valid config file
	validCfgPath := filepath.Join(testdataDir, "valid_config.yaml")
	validCfgContent := `
name: test-project
module: github.com/example/test-project
description: A test project
license: MIT
author: Test Author
type: cli
use_cmd: true
use_cobra: true
use_viper: true
`
	if err := os.WriteFile(validCfgPath, []byte(validCfgContent), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Create malformed config file
	malformedCfgPath := filepath.Join(testdataDir, "malformed_config.yaml")
	malformedCfgContent := `
name: test-project
module: github.com/example/test-project
description: "incomplete string
license: MIT
author: Test Author
type: cli
use_cmd: true
use_cobra: true
use_viper: [broken, array, syntax
`
	if err := os.WriteFile(malformedCfgPath, []byte(malformedCfgContent), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Test loading valid configuration
	cfg, err := LoadConfigFromFile(validCfgPath)
	assert.NoError(t, err)
	assert.Equal(t, "test-project", cfg.Name)
	assert.Equal(t, TypeCLI, cfg.Type)
	assert.True(t, cfg.UseCobra)
	assert.True(t, cfg.UseViper)

	// Test loading missing file
	_, err = LoadConfigFromFile(filepath.Join(testdataDir, "nonexistent.yaml"))
	assert.Error(t, err)

	// Test loading malformed YAML
	_, err = LoadConfigFromFile(malformedCfgPath)
	assert.Error(t, err)
}

func TestSaveConfigToFile(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()

	// Create config to save
	cfg := NewCLIProjectConfig()
	cfg.Name = "test-save"

	tempFile := filepath.Join(tmpDir, "test_config.yaml")

	// Test saving
	err := SaveConfigToFile(cfg, tempFile)
	assert.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(tempFile)
	assert.NoError(t, err)

	// Test loading saved file
	loadedCfg, err := LoadConfigFromFile(tempFile)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Name, loadedCfg.Name)
	assert.Equal(t, cfg.Type, loadedCfg.Type)
}

func TestProjectTypeConfig(t *testing.T) {
	// Test CLI config
	cliCfg := NewCLIProjectConfig()
	assert.Equal(t, TypeCLI, cliCfg.Type)
	assert.True(t, cliCfg.UseCobra)
	assert.True(t, cliCfg.UseViper)

	// Test API config
	apiCfg := NewAPIProjectConfig()
	assert.Equal(t, TypeAPI, apiCfg.Type)
	assert.True(t, apiCfg.UseGin)

	// Test Library config
	libCfg := NewLibraryProjectConfig()
	assert.Equal(t, TypeLibrary, libCfg.Type)
	assert.False(t, libCfg.UseCmd)

	// Test GetProjectConfigForType
	defaultCfg := GetProjectConfigForType(TypeDefault)
	assert.Equal(t, TypeDefault, defaultCfg.Type)

	// Test with unknown type (should return default)
	unknownType := ProjectType("unknown")
	unknownCfg := GetProjectConfigForType(unknownType)
	assert.Equal(t, TypeDefault, unknownCfg.Type)
}
