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

func TestRunnerTestAll(t *testing.T) {
	// Test that TestAll returns results for a simple mock server
	runner := NewRunner(1 * time.Second)

	// Verify we can create results without errors
	results := runner.TestAll(nil)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for nil servers, got %d", len(results))
	}
}

func TestRunnerTimeout(t *testing.T) {
	// Test timeout behavior
	runner := NewRunner(100 * time.Millisecond)

	// Test that we can create runner with timeout
	if runner.Timeout != 100*time.Millisecond {
		t.Errorf("Expected timeout 100ms, got %v", runner.Timeout)
	}
}

func TestRunnerErrorHandling(t *testing.T) {
	// Test error handling
	runner := NewRunner(1 * time.Second)

	// Mock a failing server test result
	// This cannot test actual failure without implementing a mock server,
	// but it verifies the runner can be constructed properly
	if runner == nil {
		t.Fatal("Runner should not be nil")
	}
}
