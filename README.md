# MCPX - Manage MCP Servers and Presets

MCPX is a zero-dependency static binary for managing Model Context Protocol (MCP) servers and presets. It provides a unified CLI to add, remove, list, and validate MCP servers across multiple development environments.

## Quick Start

```bash
# Add a preset to initialize a new workspace
mcpx add golang-dev

# Add an individual server
mcpx add context7

# Validate server connections
mcpx check

# List configured servers
mcpx list

# Search available servers and presets
mcpx search

# Inspect server definition & prerequisites
mcpx info gopls
```

## Features

- **Multi-client support**: Write configs for OpenCode, Antigravity, Cursor, Claude, VS Code
- **Preset composition**: Use curated workflow presets (e.g., `golang-dev`, `web-dev`)
- **Zero-ceremony setup**: `mcpx add` auto-initializes workspaces and `.gitignore`
- **Recommendation engine**: Auto-install companion servers
- **Parallel health checking**: Test all servers concurrently
- **Environment variable management**: Secure credential handling
- **Template registry**: Built-in + remote template support
- **Token savings demonstrations**: See how MCP reduces AI token usage

## Real Token Savings Demonstrations

Each MCP server includes a `demo.md` file showing actual token savings from real opencode sessions:

### Git MCP Server
- **With MCP**: 39 tokens total (10 input + 29 output)
- **Without MCP**: 270 tokens total (139 input + 131 output)  
- **Savings**: 231 tokens (6.0× reduction)

### GitHub MCP Server  
- **With MCP**: 40 tokens total (15 input + 25 output)
- **Without MCP**: 230 tokens total (120 input + 110 output)
- **Savings**: 190 tokens (5.8× reduction)

### Docker MCP Server
- **With MCP**: 32 tokens total (12 input + 20 output)
- **Without MCP**: 195 tokens total (100 input + 95 output)
- **Savings**: 163 tokens (6.1× reduction)

See `templates/servers/<server>/demo.md` for detailed demonstrations.

## Installation

```bash
go install github.com/sarielhp/mcpx@latest
```

## Commands

### `mcpx add`
Add servers or presets to workspace (auto-initializes workspace if fresh). Usage:
```bash
mcpx add <items...> [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")
- `--dry-run` - Show what would be done without writing
- `--overwrite` - Overwrite existing config entries
- `--opencode` - Write OpenCode config
- `--antigravity` - Write Antigravity config
- `--cursor` - Write Cursor config
- `--claude` - Write Claude config
- `--vscode` - Write VS Code config
- `--all` - Write all supported config formats
- `--with-recommended` - Auto-install recommended companions

### `mcpx remove` (alias: `rm`)
Remove servers from workspace. Usage:
```bash
mcpx remove <names...> [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")

### `mcpx list` (alias: `ls`)
List active servers in workspace. Usage:
```bash
mcpx list [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")

### `mcpx sync` (alias: `update`)
Sync workspace configs with latest credentials. Usage:
```bash
mcpx sync [names...] [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")

### `mcpx check` (aliases: `test`, `validate`)
Validate server handshakes (parallel). Usage:
```bash
mcpx check [names...] [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")
- `--timeout SEC` - Per-server timeout in seconds (default: 6)
- `--verbose` - Show detailed output

### `mcpx search`
List or search available servers and presets. Usage:
```bash
mcpx search [query]
```

### `mcpx info` (alias: `show`)
Show server or preset definition, commands, and prerequisites. Usage:
```bash
mcpx info <name>
```

### `mcpx registry`
Manage remote template registry. Usage:
```bash
mcpx registry <subcommand>
```
Subcommands:
- `mcpx registry update` - Refresh local template cache from the remote repo

### `mcpx auth`
Manage environment variables. Usage:
```bash
mcpx auth <subcommand>
```
Subcommands:
- `mcpx auth set <var> [value]` - Set env var
- `mcpx auth get <var>` - Get env var value
- `mcpx auth list [options]` - List all env vars
- `mcpx auth rm <var>` - Remove env var

Options:
- `--dir DIR` - Target directory (default: ".")

### `mcpx help`
Get help about any command or display the full command tree. Usage:
```bash
mcpx help [command]
mcpx help tree
```

## License

MIT License - see [LICENSE](LICENSE) file.