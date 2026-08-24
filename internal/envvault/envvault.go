package envvault

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// Vault manages environment variables across multiple tiers:
// system env -> ~/.config/mcpx/.env -> {targetDir}/.env
type Vault struct {
	ConfigDir string
	TargetDir string
}

// NewVault creates a Vault with default config and target directories.
func NewVault(targetDir string) *Vault {
	home := os.Getenv("HOME")
	if home == "" {
		home = "."
	}
	return &Vault{
		ConfigDir: path.Join(home, ".config", "mcpx"),
		TargetDir: targetDir,
	}
}

// LoadEnv merges environment variables from the config and target .env files.
// TargetDir/.env has the highest precedence.
func (v *Vault) LoadEnv() map[string]string {
	merged := map[string]string{}
	for k, val := range parseEnvFile(path.Join(v.ConfigDir, ".env")) {
		merged[k] = val
	}
	for k, val := range parseEnvFile(path.Join(v.TargetDir, ".env")) {
		merged[k] = val
	}
	return merged
}

// SetEnv sets an env var in the config-level .env file.
func (v *Vault) SetEnv(name, value string) error {
	if name == "" || strings.Contains(name, "=") || strings.Contains(name, "\n") {
		return fmt.Errorf("invalid env var name %q", name)
	}
	if err := os.MkdirAll(v.ConfigDir, 0o755); err != nil {
		return err
	}
	file := path.Join(v.ConfigDir, ".env")
	lines := readEnvLines(file)
	updated := false
	for i, line := range lines {
		key := line
		if idx := strings.Index(line, "="); idx >= 0 {
			key = line[:idx]
		}
		if key == name {
			lines[i] = name + "=" + value
			updated = true
			break
		}
	}
	if !updated {
		lines = append(lines, name+"="+value)
	}
	return writeLines(file, lines)
}

// GetEnv returns the resolved value of an env var, or nil if unset.
func (v *Vault) GetEnv(name string) (string, bool) {
	env := v.LoadEnv()
	if val, ok := env[name]; ok {
		return val, true
	}
	return "", false
}

// RemoveEnv removes an env var from the config-level .env file.
func (v *Vault) RemoveEnv(name string) error {
	file := path.Join(v.ConfigDir, ".env")
	lines := readEnvLines(file)
	var kept []string
	for _, line := range lines {
		key := line
		if idx := strings.Index(line, "="); idx >= 0 {
			key = line[:idx]
		}
		if key != name {
			kept = append(kept, line)
		}
	}
	return writeLines(file, kept)
}

// ListEnv returns all env vars from the merged environment.
func (v *Vault) ListEnv() map[string]string {
	return v.LoadEnv()
}

// MaskValue masks a secret value for display, keeping the last 4 chars.
func MaskValue(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return "****" + value[len(value)-4:]
}

// ResolveTemplateEnv resolves ${VAR} references in a template env value
// against the merged environment.
func (v *Vault) ResolveTemplateEnv(value string) string {
	env := v.LoadEnv()
	var out strings.Builder
	i := 0
	for i < len(value) {
		if value[i] == '$' && i+1 < len(value) && value[i+1] == '{' {
			end := strings.Index(value[i+2:], "}")
			if end >= 0 {
				name := value[i+2 : i+2+end]
				if val, ok := env[name]; ok {
					out.WriteString(val)
				}
				i += 2 + end + 1
				continue
			}
		}
		out.WriteByte(value[i])
		i++
	}
	return out.String()
}

func parseEnvFile(file string) map[string]string {
	result := map[string]string{}
	content, err := os.ReadFile(file)
	if err != nil {
		return result
	}
	text := string(content)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "="); idx >= 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			result[key] = val
		}
	}
	return result
}

func readEnvLines(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		return []string{}
	}
	return strings.Split(string(content), "\n")
}

func writeLines(file string, lines []string) error {
	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	return os.WriteFile(file, []byte(content), 0o644)
}
