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

	if err := app.Execute(applyShorthands(os.Args[1:])); err != nil {
		clihelp.PrintError(err)
		os.Exit(1)
	}
}

// applyShorthands rewrites backward-compatible shorthand invocations into
// their canonical command forms before dispatch.
func applyShorthands(args []string) []string {
	if len(args) == 0 {
		return args
	}
	first := args[0]
	switch first {
	case "ls":
		return []string{"template", "list"}
	case "status":
		return []string{"list"}
	case "remove":
		return append([]string{"rm"}, args[1:]...)
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
		Version:     "0.1.1",
		GlobalNote:  "Run 'mcpx <command> --help' for command-specific options.",
		Commands: []clihelp.Command{
			{
				Name:        "add",
				Description: "Add servers/presets to workspace",
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
				Name:        "rm",
				Description: "Remove servers from workspace",
				UsageLine:   "mcpx rm <names...> [options]",
				Args:        clihelp.MinimumNArgs(1),
				Options:     []clihelp.Option{clihelp.String(&dir, "--dir DIR", ".", "Target directory")},
				Run:         func(ctx *clihelp.Context) error { return cmdRm(ctx.Args, dir) },
			},
			{
				Name:        "list",
				Description: "List active servers in workspace",
				UsageLine:   "mcpx list [options]",
				Args:        clihelp.NoArgs,
				Options:     []clihelp.Option{clihelp.String(&dir, "--dir DIR", ".", "Target directory")},
				Run:         func(ctx *clihelp.Context) error { return cmdList(dir) },
			},
			{
				Name:        "show",
				Description: "Show template definition",
				UsageLine:   "mcpx show <name>",
				Args:        clihelp.ExactArgs(1),
				Run:         func(ctx *clihelp.Context) error { return cmdShow(ctx.Args[0]) },
			},
			{
				Name:        "test",
				Description: "Test server handshakes (parallel)",
				UsageLine:   "mcpx test [names...] [options]",
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
				Name:        "update",
				Description: "Update configs with latest credentials",
				UsageLine:   "mcpx update [names...] [options]",
				Options:     []clihelp.Option{clihelp.String(&dir, "--dir DIR", ".", "Target directory")},
				Run:         func(ctx *clihelp.Context) error { return cmdUpdate(ctx.Args, dir) },
			},
			{
				Name:        "auth",
				Description: "Manage environment variables",
				UsageLine:   "mcpx auth <subcommand>",
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
				Name:        "template",
				Description: "List and show templates",
				UsageLine:   "mcpx template <subcommand>",
				Subcommands: []clihelp.Command{
					{
						Name:        "list",
						Description: "List available templates",
						UsageLine:   "mcpx template list",
						Args:        clihelp.NoArgs,
						Run:         func(ctx *clihelp.Context) error { return cmdTemplateList() },
					},
					{
						Name:        "show",
						Description: "Show template definition",
						UsageLine:   "mcpx template show <name>",
						Args:        clihelp.ExactArgs(1),
						Run:         func(ctx *clihelp.Context) error { return cmdShow(ctx.Args[0]) },
					},
					{
						Name:        "update",
						Description: "Refresh local template cache from the remote repo",
						UsageLine:   "mcpx template update",
						Args:        clihelp.NoArgs,
						Run:         func(ctx *clihelp.Context) error { return cmdTemplateUpdate() },
					},
				},
			},
		},
	}
	return app
}

func isKnownCommand(s string) bool {
	switch s {
	case "add", "rm", "list", "show", "test", "update", "auth", "template":
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

	writeTargets(servers, dir, overwrite, all, opencode, antigravity, cursor, claude, vscode)
	fmt.Printf("Added %d server(s) to %s\n", len(servers), dir)

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

func writeTargets(servers []*types.Server, dir string, overwrite, all, opencode, antigravity, cursor, claude, vscode bool) {
	vault := envvault.NewVault(dir)
	if all || opencode {
		config.NewOpenCodeWriter(dir, overwrite).Write(servers, vault)
	}
	if all || antigravity {
		config.NewAntigravityWriter(dir, overwrite).Write(servers, vault)
	}
	if all || cursor || claude || vscode {
		config.NewStandardWriter(dir, overwrite).Write(servers, vault)
	}
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
		config.WriteJSON(file, root)
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

func cmdUpdate(names []string, dir string) error {
	servers, _ := ResolveItems(names)
	if len(servers) == 0 {
		servers = allServers()
	}
	writeTargets(servers, dir, false, true, false, false, false, false, false)
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

func cmdTemplateList() error {
	fmt.Println("Servers:")
	for _, name := range registry.ListServers() {
		fmt.Printf("  %s\n", name)
	}
	fmt.Println("Presets:")
	for _, name := range registry.ListPresets() {
		fmt.Printf("  %s\n", name)
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
