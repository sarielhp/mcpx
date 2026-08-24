# MCPX Implementation Plan

## Overview

This document outlines the implementation plan for `mcpx`, a Go rewrite of the `mcp-init` Ruby tool. The new implementation will be a zero-dependency static binary with enhanced features including parallel health checking, recommendation engine, and a hierarchical directory-based template system.

## Module & Dependencies

```
module github.com/sarielhp/mcpx
go 1.26.5
require github.com/sarielhp/clihelp v0.2.0
```

Zero external deps beyond `clihelp` (which itself pulls in `pflag`, `fatih/color`, `golang.org/x/term`). No `hujson` - embedded templates are clean JSON, and user overrides use stdlib `encoding/json`.

## Directory Structure

```
mcpx/
├── main.go
├── go.mod / go.sum
├── internal/
│   ├── types/types.go          # Core structs
│   ├── registry/
│   │   ├── registry.go         # Discovery + resolution
│   │   └── embed.go            # go:embed built-in templates
│   ├── envvault/envvault.go    # .env CRUD + multi-tier resolution
│   ├── config/
│   │   ├── opencode.go         # OpenCode format writer
│   │   ├── antigravity.go      # Antigravity format writer
│   │   └── standard.go         # Cursor/Claude/VS Code writer
│   └── runner/runner.go        # Parallel JSON-RPC health checker
└── templates/
    ├── servers/                 # 23 server dirs (git/, github/, ...)
    │   └── <name>/
    │       ├── template.json    # Server definition
    │       └── README.md        # Human docs (placeholder)
    └── presets/                 # 10 preset dirs (golang-dev/, ...)
        └── <name>/
            ├── preset.json      # Server references
            └── README.md        # Human docs (placeholder)
```

## Template Schema

**Server** (`templates/servers/<name>/template.json`):
```json
{
  "name": "semantic-scholar",
  "description": "Semantic Scholar academic paper & author search",
  "command": "uvx",
  "args": ["s2-mcp-server"],
  "env": {
    "SEMANTIC_SCHOLAR_API_KEY": {
      "value": "${SEMANTIC_SCHOLAR_API_KEY}",
      "required": true,
      "description": "API key from semanticscholar.org"
    }
  },
  "prerequisites": [
    {"binary": "uvx", "check": "uvx --version", "install": "curl -LsSf https://astral.sh/uv/install.sh | sh"}
  ],
  "recommends": [],
  "alternatives": [],
  "requires": [],
  "healthCheck": {
    "type": "stdio",
    "timeout": "5s",
    "expectTools": ["semantic_scholar_search_papers"]
  }
}
```

**Preset** (`templates/presets/<name>/preset.json`):
```json
{
  "name": "golang-dev",
  "description": "Go development environment with gopls, context7, docker, etc.",
  "servers": ["gopls", "context7", "docker", "fetch", "github", "git"]
}
```

## CLI Commands (using `clihelp`)

| Command | Args | Flags | Description |
|---|---|---|---|
| `add <items...>` | `MinimumNArgs(1)` | `--dir`, `--dry-run`, `--overwrite`, `--opencode`, `--antigravity`, `--cursor`, `--claude`, `--vscode`, `--all`, `--with-recommended` | Add servers/presets to workspace |
| `rm <names...>` | `MinimumNArgs(1)` | `--dir` | Remove servers from workspace |
| `list` | `NoArgs` | `--dir` | List active servers in workspace |
| `show <name>` | `ExactArgs(1)` | — | Show template definition |
| `test [names...]` | `MaximumNArgs(23)` | `--dir`, `--timeout`, `--verbose` | Test server handshakes (parallel) |
| `update [names...]` | — | `--dir` | Update configs with latest credentials |
| `auth set <var> [value]` | `RangeArgs(1,2)` | — | Set env var |
| `auth get <var>` | `ExactArgs(1)` | — | Get env var value |
| `auth list` | `NoArgs` | `--dir` | List all env vars |
| `auth rm <var>` | `ExactArgs(1)` | — | Remove env var |
| `template list` | `NoArgs` | — | List available templates |
| `template show <name>` | `ExactArgs(1)` | — | Show template definition |

**Backward-compatible shorthands** (handled via `BeforeRun` hook):
- `mcpx <items...>` → `mcpx add <items...>`
- `mcpx list` / `mcpx ls` → `mcpx template list`
- `mcpx status` → `mcpx list`
- `mcpx show <name>` → `mcpx template show <name>`
- `mcpx test` → `mcpx test`
- `mcpx remove` → `mcpx rm`
- `mcpx update` → `mcpx update`

## Implementation Order (13 steps)

| # | What | Files | Est. LOC |
|---|---|---|---|
| 1 | `go.mod`, `main.go` skeleton with `clihelp.App` definition | `go.mod`, `main.go` | 30 |
| 2 | Core types: `Server`, `Preset`, `EnvVar`, `Prerequisite`, `Recommendation`, `HealthCheck`, `ServerDef` | `internal/types/types.go` | 80 |
| 3 | Convert 23 servers to new `template.json` format with `recommends`/`alternatives`/`requires` metadata | `templates/servers/*/template.json` | 460 |
| 4 | Convert 10 presets to new `preset.json` format (server references only) | `templates/presets/*/preset.json` | 60 |
| 5 | `go:embed` for built-in templates | `internal/registry/embed.go` | 20 |
| 6 | Registry: `ListServers`, `ListPresets`, `GetServer`, `GetPreset`, `ResolveItems` | `internal/registry/registry.go` | 120 |
| 7 | Env vault: `LoadEnv`, `SetEnv`, `GetEnv`, `RemoveEnv`, `ListEnv`, `MaskValue` | `internal/envvault/envvault.go` | 100 |
| 8 | OpenCode config writer | `internal/config/opencode.go` | 60 |
| 9 | Antigravity config writer | `internal/config/antigravity.go` | 50 |
| 10 | Standard config writer (Cursor/Claude/VS Code) | `internal/config/standard.go` | 40 |
| 11 | Parallel health checker with goroutines + `context.WithTimeout` | `internal/runner/runner.go` | 120 |
| 12 | All subcommand handlers (`add`, `rm`, `list`, `show`, `test`, `update`, `auth`, `template`) | `internal/cli/*.go` (or inline in `main.go`) | 400 |
| 13 | Wire everything, `BeforeRun` for shorthands, recommendation output | `main.go` (finalize) | 80 |

**Total: ~1,620 LOC** (vs 1,995 in the Ruby original, but with proper modularity, parallel health checks, recommendation engine, and zero runtime deps)

## Key Design Decisions

1. **Parallel health checking**: Goroutine pool with `errgroup` + `context.WithTimeout`. All servers probed concurrently, results collected in a `sync.Map`. Default timeout 6s per server.

2. **Recommendation engine**: After `add`, iterate each server's `recommends` list and print non-blocking tips. `--with-recommended` flag auto-installs companions.

3. **Overleaf auth**: Not implemented as Go code. Each server template can have a `scripts/` field (string) pointing to an external script. For overleaf, the `template.json` will have a note: `"scripts": {"auth": "See old/mcp-templates/bin/overleaf-auth — extract Overleaf session cookie into .env"}`.

4. **Config writers**: Each writer merges with existing config (unless `--overwrite`). OpenCode uses `mcp` key with `type: "local"` / `command: [...]` / `enabled: true` format. Others use `mcpServers` key with `command`/`args`/`env` format.

5. **Env resolution**: System env → `~/.config/mcpx/.env` → `{targetDir}/.env`. Template `env` values use `${VAR}` syntax resolved against merged env.