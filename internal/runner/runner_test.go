package runner

import (
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	// Create a simple runner with a short timeout
	runner := NewRunner(1 * time.Second)

	// Test that it can be created without error
	if runner == nil {
		t.Fatal("Runner should not be nil")
	}

	// Test that timeout is set correctly
	if runner.Timeout != 1*time.Second {
		t.Errorf("Expected timeout 1s, got %v", runner.Timeout)
	}
}
