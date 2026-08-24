# MCPX - Manage MCP Servers and Presets

MCPX is a zero-dependency static binary for managing Model Context Protocol (MCP) servers and presets. It provides a unified CLI to add, remove, list, and test MCP servers across multiple development environments.

## Quick Start

```bash
# Initialize a new workspace with a preset
mcpx init --preset golang-dev

# Add a server
mcpx add context7

# Test server connections
mcpx test

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

## Token Savings Demonstrations

Each MCP server includes a `demo.md` file showing how much token usage is reduced when using the server vs. traditional shell command explanations.

For example, the Git MCP server demonstrates:
- **Without MCP**: ~150 tokens to explain `git log --oneline --graph -5`
- **With MCP**: ~20 tokens to call `git_log({"max_count": 5, "oneline": true, "graph": true})`
- **Savings**: ~130 tokens (6.5× reduction)

See `templates/servers/<server>/demo.md` for detailed demonstrations.

## Installation

```bash
go install github.com/sarielhp/mcpx@latest
```

## Commands

- `mcpx init` - Initialize workspace with preset
- `mcpx add` - Add servers/presets to workspace
- `mcpx rm` - Remove servers from workspace
- `mcpx list` - List configured servers
- `mcpx show` - Show template definition
- `mcpx test` - Test server handshakes (parallel)
- `mcpx update` - Update configs with latest credentials
- `mcpx auth` - Manage environment variables
- `mcpx template` - Manage templates
- `mcpx repos` - List project repositories

## License

MIT License - see [LICENSE](LICENSE) file.