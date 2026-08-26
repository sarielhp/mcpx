# MCPX - Manage MCP Servers and Presets

MCPX is a zero-dependency static binary for managing Model Context Protocol (MCP) servers and presets. It provides a unified CLI to add, remove, list, and validate MCP servers across multiple development environments.

## Quick Start

```bash
# Initialize a new workspace with a preset
mcpx init --preset golang-dev

# Add a server
mcpx add context7

# Validate server connections
mcpx validate

# List configured servers
mcpx list
```

## Features

- **Multi-client support**: Write configs for OpenCode, Antigravity, Cursor, Claude, VS Code
- **Preset composition**: Use curated workflow presets (e.g., `golang-dev`, `web-dev`)
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

### `mcpx init`
Initialize a new workspace with a preset. Usage:
```bash
mcpx init [--preset NAME] [--dir DIR] [options]
```
Options:
- `--preset PRESET` - Preset to initialize with (default: "golang-dev")
- `--dir DIR` - Target directory (default: ".")
- `--overwrite` - Overwrite existing config files
- `--all` - Write all supported config formats

### `mcpx add`
Add servers/presets to workspace. Usage:
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

### `mcpx rm`
Remove servers from workspace. Usage:
```bash
mcpx rm <names...> [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")

### `mcpx list`
List active servers in workspace. Usage:
```bash
mcpx list [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")

### `mcpx show`
Show template definition. Usage:
```bash
mcpx show <name>
```

### `mcpx validate`
Validate server handshakes (parallel). Usage:
```bash
mcpx validate [names...] [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")
- `--timeout SEC` - Per-server timeout in seconds (default: 6)
- `--verbose` - Show detailed output

### `mcpx update`
Update configs with latest credentials. Usage:
```bash
mcpx update [names...] [options]
```
Options:
- `--dir DIR` - Target directory (default: ".")

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

### `mcpx template`
Manage templates. Usage:
```bash
mcpx template <subcommand>
```
Subcommands:
- `mcpx template list` - List available templates with descriptions and usage
- `mcpx template show <name>` - Show template definition
- `mcpx template update` - Refresh local template cache from the remote repo

## License

MIT License - see [LICENSE](LICENSE) file.