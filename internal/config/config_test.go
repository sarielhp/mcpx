package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sarielhp/mcpx/internal/envvault"
	"github.com/sarielhp/mcpx/internal/types"
)

func TestOpenCodeWriter(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "opencode_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test server
	server := &types.Server{
		Name:        "test-server",
		Description: "Test server",
		Command:     "echo",
		Args:        []string{"hello"},
		Env:         map[string]types.EnvVar{},
	}

	// Create a vault
	vault := &envvault.Vault{
		ConfigDir: "/tmp",
		TargetDir: tmpDir,
	}

	// Create writer and write
	writer := NewOpenCodeWriter(tmpDir, false)
	err = writer.Write([]*types.Server{server}, vault)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the file was created
	content, err := os.ReadFile(filepath.Join(tmpDir, "opencode.json"))
	if err != nil {
		t.Fatal(err)
	}

	// Verify content contains expected structure
	contentStr := string(content)
	if !strings.Contains(contentStr, `"test-server"`) {
		t.Errorf("Expected server in output, got: %s", contentStr)
	}
	// The actual format is different - check for the correct structure
	if !strings.Contains(contentStr, `"command":`) {
		t.Errorf("Expected command in output, got: %s", contentStr)
	}
}

func TestAntigravityWriter(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "antigravity_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test server
	server := &types.Server{
		Name:        "test-server",
		Description: "Test server",
		Command:     "echo",
		Args:        []string{"hello"},
		Env: map[string]types.EnvVar{
			"TEST_VAR": {
				Value:       "${TEST_VAR}",
				Required:    true,
				Description: "Test variable",
			},
		},
	}

	// Create a vault
	vault := &envvault.Vault{
		ConfigDir: "/tmp",
		TargetDir: tmpDir,
	}

	// Create writer and write
	writer := NewAntigravityWriter(tmpDir, false)
	err = writer.Write([]*types.Server{server}, vault)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the file was created
	content, err := os.ReadFile(filepath.Join(tmpDir, "antigravity.json"))
	if err != nil {
		t.Fatal(err)
	}

	// Verify content contains expected structure
	contentStr := string(content)
	if !strings.Contains(contentStr, `"test-server"`) {
		t.Errorf("Expected server in output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, `"command":"echo"`) {
		t.Errorf("Expected command 'echo' in output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, `"args":["hello"]`) {
		t.Errorf("Expected args ['hello'] in output, got: %s", contentStr)
	}
}

func TestStandardWriter(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "standard_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test server
	server := &types.Server{
		Name:        "test-server",
		Description: "Test server",
		Command:     "echo",
		Args:        []string{"hello"},
		Env:         map[string]types.EnvVar{},
	}

	// Create a vault
	vault := &envvault.Vault{
		ConfigDir: "/tmp",
		TargetDir: tmpDir,
	}

	// Create writer and write
	writer := NewStandardWriter(tmpDir, false)
	err = writer.Write([]*types.Server{server}, vault)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the file was created
	content, err := os.ReadFile(filepath.Join(tmpDir, "mcp.json"))
	if err != nil {
		t.Fatal(err)
	}

	// Verify content contains expected structure
	contentStr := string(content)
	if !strings.Contains(contentStr, `"test-server"`) {
		t.Errorf("Expected server in output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, `"command":"echo"`) {
		t.Errorf("Expected command 'echo' in output, got: %s", contentStr)
	}
}
