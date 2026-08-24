# AGENTS.md

Guidance for AI agents working in this repository.

## Project Overview

`mcpx` is a zero-dependency static binary (Go) that manages MCP (Model Context
Protocol) servers and presets. It provides a hierarchical directory-based
template system (`templates/servers/`, `templates/presets/`), a parallel
JSON-RPC health checker, an env-var vault, config writers for multiple
frontends (OpenCode, Antigravity, Cursor/Claude/VS Code), and a recommendation
engine.

Key facts:
- CLI framework: `github.com/sarielhp/clihelp` v0.2.14 (pflag-based).
- The embedded template registry lives in `embed.go` at the repo root because
  `go:embed` cannot reference parent directories (`..`); templates are
  `go:embed`'d from `templates/`.
- Binary is built with `go build .` — produces `./mcpx`.

## Build & Verify

```bash
go build ./...
go vet ./...
tools/mcpxsmoke            # regression smoke tests (builds the binary)
```

The smoke suite is the primary correctness gate; run it after any change to
`main.go`, the config writers, or the template schemas.

## Scripting & Task Automation

- **Primary Scripting Language**: Always use **Ruby** for all scripts,
  automation tasks, file manipulation, data processing, and utility tooling.
- **Avoid Alternatives**: Do not write scripts in Python, Bash/Shell, Awk,
  Sed, Perl, etc. Whenever a task would typically use a shell script,
  awk/sed one-liner, or Python script, implement it using Ruby instead.
- **Exception for Python**: Python is strictly restricted and may only be used
  when it is the sole reasonable option (e.g., interacting with a Python-only
  SDK/framework with no Ruby bindings or alternative).
- **Proactive Tooling in `tools/`**:
  - Whenever a task involves repetitive operations (batch file updates,
    multi-file inspection, AST queries, log analysis, schema transformations,
    line audits), **proactively write a standalone Ruby script under the
    `tools/` directory** rather than performing the task through repeated
    manual LLM tool calls.
  - Do NOT create a `README.md` inside `tools/`. Scripts must be
    self-documenting via top-level comment headers and `--help` options (using
    Ruby's `OptionParser`).
  - Use `#!/usr/bin/env ruby` for standalone executable scripts and run them
    directly (no `.rb` extension; `chmod +x`).
  - Ensure scripts are executable, idiomatic, handle errors cleanly, and
    output compact, actionable summaries.
  - Leverage Ruby's standard library modules (such as `fileutils`, `open3`,
    `json`, `csv`, `optparse`, `uri`, `net/http`) where appropriate.

## Development Tools (`tools/`)

All tools are standalone Ruby scripts (`#!/usr/bin/env ruby`), executable
directly (no `.rb` extension), using `OptionParser`. Do NOT create a README
inside `tools/`; each script is self-documenting via `--help` and top-level
comments.

### Go API / Module discovery

These exist because Go 1.26 stdlib and module-cache APIs must be inspected
directly — never guess signatures or write throwaway probe programs.

**`tools/goapi`** — Dump exported symbols (funcs/types/methods) with
signatures from Go source trees (stdlib or module cache).

```bash
tools/goapi --path /usr/lib/go-1.26/src/io --name Write
tools/goapi --path ~/.go/pkg/mod/github.com/sarielhp/clihelp@v0.2.14 --methods
tools/goapi --stdlib --name Exists --funcs
```

Flags: `--path DIR` (repeatable), `--stdlib`, `--name SUBSTR`, `--funcs`,
`--types`, `--methods`, `--json`.

**`tools/gomodcache`** — Resolve installed module versions from the real
GOMODCACHE (`go env GOMODCACHE`, NOT stray mirrors like `~/.sandbox`).

```bash
tools/gomodcache clihelp                    # list cached versions
tools/gomodcache clihelp --latest           # newest cached version
tools/gomodcache clihelp --api v0.2.14 --name Execute --methods
tools/gomodcache --root                     # print GOMODCACHE path
```

Accepts full module paths or just the leaf package name. `--api` delegates to
`goapi`.

**`tools/modcheck`** — Verify every `require` in `go.mod` exists in the
local module cache at the required version (catches dependency drift before
compile).

