package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/sarielhp/clihelp"

	"github.com/sarielhp/mcpx/internal/config"
	"github.com/sarielhp/mcpx/internal/envvault"
	"github.com/sarielhp/mcpx/internal/runner"
	"github.com/sarielhp/mcpx/internal/types"
)

func main() {
	app := buildApp()

	if os.Getenv("CLIHELP_GEN") != "" {
		_, gerr := clihelp.RenderMarkdown(app, clihelp.MarkdownOptions{Dir: "docs/clihelp"})
		if gerr != nil {
			fmt.Fprintf(os.Stderr, "generate docs: %v\n", gerr)
			os.Exit(1)
		}
		return
	}

	// Top-level convenience flags that bypass the command dispatcher.
	rawArgs := os.Args[1:]
	if len(rawArgs) == 1 && (rawArgs[0] == "--repos" || rawArgs[0] == "--repositories") {
		if err := cmdRepos(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	args := applyShorthands(rawArgs)
	if err := app.Execute(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// applyShorthands rewrites backward-compatible shorthand invocations into
// applyShorthands rewrites backward-compatible shorthand invocations into
// their canonical command forms before dispatch.
func applyShorthands(args []string) []string {
	if len(args) == 0 {
		return args
	}
	first := args[0]
	switch first {
	case "status":
		return []string{"list"}
	case "template":
		if len(args) > 1 {
			switch args[1] {
			case "list":
				return []string{"search"}
			case "show":
				return append([]string{"info"}, args[2:]...)
			case "update":
				return []string{"registry", "update"}
			}
		}
		return []string{"search"}
	case "help", "--help", "-h", "--version", "-v", "version":
		return args
	default:
		if !isKnownCommand(first) {
			return append([]string{"add"}, args...)
		}
		return args
	}
}

func buildApp() *clihelp.App {
	var dir string
	var dryRun bool
	var overwrite bool
	var opencode bool
	var antigravity bool
	var cursor bool
	var claude bool
	var vscode bool
	var all bool
	var withRecommended bool
	var verbose bool
	var timeout int

	app := &clihelp.App{
		Name:        "mcpx",
		Description: "MCPX - Manage MCP servers and presets",
		Version:     "0.2.2",
		Commands: []clihelp.Command{
			{
				Name:        "add",
				Description: "Add servers or presets to workspace (auto-initializes if new)",
				UsageLine:   "mcpx add <items...> [options]",
				Args:        clihelp.MinimumNArgs(1),
				Options: []clihelp.Option{
					clihelp.String(&dir, "--dir DIR", ".", "Target directory"),
					clihelp.Bool(&dryRun, "--dry-run", false, "Show what would be done without writing"),
					clihelp.Bool(&overwrite, "--overwrite", false, "Overwrite existing config entries"),
					clihelp.Bool(&opencode, "--opencode", false, "Write OpenCode config"),
					clihelp.Bool(&antigravity, "--antigravity", false, "Write Antigravity config"),
					clihelp.Bool(&cursor, "--cursor", false, "Write Cursor config"),
					clihelp.Bool(&claude, "--claude", false, "Write Claude config"),
					clihelp.Bool(&vscode, "--vscode", false, "Write VS Code config"),
					clihelp.Bool(&all, "--all", false, "Write all supported config formats"),
					clihelp.Bool(&withRecommended, "--with-recommended", false, "Auto-install recommended companions"),
				},
				Run: func(ctx *clihelp.Context) error {
					return cmdAdd(ctx.Args, dir, dryRun, overwrite, opencode, antigravity, cursor, claude, vscode, all, withRecommended)
				},
			},
			{
				Name:        "remove",
				Aliases:     []string{"rm"},
				Description: "Remove servers from workspace",
				UsageLine:   "mcpx remove <names...> [options]",
				Args:        clihelp.MinimumNArgs(1),
				Options:     []clihelp.Option{clihelp.String(&dir, "--dir DIR", ".", "Target directory")},
				Run:         func(ctx *clihelp.Context) error { return cmdRm(ctx.Args, dir) },
			},
			{
				Name:        "list",
				Aliases:     []string{"ls"},
				Description: "List active servers in workspace",
				UsageLine:   "mcpx list [options]",
				Args:        clihelp.NoArgs,
				Options:     []clihelp.Option{clihelp.String(&dir, "--dir DIR", ".", "Target directory")},
				Run:         func(ctx *clihelp.Context) error { return cmdList(dir) },
			},
			{
				Name:        "sync",
				Aliases:     []string{"update"},
				Description: "Sync workspace configs with latest credentials",
				UsageLine:   "mcpx sync [names...] [options]",
				Options:     []clihelp.Option{clihelp.String(&dir, "--dir DIR", ".", "Target directory")},
				Run:         func(ctx *clihelp.Context) error { return cmdSync(ctx.Args, dir) },
			},
			{
				Name:        "check",
				Aliases:     []string{"test", "validate"},
				Description: "Validate server handshakes (parallel)",
				UsageLine:   "mcpx check [names...] [options]",
				Args:        clihelp.MaximumNArgs(23),
				Options: []clihelp.Option{
					clihelp.String(&dir, "--dir DIR", ".", "Target directory"),
					clihelp.Int(&timeout, "--timeout SEC", 6, "Per-server timeout in seconds"),
					clihelp.Bool(&verbose, "--verbose", false, "Show detailed output"),
				},
				Run: func(ctx *clihelp.Context) error {
					return cmdTest(ctx.Args, timeout, verbose)
				},
			},
			{
				Name:        "search",
				Description: "List or search available servers and presets",
				UsageLine:   "mcpx search [query]",
				Args:        clihelp.MaximumNArgs(1),
				Run: func(ctx *clihelp.Context) error {
					q := ""
					if len(ctx.Args) > 0 {
						q = ctx.Args[0]
					}
					return cmdSearch(q)
				},
			},
			{
				Name:        "info",
				Aliases:     []string{"show"},
				Description: "Show server or preset definition and prerequisites",
				UsageLine:   "mcpx info <name>",
				Args:        clihelp.ExactArgs(1),
				Run:         func(ctx *clihelp.Context) error { return cmdInfo(ctx.Args[0]) },
			},
			{
				Name:        "registry",
				Description: "Manage remote template registry",
				UsageLine:   "mcpx registry <subcommand>",
				Subcommands: []clihelp.Command{
					{
						Name:        "update",
						Description: "Refresh local template cache from the remote repo",
						UsageLine:   "mcpx registry update",
						Args:        clihelp.NoArgs,
						Run:         func(ctx *clihelp.Context) error { return cmdTemplateUpdate() },
					},
				},
			},
			{
				Name:        "auth",
				Description: "Manage environment variables for MCP servers",
				UsageLine:   "mcpx auth <subcommand>",
				Notes: []clihelp.Note{
					{
						Heading: "Environment Variable Hierarchy",
						Text:    "Variables are resolved from three tiers, with later tiers taking precedence:\n1. System environment variables\n2. `~/.config/mcpx/.env` (global config)\n3. `{targetDir}/.env` (workspace-specific)",
					},
				},
				Examples: []clihelp.Example{
					{Line: "mcpx auth set GITHUB_TOKEN ghp_abc123", Description: "Store a GitHub token in the global config"},
					{Line: "mcpx auth get GITHUB_TOKEN", Description: "Display the masked value of GITHUB_TOKEN"},
					{Line: "mcpx auth list", Description: "List all stored environment variables"},
					{Line: "mcpx auth rm GITHUB_TOKEN", Description: "Remove GITHUB_TOKEN from the vault"},
				},
				Subcommands: []clihelp.Command{
					{
						Name:        "set",
						Description: "Set env var",
						UsageLine:   "mcpx auth set <var> [value]",
						Args:        clihelp.RangeArgs(1, 2),
						Run:         func(ctx *clihelp.Context) error { return cmdAuthSet(ctx.Args) },
					},
					{
						Name:        "get",
						Description: "Get env var value",
						UsageLine:   "mcpx auth get <var>",
						Args:        clihelp.ExactArgs(1),
						Run:         func(ctx *clihelp.Context) error { return cmdAuthGet(ctx.Args[0]) },
					},
					{
						Name:        "list",
						Description: "List all env vars",
						UsageLine:   "mcpx auth list [options]",
						Args:        clihelp.NoArgs,
						Options:     []clihelp.Option{clihelp.String(&dir, "--dir DIR", ".", "Target directory")},
						Run:         func(ctx *clihelp.Context) error { return cmdAuthList(dir) },
					},
					{
						Name:        "rm",
						Description: "Remove env var",
						UsageLine:   "mcpx auth rm <var>",
						Args:        clihelp.ExactArgs(1),
						Run:         func(ctx *clihelp.Context) error { return cmdAuthRm(ctx.Args[0]) },
					},
				},
			},
			{
				Name:        "help",
				Description: "Help about any command or topic",
				UsageLine:   "mcpx help [command|tree]",
				Subcommands: []clihelp.Command{
					{
						Name:        "tree",
						Description: "Show command tree of all commands",
						UsageLine:   "mcpx help tree",
						Args:        clihelp.NoArgs,
						Run: func(ctx *clihelp.Context) error {
							ctx.App.RenderTree(clihelp.Options{Writer: ctx.Stdout, Theme: ctx.App.Theme})
							return nil
						},
					},
				},
				Run: func(ctx *clihelp.Context) error {
					if len(ctx.Args) > 0 {
						if ctx.Args[0] == "tree" {
							ctx.App.RenderTree(clihelp.Options{Writer: ctx.Stdout, Theme: ctx.App.Theme})
							return nil
						}
						if !ctx.App.RenderCommand(clihelp.Options{Writer: ctx.Stdout, Theme: ctx.App.Theme}, ctx.Args...) {
							return fmt.Errorf("unknown help topic %q", ctx.Args[0])
						}
						return nil
					}
					ctx.App.RenderGlobal(clihelp.Options{Writer: ctx.Stdout, Theme: ctx.App.Theme})
					return nil
				},
			},
		},
	}
	return app
}

func isKnownCommand(s string) bool {
	switch s {
	case "add", "remove", "rm", "list", "ls", "sync", "update", "check", "test", "validate", "search", "info", "show", "registry", "template", "auth", "help":
		return true
	default:
		return false
	}
}

func cmdAdd(items []string, dir string, dryRun, overwrite, opencode, antigravity, cursor, claude, vscode, all, withRecommended bool) error {
	vault := envvault.NewVault(dir)
	servers, unknown := ResolveItems(items)
	if len(unknown) > 0 {
		fmt.Fprintf(os.Stderr, "Unknown items: %s\n", strings.Join(unknown, ", "))
	}
	if len(servers) == 0 {
		return fmt.Errorf("no valid servers or presets to add")
	}

	for _, s := range servers {
		for _, prereq := range s.Prerequisites {
			if !binaryExists(prereq.Binary) {
				fmt.Fprintf(os.Stderr, "Missing prerequisite %q for %s: %s\n", prereq.Binary, s.Name, prereq.Install)
			}
		}
	}
	for _, s := range servers {
		for name, spec := range s.Env {
			if spec.Required {
				if _, ok := vault.GetEnv(name); !ok {
					fmt.Fprintf(os.Stderr, "Missing required env var %s for %s: %s\n", name, s.Name, spec.Description)
				}
			}
		}
	}

	if dryRun {
		fmt.Printf("Would add %d server(s) to %s\n", len(servers), dir)
		for _, s := range servers {
			fmt.Printf("  - %s\n", s.Name)
		}
		return nil
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := writeTargets(servers, dir, overwrite, all, opencode, antigravity, cursor, claude, vscode); err != nil {
		return fmt.Errorf("failed to write targets: %w", err)
	}
	fmt.Printf("Added %d server(s) to %s\n", len(servers), dir)

	ensureGitignore(dir)

	for _, s := range servers {
		for _, rec := range s.Recommends {
			if withRecommended {
				if _, err := registry.GetServer(rec); err == nil {
					fmt.Printf("  Added recommended companion: %s\n", rec)
				}
			} else {
				fmt.Printf("  Tip: %s recommends %s (use --with-recommended to auto-add)\n", s.Name, rec)
			}
		}
	}
	return nil
}

func ensureGitignore(dir string) {
	if _, err := os.Stat(path.Join(dir, ".git")); err == nil {
		gitignorePath := path.Join(dir, ".gitignore")
		if data, err := os.ReadFile(gitignorePath); err == nil {
			if !strings.Contains(string(data), ".env") {
				f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
				if err == nil {
					defer f.Close()
					f.WriteString("\n.env\n")
				}
			}
		} else if os.IsNotExist(err) {
			_ = os.WriteFile(gitignorePath, []byte(".env\n"), 0644)
		}
	}
}

func writeTargets(servers []*types.Server, dir string, overwrite, all, opencode, antigravity, cursor, claude, vscode bool) error {
	vault := envvault.NewVault(dir)
	defaultTarget := !all && !opencode && !antigravity && !cursor && !claude && !vscode
	if all || opencode || defaultTarget {
		if err := config.NewOpenCodeWriter(dir, overwrite).Write(servers, vault); err != nil {
			return err
		}
	}
	if all || antigravity {
		if err := config.NewAntigravityWriter(dir, overwrite).Write(servers, vault); err != nil {
			return err
		}
	}
	if all || cursor || claude || vscode {
		if err := config.NewStandardWriter(dir, overwrite).Write(servers, vault); err != nil {
			return err
		}
	}
	return nil
}

func cmdRm(names []string, dir string) error {
	file := path.Join(dir, "opencode.json")
	root := config.LoadConfig(file)
	mcp, ok := root["mcp"]
	mcpObj := map[string]any{}
	if ok && config.IsObject(mcp) {
		mcpObj = mcp.(map[string]any)
	}
	changed := false
	for _, name := range names {
		if hasKey(mcpObj, name) {
			delete(mcpObj, name)
			fmt.Printf("Removed %s\n", name)
			changed = true
		} else {
			fmt.Fprintf(os.Stderr, "Not found: %s\n", name)
		}
	}
	if changed {
		root["mcp"] = mcpObj
		if err := config.WriteJSON(file, root); err != nil {
			return fmt.Errorf("failed to write config file %s: %w", file, err)
		}
	}
	return nil
}

func cmdList(dir string) error {
	file := path.Join(dir, "opencode.json")
	root := config.LoadConfig(file)
	mcp, ok := root["mcp"]
	if !ok || !config.IsObject(mcp) {
		fmt.Println("No servers configured.")
		return nil
	}
	mcpObj := mcp.(map[string]any)
	for _, name := range config.Keys(mcpObj) {
		fmt.Println(name)
	}
	return nil
}

func cmdInfo(name string) error {
	return cmdShow(name)
}

func cmdShow(name string) error {
	if server, err := registry.GetServer(name); err == nil {
		fmt.Printf("%s — %s\n", server.Name, server.Description)
		if server.URL != "" {
			fmt.Printf("  url: %s\n", server.URL)
		}
		fmt.Printf("  command: %s %s\n", server.Command, strings.Join(server.Args, " "))
		for _, prereq := range server.Prerequisites {
			fmt.Printf("  requires: %s\n", prereq.Binary)
		}
		return nil
	}
	if preset, err := registry.GetPreset(name); err == nil {
		fmt.Printf("Preset %s — %s\n", preset.Name, preset.Description)
		fmt.Printf("  servers: %s\n", strings.Join(preset.Servers, ", "))
		return nil
	}
	return fmt.Errorf("unknown template %q", name)
}

func cmdTest(names []string, timeoutSec int, verbose bool) error {
	timeout := time.Duration(timeoutSec) * time.Second
	var servers []*types.Server
	if len(names) == 0 {
		servers = allServers()
	} else {
		var unknown []string
		servers, unknown = ResolveItems(names)
		if len(unknown) > 0 {
			fmt.Fprintf(os.Stderr, "Unknown items: %s\n", strings.Join(unknown, ", "))
		}
	}
	if len(servers) == 0 {
		return fmt.Errorf("no servers to test")
	}
	r := runner.NewRunner(timeout)
	results := r.TestAll(servers)
	for _, res := range results {
		status := "OK"
		if !res.OK {
			status = "FAIL"
		}
		line := fmt.Sprintf("%-20s %-5s %6.2fs", res.Name, status, res.Elapsed.Seconds())
		if verbose || !res.OK {
			line += "  " + res.Message
		}
		fmt.Println(line)
	}
	return nil
}

func cmdSync(names []string, dir string) error {
	return cmdUpdate(names, dir)
}

func cmdUpdate(names []string, dir string) error {
	servers, _ := ResolveItems(names)
	if len(servers) == 0 {
		servers = allServers()
	}
	if err := writeTargets(servers, dir, false, true, false, false, false, false, false); err != nil {
		return err
	}
	fmt.Printf("Updated %d server(s) in %s\n", len(servers), dir)
	return nil
}

func cmdAuthSet(args []string) error {
	vault := envvault.NewVault(".")
	name := args[0]
	value := ""
	if len(args) > 1 {
		value = args[1]
	} else {
		value = os.Getenv(name)
	}
	if err := vault.SetEnv(name, value); err != nil {
		return err
	}
	fmt.Printf("Set %s\n", name)
	return nil
}

func cmdAuthGet(name string) error {
	vault := envvault.NewVault(".")
	if val, ok := vault.GetEnv(name); ok {
		fmt.Println(envvault.MaskValue(val))
	} else {
		fmt.Println("(unset)")
	}
	return nil
}

func cmdAuthList(dir string) error {
	vault := envvault.NewVault(dir)
	env := vault.ListEnv()
	var names []string
	for k := range env {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s=%s\n", name, envvault.MaskValue(env[name]))
	}
	return nil
}

func cmdAuthRm(name string) error {
	vault := envvault.NewVault(".")
	if err := vault.RemoveEnv(name); err != nil {
		return err
	}
	fmt.Printf("Removed %s\n", name)
	return nil
}

func cmdSearch(query string) error {
	query = strings.ToLower(strings.TrimSpace(query))

	servers := registry.ListServers()
	var matchedServers []*types.Server
	for _, name := range servers {
		server, err := registry.GetServer(name)
		if err != nil {
			continue
		}
		if query == "" || strings.Contains(strings.ToLower(server.Name), query) || strings.Contains(strings.ToLower(server.Description), query) {
			matchedServers = append(matchedServers, server)
		}
	}

	presets := registry.ListPresets()
	var matchedPresets []*types.Preset
	for _, name := range presets {
		preset, err := registry.GetPreset(name)
		if err != nil {
			continue
		}
		if query == "" || strings.Contains(strings.ToLower(preset.Name), query) || strings.Contains(strings.ToLower(preset.Description), query) {
			matchedPresets = append(matchedPresets, preset)
		}
	}

	if len(matchedServers) == 0 && len(matchedPresets) == 0 {
		fmt.Printf("No servers or presets matching %q\n", query)
		return nil
	}

	if len(matchedServers) > 0 {
		fmt.Println("Servers:")
		for _, server := range matchedServers {
			fmt.Printf("  \033[36m%s\033[0m - %s\n", server.Name, server.Description)
			cmdLine := server.Command
			if len(server.Args) > 0 {
				cmdLine += " " + strings.Join(server.Args, " ")
			}
			fmt.Printf("    Command: %s\n", cmdLine)
			if len(server.Prerequisites) > 0 {
				prereqs := make([]string, len(server.Prerequisites))
				for i, prereq := range server.Prerequisites {
					prereqs[i] = prereq.Binary
				}
				fmt.Printf("    Prerequisites: %s\n", strings.Join(prereqs, ", "))
			}
			if len(server.Recommends) > 0 {
				fmt.Printf("    Recommends: %s\n", strings.Join(server.Recommends, ", "))
			}
			fmt.Println()
		}
	}

	if len(matchedPresets) > 0 {
		fmt.Println("Presets:")
		for _, preset := range matchedPresets {
			fmt.Printf("  \033[33m%s\033[0m - %s\n", preset.Name, preset.Description)
			if len(preset.Servers) > 0 {
				fmt.Printf("    Includes: %s\n", strings.Join(preset.Servers, ", "))
			}
			fmt.Println()
		}
	}

	return nil
}

func cmdTemplateUpdate() error {
	n, err := registry.Update()
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d template(s) in %s\n", n, cacheDir())
	return nil
}

func binaryExists(binary string) bool {
	_, err := exec.Command(binary, "--version").Output()
	return err == nil
}

func allServers() []*types.Server {
	var servers []*types.Server
	for _, name := range registry.ListServers() {
		if s, err := registry.GetServer(name); err == nil {
			servers = append(servers, s)
		}
	}
	return servers
}

func hasKey(m map[string]any, key string) bool {
	for k := range m {
		if k == key {
			return true
		}
	}
	return false
}
