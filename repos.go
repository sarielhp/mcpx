package main

import "fmt"

// Repo describes a source repository maintained by this project.
type Repo struct {
	Name        string
	Description string
	URL         string
}

// projectRepos lists the two repositories that make up the mcpx ecosystem.
var projectRepos = []Repo{
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
}

// cmdRepos prints the project repositories with HTML links.
func cmdRepos() error {
	for _, r := range projectRepos {
		fmt.Printf("%s — <a href=%q>%s</a>\n  %s\n", r.Name, r.URL, r.URL, r.Description)
	}
	return nil
}
