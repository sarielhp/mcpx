package runner

import (
	"os/exec"
	"sync"
	"time"

	"github.com/sarielhp/mcpx/internal/types"
)

// Result represents the outcome of a server health check.
type Result struct {
	Name    string
	OK      bool
	Message string
	Elapsed time.Duration
}

// Runner probes MCP servers in parallel to verify their handshakes.
type Runner struct {
	Timeout time.Duration
}

// NewRunner creates a Runner with the given per-server timeout.
func NewRunner(timeout time.Duration) *Runner {
	return &Runner{Timeout: timeout}
}

// TestAll probes all servers concurrently and returns results.
func (r *Runner) TestAll(servers []*types.Server) []Result {
	results := make([]Result, len(servers))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i, s := range servers {
		wg.Add(1)
		go func(idx int, srv *types.Server) {
			defer wg.Done()
			res := r.testOne(srv)
			mu.Lock()
			results[idx] = res
			mu.Unlock()
		}(i, s)
	}
	wg.Wait()
	return results
}

func (r *Runner) testOne(s *types.Server) Result {
	start := time.Now()
	ok, msg := handshake(s, r.Timeout)
	return Result{
		Name:    s.Name,
		OK:      ok,
		Message: msg,
		Elapsed: time.Since(start),
	}
}

// handshake spawns the server command and performs a JSON-RPC initialize
// handshake, returning whether the server responded.
func handshake(s *types.Server, timeout time.Duration) (bool, string) {
	cmd := exec.Command(s.Command, s.Args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return false, "stdin pipe failed: " + err.Error()
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, "stdout pipe failed: " + err.Error()
	}
	if err := cmd.Start(); err != nil {
		return false, "failed to start: " + err.Error()
	}

	req := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"mcpx","version":"1.0"}}}`
	if _, err := stdin.Write([]byte(req + "\n")); err != nil {
		cmd.Process.Kill()
		return false, "write failed: " + err.Error()
	}
	stdin.Close()

	deadline := time.Now().Add(timeout)
	var buf []byte
	for time.Now().Before(deadline) {
		chunk := make([]byte, 4096)
		n, err := stdout.Read(chunk)
		if err != nil {
			break
		}
		if n == 0 {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		buf = append(buf, chunk[:n]...)
		if containsNewline(buf) {
			break
		}
	}
	cmd.Process.Kill()
	if len(buf) == 0 {
		return false, "no response within timeout"
	}
	if !contains(buf, `"result"`) {
		return false, "unexpected response"
	}
	return true, "ok"
}

func containsNewline(b []byte) bool {
	for _, c := range b {
		if c == '\n' {
			return true
		}
	}
	return false
}

func contains(b []byte, needle string) bool {
	if len(needle) == 0 {
		return true
	}
	for i := 0; i+len(needle) <= len(b); i++ {
		match := true
		for j := 0; j < len(needle); j++ {
			if b[i+j] != byte(needle[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
