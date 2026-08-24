package config

import (
	"path"

	"github.com/sarielhp/mcpx/internal/envvault"
	"github.com/sarielhp/mcpx/internal/types"
)

// AntigravityWriter writes MCP server configs in Antigravity's format.
// Antigravity uses an "mcpServers" key with command/args/env format.
type AntigravityWriter struct {
	Dir       string
	Overwrite bool
}

// NewAntigravityWriter creates a writer targeting the given directory.
func NewAntigravityWriter(dir string, overwrite bool) *AntigravityWriter {
	return &AntigravityWriter{Dir: dir, Overwrite: overwrite}
}

// Write adds the given servers to the Antigravity config file.
func (w *AntigravityWriter) Write(servers []*types.Server, vault *envvault.Vault) error {
	file := path.Join(w.Dir, "antigravity.json")
	root := LoadConfig(file)
	mcp, ok := root["mcpServers"]
	if !ok || !IsObject(mcp) {
		mcp = map[string]any{}
		root["mcpServers"] = mcp
	}
	mcpObj := mcp.(map[string]any)
	for _, s := range servers {
		entry := map[string]any{
			"command": s.Command,
			"args":    s.Args,
			"env":     resolveEnv(s, vault),
		}
		mcpObj[s.Name] = entry
	}
	return WriteJSON(file, root)
}

func resolveEnv(s *types.Server, vault *envvault.Vault) map[string]string {
	env := map[string]string{}
	for name, spec := range s.Env {
		env[name] = vault.ResolveTemplateEnv(spec.Value)
	}
	return env
}
