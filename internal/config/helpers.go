package config

import (
	"encoding/json"
	"os"
	"path"
	"sort"
)

// LoadConfig reads a JSON config file into a map, returning an empty map if
// the file does not exist or is invalid.
func LoadConfig(file string) map[string]any {
	content, err := os.ReadFile(file)
	if err != nil {
		return map[string]any{}
	}
	var data map[string]any
	if err := json.Unmarshal(content, &data); err != nil {
		return map[string]any{}
	}
	return data
}

// WriteJSON writes a map as pretty JSON to the given file, creating parent
// directories as needed.
func WriteJSON(file string, data map[string]any) error {
	if err := os.MkdirAll(path.Dir(file), 0o755); err != nil {
		return err
	}
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(file, content, 0o644)
}

// IsObject reports whether v is a JSON object (map).
func IsObject(v any) bool {
	_, ok := v.(map[string]any)
	return ok
}

// Keys returns the sorted keys of a map.
func Keys(m map[string]any) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
