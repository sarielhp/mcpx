package main

import (
	"testing"
)

func TestResolveItems(t *testing.T) {
	// Test that ResolveItems can be called without error
	// This is a basic test - more comprehensive tests would require
	// mocking the registry
	items := []string{"gopls"}
	servers, unknown := ResolveItems(items)

	if len(servers) == 0 {
		t.Error("Expected at least one server")
	}
	if len(unknown) != 0 {
		t.Errorf("Expected no unknown items, got: %v", unknown)
	}
}

func TestApplyShorthands(t *testing.T) {
	// Test that applyShorthands can be called without error
	// This is a basic test - more comprehensive tests would require
	// mocking the command line parsing
	args := []string{"ls"}
	result := applyShorthands(args)

	// Should expand to template list
	if len(result) < 2 || result[0] != "template" || result[1] != "list" {
		t.Errorf("Expected shorthand expansion to template list, got: %v", result)
	}
}
