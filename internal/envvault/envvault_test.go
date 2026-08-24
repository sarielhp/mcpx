package envvault

import (
	"os"
	"testing"
)

func TestResolveTemplateEnv(t *testing.T) {
	// Create a test vault
	vault := &Vault{
		ConfigDir: "/tmp",
		TargetDir: "/tmp",
	}

	// Set up some environment variables for testing
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("ANOTHER", "another_value")
	defer os.Unsetenv("TEST_VAR")
	defer os.Unsetenv("ANOTHER")

	// Test basic substitution
	result := vault.ResolveTemplateEnv("${TEST_VAR}")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test no substitution
	result = vault.ResolveTemplateEnv("no_substitution")
	if result != "no_substitution" {
		t.Errorf("Expected 'no_substitution', got '%s'", result)
	}

	// Test mixed content
	result = vault.ResolveTemplateEnv("prefix_${TEST_VAR}_suffix")
	if result != "prefix_test_value_suffix" {
		t.Errorf("Expected 'prefix_test_value_suffix', got '%s'", result)
	}

	// Test undefined variable
	result = vault.ResolveTemplateEnv("${UNDEFINED}")
	if result != "${UNDEFINED}" {
		t.Errorf("Expected '${UNDEFINED}', got '%s'", result)
	}
}

func TestMaskValue(t *testing.T) {
	// Test masking of short value
	result := MaskValue("1234")
	if result != "****4" {
		t.Errorf("Expected '****4', got '%s'", result)
	}

	// Test masking of long value
	result = MaskValue("1234567890")
	if result != "****7890" {
		t.Errorf("Expected '****7890', got '%s'", result)
	}
}

func TestParseEnvFile(t *testing.T) {
	// Create a temporary .env file
	content := `TEST_VAR=test_value
# This is a comment
ANOTHER_VAR=another_value

# Another comment
EMPTY_VAR=
`

	tmpFile, err := os.CreateTemp("", "test.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Parse the file
	result := parseEnvFile(tmpFile.Name())

	// Check results
	if result["TEST_VAR"] != "test_value" {
		t.Errorf("Expected TEST_VAR=test_value, got %s", result["TEST_VAR"])
	}
	if result["ANOTHER_VAR"] != "another_value" {
		t.Errorf("Expected ANOTHER_VAR=another_value, got %s", result["ANOTHER_VAR"])
	}
	if result["EMPTY_VAR"] != "" {
		t.Errorf("Expected EMPTY_VAR=, got %s", result["EMPTY_VAR"])
	}
	if _, exists := result["# This is a comment"]; exists {
		t.Error("Comment should not be parsed as variable")
	}
}
