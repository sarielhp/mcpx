package config

import (
	"path"

	"github.com/sarielhp/mcpx/internal/envvault"
	"github.com/sarielhp/mcpx/internal/types"
)

// OpenCodeWriter writes MCP server configs in OpenCode's format.
// OpenCode uses an "mcp" key with type: "local" / command: [...] / enabled: true.
type OpenCodeWriter struct {
	Dir       string
	Overwrite bool
}

// NewOpenCodeWriter creates a writer targeting the given directory.
func NewOpenCodeWriter(dir string, overwrite bool) *OpenCodeWriter {
	return &OpenCodeWriter{Dir: dir, Overwrite: overwrite}
}

// Write adds the given servers to the OpenCode config file.
func (w *OpenCodeWriter) Write(servers []*types.Server, vault *envvault.Vault) error {
	file := path.Join(w.Dir, "opencode.json")
	root := LoadConfig(file)
	mcp, ok := root["mcp"]
	if !ok || !IsObject(mcp) {
		mcp = map[string]any{}
		root["mcp"] = mcp
	}
	mcpObj := mcp.(map[string]any)
	for _, s := range servers {
		entry := map[string]any{
			"type":    "local",
			"command": buildCommand(s, vault),
			"enabled": true,
		}
		mcpObj[s.Name] = entry
	}
	return WriteJSON(file, root)
}

func buildCommand(s *types.Server, vault *envvault.Vault) []any {
	cmd := []any{s.Command}
	for _, arg := range s.Args {
		cmd = append(cmd, arg)
	}
	return cmd
}
