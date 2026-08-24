package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/sarielhp/mcpx/internal/types"
)

// Embedded templates are stored under templates/ in the source tree.
//
//go:embed templates/servers/*/template.json
var serverTemplates embed.FS

//go:embed templates/presets/*/preset.json
var presetTemplates embed.FS

// ListServers returns the names of all embedded server templates.
func ListServers() []string {
	return listNames(serverTemplates, "templates/servers", "template.json")
}

// ListPresets returns the names of all embedded preset templates.
func ListPresets() []string {
	return listNames(presetTemplates, "templates/presets", "preset.json")
}

// GetServer loads a server template by name.
func GetServer(name string) (*types.Server, error) {
	content, err := readTemplate(serverTemplates, "templates/servers", name, "template.json")
	if err != nil {
		return nil, err
	}
	var server types.Server
	if err := json.Unmarshal(content, &server); err != nil {
		return nil, fmt.Errorf("invalid template for server %q: %v", name, err)
	}
	return &server, nil
}

// GetPreset loads a preset template by name.
func GetPreset(name string) (*types.Preset, error) {
	content, err := readTemplate(presetTemplates, "templates/presets", name, "preset.json")
	if err != nil {
		return nil, err
	}
	var preset types.Preset
	if err := json.Unmarshal(content, &preset); err != nil {
		return nil, fmt.Errorf("invalid template for preset %q: %v", name, err)
	}
	return &preset, nil
}

// ResolveItems expands a list of item names into concrete servers.
// Preset names are expanded into their member servers; unknown names are
// returned as-is so the caller can report them.
func ResolveItems(items []string) ([]*types.Server, []string) {
	var servers []*types.Server
	var unknown []string
	for _, item := range items {
		if server, err := GetServer(item); err == nil {
			servers = append(servers, server)
			continue
		}
		if preset, err := GetPreset(item); err == nil {
			for _, member := range preset.Servers {
				if server, err := GetServer(member); err == nil {
					servers = append(servers, server)
				} else {
					unknown = append(unknown, member)
				}
			}
			continue
		}
		unknown = append(unknown, item)
	}
	return servers, unknown
}

func listNames(fsys fs.FS, base, file string) []string {
	var names []string
	entries, _ := fs.ReadDir(fsys, base)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := fsys.Open(path.Join(base, entry.Name(), file)); err == nil {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names
}

func readTemplate(fsys fs.FS, base, name, file string) ([]byte, error) {
	if name == "" || strings.Contains(name, "/") || strings.Contains(name, "..") {
		return nil, fmt.Errorf("invalid template name %q", name)
	}
	return fs.ReadFile(fsys, path.Join(base, name, file))
}
