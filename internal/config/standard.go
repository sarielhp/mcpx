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
	Subpath   string // Optional subdirectory path (e.g., ".cursor", ".claude", ".vscode")
}

// NewStandardWriter creates a writer targeting the given directory.
func NewStandardWriter(dir string, overwrite bool) *StandardWriter {
	return &StandardWriter{Dir: dir, Overwrite: overwrite}
}

// NewStandardWriterWithSubpath creates a writer targeting a subdirectory.
func NewStandardWriterWithSubpath(dir string, overwrite bool, subpath string) *StandardWriter {
	return &StandardWriter{Dir: dir, Overwrite: overwrite, Subpath: subpath}
}

// Write adds the given servers to the standard config file.
func (w *StandardWriter) Write(servers []*types.Server, vault *envvault.Vault) error {
	var file string
	if w.Subpath != "" {
		file = path.Join(w.Dir, w.Subpath, "mcp.json")
	} else {
		file = path.Join(w.Dir, "mcp.json")
	}
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
