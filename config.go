package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
)

// mcpxConfig is the user-facing config file (default ~/.config/mcpx/config.json).
// It stores the repository list shown by --repos and the template catalog base
// URL so users can point at their own forks.
type mcpxConfig struct {
	TemplateRepo string `json:"templateRepo"`
	Repos        []Repo `json:"repos"`
}

// defaultConfig returns the built-in repos and template repo, used when no
// config file exists yet.
func defaultConfig() *mcpxConfig {
	return &mcpxConfig{
		TemplateRepo: TemplateRepoBase,
		Repos: []Repo{
			{
				Name:        "mcpx",
				Description: "MCP server manager — this CLI",
				URL:         "https://github.com/sarielhp/mcpx",
			},
			{
				Name:        "mcp-templates",
				Description: "MCP server & preset template catalog",
				URL:         "https://github.com/sarielhp/mcp-templates",
			},
		},
	}
}

// configPath returns the location of the user config file.
func configPath() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = "."
	}
	return path.Join(home, ".config", "mcpx", "config.json")
}

// loadConfig reads the user config, falling back to defaults when the file is
// missing or partially empty.
func loadConfig() (*mcpxConfig, error) {
	def := defaultConfig()
	content, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return def, nil
		}
		return nil, err
	}
	var cfg mcpxConfig
	if err := json.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}
	if cfg.TemplateRepo == "" {
		cfg.TemplateRepo = def.TemplateRepo
	}
	if len(cfg.Repos) == 0 {
		cfg.Repos = def.Repos
	}
	return &cfg, nil
}

// saveConfig writes the user config file, creating parent dirs.
func saveConfig(c *mcpxConfig) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c); err != nil {
		return err
	}
	dest := configPath()
	if err := os.MkdirAll(path.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, buf.Bytes(), 0o644)
}

// ensureConfig loads the config and writes the defaults out if no file exists
// yet, so the user can edit a concrete file. First-run only.
func ensureConfig() (*mcpxConfig, error) {
	if _, err := os.Stat(configPath()); err == nil {
		return loadConfig()
	}
	cfg := defaultConfig()
	if err := saveConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// templateRepoBase returns the configured template catalog base URL, falling
// back to the built-in constant.
func templateRepoBase() string {
	cfg, err := loadConfig()
	if err != nil || cfg.TemplateRepo == "" {
		return TemplateRepoBase
	}
	return cfg.TemplateRepo
}
