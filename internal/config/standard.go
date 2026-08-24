package config

import (
	"path"

	"github.com/sarielhp/mcpx/internal/envvault"
	"github.com/sarielhp/mcpx/internal/types"
)

// StandardWriter writes MCP server configs in the standard mcpServers format
// used by Cursor, Claude Desktop, and VS Code.
type StandardWriter struct {
	Dir       string
	Overwrite bool
}

// NewStandardWriter creates a writer targeting the given directory.
func NewStandardWriter(dir string, overwrite bool) *StandardWriter {
	return &StandardWriter{Dir: dir, Overwrite: overwrite}
}

// Write adds the given servers to the standard config file.
func (w *StandardWriter) Write(servers []*types.Server, vault *envvault.Vault) error {
	file := path.Join(w.Dir, "mcp.json")
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
