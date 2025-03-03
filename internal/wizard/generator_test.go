package wizard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oculus-core/gogo/pkg/config"
)

func TestGenerateInitialCodeByType(t *testing.T) {
	tests := []struct {
		name        string
		appType     config.ProjectType
		checkFiles  []string
		checkAbsent []string
	}{
		{
			name:    "CLI Project",
			appType: config.TypeCLI,
			checkFiles: []string{
				"cmd/{{.Name}}/main.go",
				"cmd/{{.Name}}/cmd/root.go",
				"cmd/{{.Name}}/cmd/version.go",
			},
			checkAbsent: []string{
				"internal/api/server.go",
			},
		},
		{
			name:    "API Project",
			appType: config.TypeAPI,
			checkFiles: []string{
				"cmd/{{.Name}}/main.go",
				"internal/config/config.go",
				"internal/api/server.go",
			},
			checkAbsent: []string{
				"cmd/{{.Name}}/cmd/root.go",
			},
		},
		{
			name:    "Library Project",
			appType: config.TypeLibrary,
			checkFiles: []string{
				"pkg/{{.Name}}/{{.Name}}.go",
				"pkg/{{.Name}}/{{.Name}}_test.go",
			},
			checkAbsent: []string{
				"cmd/{{.Name}}/main.go",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir := t.TempDir()

			// Create project config
			cfg := config.GetProjectConfigForType(tc.appType)
			cfg.Name = "testproj"
			cfg.Module = "github.com/example/testproj"

			// Generate project
			projectDir := filepath.Join(tmpDir, cfg.Name)
			if err := os.MkdirAll(projectDir, 0755); err != nil {
				t.Fatalf("failed to create project directory: %v", err)
			}

			// Create necessary parent directories
			for _, dir := range []string{"cmd", "internal", "pkg"} {
				if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
			}

			// Generate initial code by type
			err := generateInitialCodeByType(cfg, projectDir)
			assert.NoError(t, err)

			// Check expected files exist
			for _, file := range tc.checkFiles {
				filePath := strings.Replace(file, "{{.Name}}", cfg.Name, -1)
				fullPath := filepath.Join(projectDir, filePath)
				_, err := os.Stat(fullPath)
				assert.NoError(t, err, "File should exist: "+filePath)

				// Verify file has content
				content, err := os.ReadFile(fullPath)
				assert.NoError(t, err)
				assert.NotEmpty(t, content, "File should not be empty: "+filePath)
			}

			// Check absent files
			for _, file := range tc.checkAbsent {
				filePath := strings.Replace(file, "{{.Name}}", cfg.Name, -1)
				fullPath := filepath.Join(projectDir, filePath)
				_, err := os.Stat(fullPath)
				assert.True(t, os.IsNotExist(err), "File should not exist: "+filePath)
			}
		})
	}
}

func TestGenerateConfigFile(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()

	// Create project config
	cfg := config.NewCLIProjectConfig()
	cfg.Name = "configtest"
	cfg.Module = "github.com/example/configtest"

	// Generate project
	projectDir := filepath.Join(tmpDir, cfg.Name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project directory: %v", err)
	}

	// Generate config file
	err := generateConfigFile(cfg, projectDir)
	assert.NoError(t, err)

	// Check if config file exists
	configPath := filepath.Join(projectDir, "gogo.yaml")
	_, err = os.Stat(configPath)
	assert.NoError(t, err, "Config file should exist")

	// Check content of config file
	content, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), `project:
  name: "configtest"`)
	assert.Contains(t, string(content), `dependencies:
  use_cobra: true
  use_viper: true`)
}

func TestGenerateProject(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()

	// Test different project types
	projTypes := []config.ProjectType{
		config.TypeCLI,
		config.TypeAPI,
		config.TypeLibrary,
	}

	for _, projType := range projTypes {
		t.Run(string(projType), func(t *testing.T) {
			// Create project config
			cfg := config.GetProjectConfigForType(projType)
			cfg.Name = "testproject-" + string(projType)
			cfg.Module = "github.com/example/" + cfg.Name

			// Generate project
			err := GenerateProject(cfg, tmpDir)
			assert.NoError(t, err)

			// Check project directory was created
			projectDir := filepath.Join(tmpDir, cfg.Name)
			_, err = os.Stat(projectDir)
			assert.NoError(t, err, "Project directory should exist")

			// Check gogo.yaml config file
			configPath := filepath.Join(projectDir, "gogo.yaml")
			_, err = os.Stat(configPath)
			assert.NoError(t, err, "Config file should exist")

			// Check go.mod file
			goModPath := filepath.Join(projectDir, "go.mod")
			_, err = os.Stat(goModPath)
			assert.NoError(t, err, "go.mod file should exist")

			// Check for type-specific files
			switch projType {
			case config.TypeCLI:
				rootPath := filepath.Join(projectDir, "cmd", cfg.Name, "cmd", "root.go")
				_, err = os.Stat(rootPath)
				assert.NoError(t, err, "root.go should exist")

			case config.TypeAPI:
				serverPath := filepath.Join(projectDir, "internal", "api", "server.go")
				_, err = os.Stat(serverPath)
				assert.NoError(t, err, "server.go should exist")

			case config.TypeLibrary:
				libPath := filepath.Join(projectDir, "pkg", cfg.Name, cfg.Name+".go")
				_, err = os.Stat(libPath)
				assert.NoError(t, err, "library file should exist")
			}
		})
	}
}
