package types

// Server represents a MCP server definition
type Server struct {
	Name          string            `json:"name"`
	URL           string            `json:"url,omitempty"`
	Description   string            `json:"description"`
	Command       string            `json:"command"`
	Args          []string          `json:"args,omitempty"`
	Env           map[string]EnvVar `json:"env,omitempty"`
	Prerequisites []Prerequisite    `json:"prerequisites,omitempty"`
	Recommends    []string          `json:"recommends,omitempty"`
	Alternatives  []string          `json:"alternatives,omitempty"`
	Requires      []string          `json:"requires,omitempty"`
	HealthCheck   *HealthCheck      `json:"healthCheck,omitempty"`
	Scripts       map[string]string `json:"scripts,omitempty"`
}

// Preset represents a collection of servers
type Preset struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Servers     []string `json:"servers"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Value       string `json:"value"`
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
}

// Prerequisite represents a prerequisite for a server
type Prerequisite struct {
	Binary   string `json:"binary"`
	Check    string `json:"check"`
	Install  string `json:"install"`
	Required bool   `json:"required"`
}

// HealthCheck represents a server health check
type HealthCheck struct {
	Type          string   `json:"type"`
	Timeout       string   `json:"timeout"`
	ExpectTools   []string `json:"expectTools,omitempty"`
	ExpectMethods []string `json:"expectMethods,omitempty"`
}
