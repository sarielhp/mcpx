package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"

	"github.com/sarielhp/mcpx/internal/types"
)

// TemplateRepoBase is the base URL for the official template catalog on GitHub.
const TemplateRepoBase = "https://raw.githubusercontent.com/sarielhp/mcp-templates/main"

// TemplateIndex is the shape of INDEX.json served by the template repo.
type TemplateIndex struct {
	Version string `json:"version"`
	Updated string `json:"updated"`
	Servers map[string]struct {
		Description string `json:"description"`
		URL         string `json:"url"`
	} `json:"servers"`
	Presets map[string]struct {
		Description string   `json:"description"`
		Servers     []string `json:"servers"`
	} `json:"presets"`
}

// RemoteRegistry provides template discovery from the remote GitHub catalog,
// cached locally under ~/.config/mcpx/templates/. It falls back to the
// embedded set so the binary stays functional offline.
type RemoteRegistry struct{}

// registry is the shared template registry used across all commands.
var registry = &RemoteRegistry{}

// cacheDir returns the local template cache directory.
func cacheDir() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = "."
	}
	return path.Join(home, ".config", "mcpx", "templates")
}

// fetchRemote downloads a path from the remote template repo.
func fetchRemote(rel string) ([]byte, error) {
	url := TemplateRepoBase + "/" + rel
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: HTTP %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// cachedRead returns the cached bytes at rel, or nil if absent.
func cachedRead(rel string) ([]byte, error) {
	return os.ReadFile(path.Join(cacheDir(), rel))
}

// cachedWrite stores bytes at rel in the local cache (creating parent dirs).
func cachedWrite(rel string, content []byte) error {
	dest := path.Join(cacheDir(), rel)
	if err := os.MkdirAll(path.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, content, 0o644)
}

// fetchDef returns the raw JSON for a server/preset definition, trying cache
// then remote.
func fetchDef(rel string) ([]byte, error) {
	if content, err := cachedRead(rel); err == nil {
		return content, nil
	}
	content, err := fetchRemote(rel)
	if err != nil {
		return nil, err
	}
	_ = cachedWrite(rel, content)
	return content, nil
}

func serverRel(name string) string { return "servers/" + name + "/template.json" }
func presetRel(name string) string { return "presets/" + name + "/preset.json" }

// fetchIndex returns the (cached-or-remote) template catalog.
func fetchIndex() (*TemplateIndex, error) {
	content, err := fetchDef("INDEX.json")
	if err != nil {
		return nil, err
	}
	var idx TemplateIndex
	if err := json.Unmarshal(content, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}

// ListServers returns remote server names merged with the embedded set.
func (r *RemoteRegistry) ListServers() []string {
	var names []string
	if idx, err := fetchIndex(); err == nil && idx.Servers != nil {
		for n := range idx.Servers {
			names = append(names, n)
		}
	}
	names = mergeStrings(names, ListServers())
	sort.Strings(names)
	return names
}

// ListPresets returns remote preset names merged with the embedded set.
func (r *RemoteRegistry) ListPresets() []string {
	var names []string
	if idx, err := fetchIndex(); err == nil && idx.Presets != nil {
		for n := range idx.Presets {
			names = append(names, n)
		}
	}
	names = mergeStrings(names, ListPresets())
	sort.Strings(names)
	return names
}

// GetServer returns a server definition, trying remote then embedded.
func (r *RemoteRegistry) GetServer(name string) (*types.Server, error) {
	content, err := fetchDef(serverRel(name))
	if err == nil {
		var s types.Server
		if err := json.Unmarshal(content, &s); err == nil {
			return &s, nil
		}
	}
	return GetServer(name)
}

// GetPreset returns a preset definition, trying remote then embedded.
func (r *RemoteRegistry) GetPreset(name string) (*types.Preset, error) {
	content, err := fetchDef(presetRel(name))
	if err == nil {
		var p types.Preset
		if err := json.Unmarshal(content, &p); err == nil {
			return &p, nil
		}
	}
	return GetPreset(name)
}

// Update refreshes the local template cache from the remote catalog.
// Returns the number of files synced (0 is valid when nothing new).
func (r *RemoteRegistry) Update() (int, error) {
	idx, err := fetchIndex()
	if err != nil {
		return 0, err
	}
	updated := 0
	for name := range idx.Servers {
		if err := r.syncDef(serverRel(name)); err == nil {
			updated++
		}
	}
	for name := range idx.Presets {
		if err := r.syncDef(presetRel(name)); err == nil {
			updated++
		}
	}
	return updated, nil
}

// syncDef downloads a definition into the local cache unconditionally.
func (r *RemoteRegistry) syncDef(rel string) error {
	content, err := fetchRemote(rel)
	if err != nil {
		return err
	}
	return cachedWrite(rel, content)
}

func mergeStrings(a, b []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, s := range append(a, b...) {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}
