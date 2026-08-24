package types

import (
	"testing"
)

func TestServer(t *testing.T) {
	s := &Server{
		Name:        "test",
		Description: "Test server",
		Command:     "echo",
		Args:        []string{"hello"},
		Env: map[string]EnvVar{
			"TEST_VAR": {
				Value:       "${TEST_VAR}",
				Required:    true,
				Description: "Test variable",
			},
		},
		Prerequisites: []Prerequisite{
			{
				Binary:   "echo",
				Check:    "echo --version",
				Install:  "install echo",
				Required: true,
			},
		},
		Recommends:   []string{"other"},
		Alternatives: []string{},
		Requires:     []string{},
		HealthCheck: &HealthCheck{
			Type:          "stdio",
			Timeout:       "5s",
			ExpectTools:   []string{"test_tool"},
			ExpectMethods: []string{},
		},
	}

	if s.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", s.Name)
	}
	if s.Command != "echo" {
		t.Errorf("Expected command 'echo', got '%s'", s.Command)
	}
	if len(s.Args) != 1 || s.Args[0] != "hello" {
		t.Errorf("Expected args ['hello'], got %v", s.Args)
	}
}
