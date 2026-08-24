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
ruby tools/mcpxsmoke.rb            # regression smoke tests (builds the binary)
```

The smoke suite is the primary correctness gate; run it after any change to
`main.go`, the config writers, or the template schemas.

## Development Tools (`tools/`)

All tools are standalone Ruby scripts (`#!/usr/bin/env ruby`) with `OptionParser`.
Do NOT create a README inside `tools/`; each script is self-documenting via
`--help` and top-level comments.

### Go API / Module discovery

These exist because Go 1.26 stdlib and module-cache APIs must be inspected
directly — never guess signatures or write throwaway probe programs.

**`tools/goapi.rb`** — Dump exported symbols (funcs/types/methods) with
signatures from Go source trees (stdlib or module cache).

```bash
ruby tools/goapi.rb --path /usr/lib/go-1.26/src/io --name Write
ruby tools/goapi.rb --path ~/.go/pkg/mod/github.com/sarielhp/clihelp@v0.2.14 --methods
ruby tools/goapi.rb --stdlib --name Exists --funcs
```

Flags: `--path DIR` (repeatable), `--stdlib`, `--name SUBSTR`, `--funcs`,
`--types`, `--methods`, `--json`.

**`tools/gomodcache.rb`** — Resolve installed module versions from the real
GOMODCACHE (`go env GOMODCACHE`, NOT stray mirrors like `~/.sandbox`).

```bash
ruby tools/gomodcache.rb clihelp                    # list cached versions
ruby tools/gomodcache.rb clihelp --latest           # newest cached version
ruby tools/gomodcache.rb clihelp --api v0.2.14 --name Execute --methods
ruby tools/gomodcache.rb --root                     # print GOMODCACHE path
```

Accepts full module paths or just the leaf package name. `--api` delegates to
`goapi.rb`.

**`tools/modcheck.rb`** — Verify every `require` in `go.mod` exists in the
local module cache at the required version (catches dependency drift before
compile).

```bash
ruby tools/modcheck.rb          # check ./go.mod
ruby tools/modcheck.rb --json
ruby tools/modcheck.rb --dry-run  # show what `go mod tidy` would change
ruby tools/modcheck.rb --fix      # run `go mod tidy`
```

### Correctness & validation

**`tools/mcpxsmoke.rb`** — Regression smoke test suite. Builds mcpx and runs
isolated scenarios in scratch dirs against `add`/`rm`/`list`/`show`/`test`/
`auth`/`template` and the shorthand rewrites (`ls`, `status`, `remove`, bare
items → `add`). Asserts on exit codes, stdout/stderr substrings, and JSON
config contents (with `sub:` to navigate nested keys like `mcp`/`mcpServers`).

```bash
ruby tools/mcpxsmoke.rb                          # full suite (builds binary)
ruby tools/mcpxsmoke.rb --bin /path/to/mcpx      # use prebuilt binary
ruby tools/mcpxsmoke.rb --select template,list   # run matching tests
ruby tools/mcpxsmoke.rb --verbose --fail-quick
```

Tests are a data table (`TESTS` constant) with `steps: [{args:, expect:}]`
and optional `post:` JSON checks. Add a new entry when adding a command or
flag.

**`tools/goldendif.rb`** — Golden-file diff for the config writers. Runs
`mcpx add <server> --all --overwrite` in a scratch dir and compares the three
generated files (`opencode.json`, `antigravity.json`, `mcp.json`) against
committed goldens under `testdata/golden/<server>/`. Catches structural output
regressions (not just key presence).

```bash
ruby tools/goldendif.rb                        # check all servers
ruby tools/goldendif.rb --server gopls         # one server
ruby tools/goldendif.rb --update               # (re)write goldens
ruby tools/goldendif.rb --pretty --verbose     # order-insensitive compare
ruby tools/goldendif.rb --golden /path/to/goldens
```

Flags: `--server NAME`, `--update`, `--pretty`, `--bin PATH`, `--golden DIR`,
`--json`, `--fail-fast`.

**`tools/mcpxprobe.rb`** — Standalone JSON-RPC initialize handshake probe.
Spawns a server command, sends the `initialize` request over stdin, reads the
response with a timeout. Mirrors the Go runner's handshake so you can debug a
server without rebuilding Go.

```bash
ruby tools/mcpxprobe.rb --command npx --args '-y @upstash/context7-mcp'
ruby tools/mcpxprobe.rb --cmd "gopls mcp" --timeout 8 --verbose --json
ruby tools/mcpxprobe.rb --command ./my-server --expect-tools initialize
```

Flags: `--command CMD`, `--args "S"` (repeatable), `--arg A` (repeatable),
`--cmd "LINE"`, `--timeout SEC` (default 6), `--verbose`, `--json`,
`--expect-tools LIST`.

### Git workflow

**`tools/snapshot.rb`** — Stage, commit, and push the working tree in one step.
Runs `git add -A`, creates a commit (from `--message` or an auto-generated
summary of the changed paths), and pushes to the branch's upstream. Exit 0 on
success, including a clean tree (benign no-op).

```bash
ruby tools/snapshot.rb                           # auto-message + push
ruby tools/snapshot.rb "fix: update writers"     # explicit message
ruby tools/snapshot.rb --no-push                 # commit only
ruby tools/snapshot.rb --dry-run --verbose       # preview intent
```

Flags: `-m/--message MSG` (positional arg also accepted), `--no-push`,
`--dry-run`, `--verbose`.

## Conventions

- **Scripting**: use Ruby for all scripts; never shell/awk/sed one-liners or
  Python for tooling.
- **No `tools/README.md`**: scripts self-document via top-level comments and
  `--help`.
- **Commit/push often**; the remote is `origin/master` on GitHub
  (`github.com/sarielhp/mcpx`).
- When adding a server template, also consider adding a `recommends` entry and
  a smoke/golden case.
