package main

import "fmt"

// Repo describes a source repository maintained by this project.
type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// cmdRepos prints the configured project repositories with HTML links.
func cmdRepos() error {
	cfg, err := ensureConfig()
	if err != nil {
		return err
	}
	for _, r := range cfg.Repos {
		fmt.Printf("%s — <a href=%q>%s</a>\n  %s\n", r.Name, r.URL, r.URL, r.Description)
	}
	return nil
}
