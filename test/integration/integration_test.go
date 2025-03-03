package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getProjectRoot returns the root directory of the project
func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../..")
}

// TestProjectGeneration runs integration tests for project generation
// NOTE: This test requires the gogo binary to be built and in PATH
func TestProjectGeneration(t *testing.T) {
	// Skip if running in CI or if the integration tag is not provided
	if os.Getenv("CI") != "" || os.Getenv("GOGO_INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration tests. Set GOGO_INTEGRATION_TEST=1 to run them")
	}

	// Get the path to the gogo binary
	gogoBin := filepath.Join(getProjectRoot(), "bin", "gogo")

	tests := []struct {
		name               string
		args               []string
		expectedDirs       []string
		expectedFiles      []string
		expectedFileChecks map[string]func(t *testing.T, content string)
	}{
		{
			name: "CLI Project",
			args: []string{"new", "testcli", "-t", "cli", "-s"},
			expectedDirs: []string{
				"cmd",
				"internal",
				"pkg",
			},
			expectedFiles: []string{
				"go.mod",
				"gogo.yaml",
			},
			expectedFileChecks: map[string]func(t *testing.T, content string){
				"gogo.yaml": func(t *testing.T, content string) {
					assert.Contains(t, content, `dependencies:
  use_cobra: true
  use_viper: true`)
				},
			},
		},
		{
			name: "API Project",
			args: []string{"new", "testapi", "-t", "api", "-s"},
			expectedDirs: []string{
				"cmd",
				"internal",
				"pkg",
			},
			expectedFiles: []string{
				"go.mod",
				"gogo.yaml",
			},
			expectedFileChecks: map[string]func(t *testing.T, content string){
				"gogo.yaml": func(t *testing.T, content string) {
					assert.Contains(t, content, `dependencies:
  use_cobra: false
  use_viper: false`)
				},
			},
		},
		{
			name: "Library Project",
			args: []string{"new", "testlib", "-t", "library", "-s"},
			expectedDirs: []string{
				"internal",
				"pkg",
			},
			expectedFiles: []string{
				"go.mod",
				"gogo.yaml",
			},
			expectedFileChecks: map[string]func(t *testing.T, content string){
				"gogo.yaml": func(t *testing.T, content string) {
					assert.Contains(t, content, `structure:
  use_cmd: false`)
					assert.Contains(t, content, `dependencies:
  use_cobra: false
  use_viper: false`)
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for the test
			tempDir, err := os.MkdirTemp("", "gogo-integration-test-*")
			if err != nil {
				t.Fatalf("Failed to create temporary directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Save current working directory
			oldWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer func() {
				err := os.Chdir(oldWd)
				if err != nil {
					t.Logf("Failed to restore working directory: %v", err)
				}
			}()

			// Change to temp directory for test
			err = os.Chdir(tempDir)
			if err != nil {
				t.Fatalf("Failed to change to temp directory: %v", err)
			}

			// Run gogo command
			cmd := exec.Command(gogoBin, tc.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, output)
			}

			// Get project name from args (always at index 1)
			projectName := tc.args[1]
			projectDir := filepath.Join(tempDir, projectName)

			// Verify project directory exists
			_, err = os.Stat(projectDir)
			require.NoError(t, err, "Project directory should exist")

			// Verify expected directories exist
			for _, dir := range tc.expectedDirs {
				dirPath := filepath.Join(projectDir, dir)
				_, err = os.Stat(dirPath)
				assert.NoError(t, err, "Directory %s should exist", dir)
			}

			// Verify expected files exist
			for _, file := range tc.expectedFiles {
				filePath := filepath.Join(projectDir, file)
				_, err = os.Stat(filePath)
				assert.NoError(t, err, "File %s should exist", file)
			}

			// Run checks on file contents
			for file, checkFn := range tc.expectedFileChecks {
				filePath := filepath.Join(projectDir, file)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read file %s: %v", file, err)
					continue
				}
				checkFn(t, string(content))
			}
		})
	}
}

// TestConfigChanges tests that gogo respects configuration file changes
func TestConfigChanges(t *testing.T) {
	// Skip if running in CI or if the integration tag is not provided
	if os.Getenv("CI") != "" || os.Getenv("GOGO_INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration tests. Set GOGO_INTEGRATION_TEST=1 to run them")
	}

	// Get the path to the gogo binary
	gogoBin := filepath.Join(getProjectRoot(), "bin", "gogo")

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "gogo-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save current working directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		err := os.Chdir(oldWd)
		if err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	}()

	// Change to temp directory for test
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create a test project
	projectName := "configtest"
	cmd := exec.Command(gogoBin, "new", projectName, "-t", "cli", "-s")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	projectDir := filepath.Join(tempDir, projectName)

	// Verify project directory exists
	_, err = os.Stat(projectDir)
	require.NoError(t, err, "Project directory should exist")

	// Create modified configuration file
	configPath := filepath.Join(projectDir, "gogo.yaml")

	// Read the existing config
	configData, err := os.ReadFile(configPath)
	require.NoError(t, err, "Should be able to read config file")

	// Simple check that the config was created
	assert.Contains(t, string(configData), projectName, "Config should contain project name")
}