```bash
tools/modcheck          # check ./go.mod
tools/modcheck --json
tools/modcheck --dry-run  # show what `go mod tidy` would change
tools/modcheck --fix      # run `go mod tidy`
```

### Correctness & validation

**`tools/mcpxsmoke`** — Regression smoke test suite. Builds mcpx and runs
isolated scenarios in scratch dirs against `add`/`rm`/`list`/`show`/`test`/
`auth`/`template` and the shorthand rewrites (`ls`, `status`, `remove`, bare
items → `add`). Asserts on exit codes, stdout/stderr substrings, and JSON
config contents (with `sub:` to navigate nested keys like `mcp`/`mcpServers`).

```bash
tools/mcpxsmoke                          # full suite (builds binary)
tools/mcpxsmoke --bin /path/to/mcpx      # use prebuilt binary
tools/mcpxsmoke --select template,list   # run matching tests
tools/mcpxsmoke --verbose --fail-quick
```

Tests are a data table (`TESTS` constant) with `steps: [{args:, expect:}]`
and optional `post:` JSON checks. Add a new entry when adding a command or
flag.

**`tools/goldendif`** — Golden-file diff for the config writers. Runs
`mcpx add <server> --all --overwrite` in a scratch dir and compares the three
generated files (`opencode.json`, `antigravity.json`, `mcp.json`) against
committed goldens under `testdata/golden/<server>/`. Catches structural output
regressions (not just key presence).

```bash
tools/goldendif                        # check all servers
tools/goldendif --server gopls         # one server
tools/goldendif --update               # (re)write goldens
tools/goldendif --pretty --verbose     # order-insensitive compare
tools/goldendif --golden /path/to/goldens
```

Flags: `--server NAME`, `--update`, `--pretty`, `--bin PATH`, `--golden DIR`,
`--json`, `--fail-fast`.

**`tools/mcpxprobe`** — Standalone JSON-RPC initialize handshake probe.
Spawns a server command, sends the `initialize` request over stdin, reads the
response with a timeout. Mirrors the Go runner's handshake so you can debug a
server without rebuilding Go.

```bash
tools/mcpxprobe --command npx --args '-y @upstash/context7-mcp'
tools/mcpxprobe --cmd "gopls mcp" --timeout 8 --verbose --json
tools/mcpxprobe --command ./my-server --expect-tools initialize
```

Flags: `--command CMD`, `--args "S"` (repeatable), `--arg A` (repeatable),
`--cmd "LINE"`, `--timeout SEC` (default 6), `--verbose`, `--json`,
`--expect-tools LIST`.

### Git workflow

**`tools/snapshot`** — Stage, commit, and push the working tree in one step.
Runs `git add -A`, creates a commit (from `--message` or an auto-generated
summary of the changed paths), and pushes to the branch's upstream. Exit 0 on
success, including a clean tree (benign no-op).

```bash
tools/snapshot                           # auto-message + push
tools/snapshot "fix: update writers"     # explicit message
tools/snapshot --no-push                 # commit only
tools/snapshot --dry-run --verbose       # preview intent
```

Flags: `-m/--message MSG` (positional arg also accepted), `--no-push`,
`--dry-run`, `--check`, `--verbose`.

### Security

**`tools/scan_secrets`** — Security & privacy violation scanner. Checks for
leaked API keys, passwords, private keys, tokens, credit-card numbers, tracked
or present `.env` files, and literal (non-`${VAR}`) values in `template.json`
env fields, plus preset→server reference validation.

```bash
tools/scan_secrets                # scan whole tree (exit 1 on violations)
tools/scan_secrets --path dir      # scan a directory
tools/scan_secrets --json          # CI-friendly JSON output
```

Installed as a pre-commit hook (`.githooks`), also invoked by
`tools/snapshot --check`. Run it after any change to templates or config
writers, before pushing.

## Conventions

- See the **Scripting & Task Automation** section above for the Ruby-only
  scripting rules and the `tools/` conventions.
- **Commit/push often**; the remote is `origin/master` on GitHub
  (`github.com/sarielhp/mcpx`).
- When adding a server template, also consider adding a `recommends` entry and
  a smoke/golden case.
- Templates live in the **`mcp-templates` GitHub repo**
  (`github.com/sarielhp/mcp-templates`) and are mirrored in `templates/`.
  Run `mcpx template update` to refresh the local cache from that repo.
